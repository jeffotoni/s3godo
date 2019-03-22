// Go in action
// @jeffotoni
// 2019-03-14

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

var (
	pat = os.Getenv("DO_PAT")
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

func main() {
	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)

	fmt.Println("auth ok")
	// fmt.Println(client)

	fmt.Println("###################################")

	fmt.Println("Create new Droplet.")
	dropletName := "super-cool-dit10-droplet"

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: "nyc1",
		Size:   "s-1vcpu-1gb",
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-18-04-x64",
		},
		SSHKeys: []godo.DropletCreateSSHKey{
			{Fingerprint: "7d:c1:35:d1:2d:90:2a:16:ec:ff:d1:22:b8:c7:e2:27"},
		},
	}

	ctx := context.TODO()
	newDroplet, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		fmt.Printf("Something bad happened: %s\n\n", err)
		return
	}

	fmt.Println("create: ", newDroplet)
	fmt.Println("Publish in prepare in Droplet.")
	time.Sleep(time.Second * 5)
	fmt.Println("Success!")

}
