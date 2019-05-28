package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "io/ioutil"
    "os/user"
)

var (
    BUCKET   = ""
    HOME_DIR = ""
)

// DOKey contem dados para autenticacao na Digital Ocean(acho).
type DOKey struct {
    Key      string `json:"key"`
    Secret   string `json:"secret"`
    Endpoint string `json:"endpoint"`
    Region   string `json:"region"`
    Bucket   string `json:"bucket"`
}

func init() {
    user, err := user.Current()
    if err != nil {
        return
    }
    HOME_DIR = user.HomeDir
}

func main() {

    // abir file secrey
    key, err := ReadKey()
    if err != nil {
        fmt.Println("Erro ao montar suas credenciais de acesso ao DigitalOcean Space!")
        return
    }

    fbucket := flag.String("bucket", "", "o nome do seu bucket")
    flag.Parse()

    if len(*fbucket) > 0 {
        BUCKET = *fbucket
    } else {
        fmt.Println("[Erro Bucket é obrigatorio!]")
        flag.PrintDefaults()
        return
    }

    // Initialize a client using Spaces
    s3Config := &aws.Config{
        Credentials: credentials.NewStaticCredentials(key.Key, key.Secret, ""),
        Endpoint:    aws.String(key.Endpoint),
        Region:      aws.String(key.Region), // This is counter intuitive, but it will fail with a non-AWS region name.
    }

    newSession := session.New(s3Config)
    s3Client := s3.New(newSession)

    input := &s3.ListObjectsInput{
        Bucket: aws.String(BUCKET),
        //Key:    aws.String("/s3wfcore/config/configambiente.conf.php"),
        //MaxKeys: aws.Int64(0),
    }

    result, err := s3Client.ListObjects(input)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case s3.ErrCodeNoSuchBucket:
                fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
        return
    }

    fmt.Println("Listando Objects:")
    //fmt.Println(result)

    for _, item := range result.Contents {
        fmt.Println("Name:         ", *item.Key)
        fmt.Println("Last modified:", *item.LastModified)
        fmt.Println("Size:         ", *item.Size)
        fmt.Println("Storage class:", *item.StorageClass)
        fmt.Println("")
    }

    fmt.Println("Found", len(result.Contents), "items in bucket", BUCKET)
    fmt.Println("")
}

func ReadKey() (*DOKey, error) {
    // user, err := user.Current()
    // if err != nil {
    //  return nil, err
    // }
    cfgFile := fmt.Sprintf("%s/.dokeys", HOME_DIR)
    b, err := ioutil.ReadFile(cfgFile)
    if err != nil {
        return nil, err
    }

    key := &DOKey{}
    if err := json.Unmarshal(b, key); err != nil {
        return nil, err
    }
    return key, nil
}
