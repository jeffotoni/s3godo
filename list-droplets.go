// Go in action
// @jeffotoni
// 2019-03-14

package main

import (
	"context"
	"fmt"
	"os"

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

	fmt.Println("...........................")
	fmt.Println("Lista all Droplets")

	ctx := context.TODO()

	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	droplets, _, err := client.Droplets.List(ctx, opt)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(droplets)
	// all droplet
	fmt.Println("...........................")
	for _, v := range droplets {
		fmt.Println("Id: ", v.ID)
		fmt.Println("Name: ", v.Name)
		fmt.Println("Memory: ", v.Memory)
	}

}
