package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
)

const (
    ACL  = "public-read-write"
    ACLp = "private"
)

var (
    ACL_AP = ACL
    BUCKET = ""
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
    aclSend := flag.String("acl", "", "permissao: public or private")
    fbucket := flag.String("bucket", "", "o nome do seu bucket")
    flag.Parse()

    if len(*pathFile) == 0 {
        flag.PrintDefaults()
        return
    }

    if len(*aclSend) > 0 && *aclSend != "public" {
        ACL_AP = ACLp
    }

    if len(*fbucket) > 0 {
        BUCKET = *fbucket
    } else {
        BUCKET = bucket
    }

    f, err := os.Open(*pathFile)
    if err != nil {
        fmt.Print(err)
        return
    }
    defer f.Close()

    // size file...
    // fi, err := f.Stat()
    // if err != nil {
    //     fmt.Println(err)
    //     return
    // }

    //// Use bufio.NewReader to get a Reader.
    // ... Then use ioutil.ReadAll to read the entire content.
    reader := bufio.NewReader(f)
    b, err := ioutil.ReadAll(reader)
    if err != nil {
        fmt.Println("readAll:", err)
        return
    }

    contentType, err := GetFileContentType(f)
    if err != nil {
        fmt.Println("contentType: ", err)
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
        ACL:         aws.String(ACL_AP),
        Body:        strings.NewReader(string(b)),
        Bucket:      aws.String(BUCKET),
        Key:         aws.String(nameFileSpace),
        ContentType: aws.String(contentType),
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

func GetFileContentType(out *os.File) (string, error) {

    // garante que irá
    // ler do inicio
    out.Seek(0, 0)
    // Only the first 512 bytes are used to sniff the content type.
    buffer := make([]byte, 512)
    _, err := out.Read(buffer)
    if err != nil {
        return "", err
    }

    // Use the net/http package's handy DectectContentType function. Always returns a valid
    // content-type by returning "application/octet-stream" if no others seemed to match.
    contentType := http.DetectContentType(buffer)
    return contentType, nil
}
