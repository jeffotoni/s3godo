# s3godo

client Digital Ocean in Go

Godo is a sdk, a client so we can work natively using Go.
The lib is very light and lean, we managed to handle via API all the features we have in Cloud DigitalOcean.

We were able to create, remove, list Droplets dynamically, provision them as needed using Godo, and integrate with Terraform or create their own apis.

Possibilities of working with Load Balance, Network, Spaces, Volumes, and so on.

The API is fantastic, network latency makes it all work like a rocket.

The programs in Go fit like a glove due to its simplicity and performance, consuming little hardwares we have lean costs for our projects.


Example of how to connect to DigitalOcean using Godo:

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
## DROPLET

```go
...
mt.Println("Create new Droplet.")
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
...
```

## SPACES
To access space, we will not use godo to use sdk aws.
in the simple space will find the source of the copyspace a program that will send files to the space.


To access space, we will not use godo to use sdk aws.
in the simple space will find the source of the copyspace a program that will send files to the space.


### Install with wget

You need to set up a hidden file with the name of .dokeys in your $HOME or ~/.dokeys
Without it nothing will work, to generate the keys you need at 
digitalocean.com in API -> Spaces access keys and generate your key.

```bash

$ sh -c "$(wget https://raw.githubusercontent.com/jeffotoni/s3godo/master/spaces/v1/install.sh -O -)"

```

### Install manual

For copyspace to work you need to generate a json file with the hidden .dokeys name in your $HOME, and its contents are:

```bash

$ echo "
{
     "key": "key-digitalocean",
     "secret": "secret-digitalocean",
     "endpoint": "https://your-space.digitaloceanspaces.com",
     "region": "us-east-1",
     "bucket": "your-bucket-default"
} " > $HOME/.dokeys

```

The bucket field is not required, and the keys you will be able to generate from the DigitalOcean control panel in Spaces access keys at https://cloud.digitalocean.com

### Install with Go
It is now install and use.

```bash

$ git clone https://github.com/jeffotoni/s3godo.git
$ cd s3godo/space/copyspace
$ go install

```

### Go build

```bash

$ git clone https://github.com/jeffotoni/s3godo.git
$ cd s3godo/space/copyspace
$ go build -ldflags="-s -w" -o copyspace main.go
$ cp copyspace $GOPATH/bin
$ copyspace -h

```

```bash

# The parameters are:
# file: filename
# acl: public or private
# bucket: the name of your bucket
# worker: simultaneous works
$ copyspace --file=your-file.pdf --acl=public --bucket=your-bucket --worker=100

```
