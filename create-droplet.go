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
	dropletName := "super-cool-jeff-droplet"

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: "nyc1",
		Size:   "s-1vcpu-1gb",
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-18-04-x64",
		},
		SSHKeys: []godo.DropletCreateSSHKey{
			{Fingerprint: "3d:c9:31:d4:4d:58:3a:66:ea:ff:d2:34:b9:c8:e9:36"},
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
