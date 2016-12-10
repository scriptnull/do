package main

import (
	"encoding/json"
	"fmt"
	"github.com/devbot-xyz/do/actions"
	"github.com/devbot-xyz/do/doproxy"
	"github.com/digitalocean/godo"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Message struct {
	Token   string      `json:"token"`
	Action  string      `json:"action"`
	UserId  string 			`json:"userid"`
	Payload string	    `json:"payload"`
}

func main() {
	conn, err := amqp.Dial("amqp://devbotuser:devbotpass@172.17.0.1:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"worker", // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"reactor.do", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		"worker", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var inMessage Message

			err = json.Unmarshal(d.Body, &inMessage)
			if err != nil {
				fmt.Println("Error converting json")
				fmt.Println(err)
				return
			}

			log.Printf("STATE: starting, ACTION: %s, SLACK_USER_ID: %s", inMessage.Action, inMessage.UserId)
			doClient := doproxy.GetDoClient(inMessage.Token)
			if doClient == nil {
				fmt.Println("Client instantiation failed. Unable to access digital ocean API")
				return
			}
			actionResult := doAction(doClient, inMessage.Action, inMessage.Payload)
			actionResult.UserId = inMessage.UserId
			log.Printf("STATE: ended, ACTION: %s, SLACK_USER_ID: %s", inMessage.Action, inMessage.UserId)

			messageToBeSent, _ := json.Marshal(actionResult)
			publishActionResult(ch, messageToBeSent)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func doAction(doClient *godo.Client, action string, payload string) (actionResult actions.ActionResult) {
	if action == "getDroplets" {
		actionResult = actions.GetDroplets(doClient, payload)
	}

	if action == "createDroplet" {
		actionResult = actions.CreateDroplet(doClient, payload)
	}

	return
}

func publishActionResult(ch *amqp.Channel, messageToBeSent []byte) {
	q, err := ch.QueueDeclare(
		"zender", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        messageToBeSent,
		})
	failOnError(err, "Failed to publish a message")
}
