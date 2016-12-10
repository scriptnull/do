package doproxy

import (
  "github.com/digitalocean/godo"
  "golang.org/x/oauth2"
  "strings"
)

type TokenSource struct {
  AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
    token := &oauth2.Token{
        AccessToken: t.AccessToken,
    }
    return token, nil
}

func GetDoClient(token string) *godo.Client {
    tokenSource := &TokenSource{
        AccessToken: token,
    }
    oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
    client := godo.NewClient(oauthClient)
    return client
}

func GetFormattedDroplet (box *godo.Droplet) string {
  boxDetails := ""
  boxDetails += " - " + box.Name + " ("
  boxDetails += "Memory: " + box.SizeSlug + "; "
  boxDetails += "Region: " + box.Region.Name + "; "
  boxDetails += "Image: " + box.Image.Slug
  if len(box.Tags) > 0 {
    boxDetails += "; Tags: " + strings.Join(box.Tags, ",")
  }
  boxDetails += " )"
  return boxDetails
}
