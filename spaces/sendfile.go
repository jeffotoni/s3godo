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
    endpoint, region, bucket, key, secret, err := ReadKey()

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

    // agora capturando dados..

    pathFile := flag.String("file", "", "nome do arquivo a ser enviado")
    flag.Parse()
    if len(*pathFile) == 0 {
        flag.PrintDefaults()
        return
    }

    b, err := ioutil.ReadFile(*pathFile) // just pass the file name
    if err != nil {
        fmt.Print(err)
        return
    }

    if len(string(b)) == 0 {
        fmt.Println("Error file está vazio..")
        return
    }

    //nome de arquivo...
    pathV := strings.Split(*pathFile, "/")
    lastp := len(pathV)
    nameFileSpace := pathV[lastp-1]

    // Upload a file to the Space
    object := s3.PutObjectInput{
        Body:   strings.NewReader(string(b)),
        Bucket: aws.String(bucket),
        Key:    aws.String(nameFileSpace),
    }
    msgs3, err := s3Client.PutObject(&object)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println("Enviando com sucesso!")
    fmt.Println(msgs3)
}

func ReadKey() (endpoint, region, bucket, key, secret string, err error) {

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
        Bucket   string `json:"bucket"`
    }

    sk := &skey{}
    if err = json.Unmarshal(b, sk); err != nil {
        return
    }

    key = sk.Key
    secret = sk.Secret
    endpoint = sk.Endpoint
    region = sk.Region
    bucket = sk.Bucket
    return
}
