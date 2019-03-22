# s3godo
client Digital Ocean in Go


Godo is a sdk, a client so we can work natively using Go.
The lib is very light and lean, we managed to handle via API all the features we have in Cloud DigitalOcean.

We were able to create, remove, list Droplets dynamically, provision them as needed using Godo, and integrate with Terraform or create their own apis.

Possibilities of working with Load Balance, Network, Spaces, Volumes, and so on.

The API is fantastic, network latency makes it all work like a rocket.

The programs in Go fit like a glove due to its simplicity and performance, consuming little hardwares we have lean costs for our projects.


```go

package main

import (
	"context"
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
}
```