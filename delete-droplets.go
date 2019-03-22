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

	fmt.Println("Del Droplets")

	ctx := context.TODO()

	_, err := client.Droplets.Delete(ctx, xxxxxxxxxxx)
	if err != nil {
		fmt.Println(err)
	}

}