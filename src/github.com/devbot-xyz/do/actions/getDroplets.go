package actions

import (
	"github.com/digitalocean/godo"
	"github.com/devbot-xyz/do/doproxy"
)

// Endpoint is https://api.digitalocean.com/v2/droplets?page=1&per_page=200

func GetDroplets(client *godo.Client, payload string) ActionResult {
	// log.Println("Inside GetDroplets")
	// defer log.Println("Completed GetDroplets")

  actionResult := ActionResult{}

	list := []godo.Droplet{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}

	for {
		droplets, resp, err := client.Droplets.List(opt)
		if err != nil {
      actionResult.Error = err
			return actionResult
		}

		// append the current page's droplets to our list
		for _, d := range droplets {
			list = append(list, d)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
    if err != nil {
      actionResult.Error = err
			return actionResult
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

  actionResult.Response = "Here are your machines \n"

  for _, box := range list {
    boxDetails := doproxy.GetFormattedDroplet(&box)
    actionResult.Response += boxDetails
  }

	return actionResult
}
