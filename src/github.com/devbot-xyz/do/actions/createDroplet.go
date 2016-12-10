package actions

import (
	"github.com/digitalocean/godo"
  "github.com/devbot-xyz/do/doproxy"
  "encoding/json"
  "fmt"
)

type DropletCreatePayload struct {
  Name string `json:"name"`
  Region string `json:"region"`
  Size string `json:"size"`
  Image string `json:"image"`
}

// Endpoint is https://api.digitalocean.com/v2/droplets?page=1&per_page=200

func CreateDroplet(client *godo.Client, payload string) ActionResult {
  actionResult := ActionResult{}
  var dropletPayoad = DropletCreatePayload{}

  json.Unmarshal([]byte(payload), &dropletPayoad)

  fmt.Printf("%+v", dropletPayoad)

  createRequest := &godo.DropletCreateRequest{
      Name:   dropletPayoad.Name,
      Region: dropletPayoad.Region,
      Size:   dropletPayoad.Size,
      Image: godo.DropletCreateImage{
          Slug: dropletPayoad.Image,
      },
  }

  newDroplet, _, err := client.Droplets.Create(createRequest)
  if err != nil {
    actionResult.Error = err
    return actionResult
  }

  actionResult.Response = "Created Droplet :: " + doproxy.GetFormattedDroplet(newDroplet)

	return actionResult
}
