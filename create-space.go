package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "strings"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

func main() {

    // abir file secrey
    endpoint, region, key, secret, err := ReadKey()

    if err != nil {
        fmt.Println("Erro ao montar suas credenciais de acesso ao DigitalOcean Space!")
        return
    }

    // Initialize a client using Spaces
    s3Config := &aws.Config{
        Credentials: credentials.NewStaticCredentials(key, secret, ""),
        Endpoint:    aws.String(endpoint),
        Region:      aws.String(region), // This is counter intuitive, but it will fail with a non-AWS region name.
    }

    newSession := session.New(s3Config)
    s3Client := s3.New(newSession)

    bucketSapce := "your-bucket"
    //Create a new Space
    params := &s3.CreateBucketInput{
        Bucket: aws.String(bucketSapce),
    }

    _, err := s3Client.CreateBucket(params)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // List all Spaces in the region
    spaces, err := s3Client.ListBuckets(nil)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    for _, b := range spaces.Buckets {
        fmt.Printf("%s\n", aws.StringValue(b.Name))
    }
}

func ReadKey() (endpoint, region, key, secret string, err error) {

    b, err := ioutil.ReadFile("./.dokeys") // just pass the file name
    if err != nil {
        fmt.Print(err)
        return
    }

    //jsonkey := string(b) // convert content to a 'string'
    type skey struct {
        Key      string `json:"key"`
        Secret   string `json:"secret"`
        Endpoint string `json:"endpoint"`
        Region   string `json:"region"`
    }

    sk := &skey{}
    if err = json.Unmarshal(b, sk); err != nil {
        return
    }

    key = sk.Key
    secret = sk.Secret
    endpoint = sk.Endpoint
    region = sk.Region
    return
}
