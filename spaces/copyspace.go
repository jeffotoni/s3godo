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
    . "github.com/jeffotoni/gcolor"
    "io/ioutil"
    "net/http"
    "os"
    "os/signal"
    "path/filepath"
    "strings"
    "sync"
    "time"
)

const (
    ACL  = "public-read-write"
    ACLp = "private"
)

var (
    ACL_AP = ACL
    BUCKET = ""
)

// -ldflags "-X main.k=your-key -X main.s=your-secret" main.go
// var (
//     endpoint, region, bucket, key, secret string
// )

// func int() {
//     if len(secret) == 0 || len(key) == 0 || len(endpoint) == 0 || len(region) == 0 {
//         fmt.Println("Error fatal, sua chave não estão compiladas aqui.")
//         os.Exit(0)
//     }
// }

func main() {

    //println(k, s)
    //return

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
    var pathFile string

    //pathFile := flag.String("file", "", "nome do arquivo ou diretorio a ser enviado")
    flag.StringVar(&pathFile, "file", "", "nome do arquivo ou diretorio a ser enviado")
    aclSend := flag.String("acl", "", "permissao: public or private")
    fbucket := flag.String("bucket", "", "o nome do seu bucket")
    flag.Parse()

    if len(pathFile) == 0 {
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

    var wg sync.WaitGroup

    if DirExist(pathFile) {

        type sendS3 struct {
            Path     string
            Pbucket  string
            S3Client *s3.S3
        }

        //c := make(chan sendS3)
        done := make(chan string)
        dir := pathFile
        go func() {
            err := filepath.Walk(dir,
                func(path string, info os.FileInfo, err error) error {
                    if err != nil {
                        return err
                    }
                    pbucket := strings.Replace(path, os.Getenv("HOME"), "", -1)
                    //cy := sendS3{Path: path, Pbucket: pbucket, S3Client: s3Client}
                    //c <- cy

                    wg.Add(1)
                    SendFileDo(cx.Path, cx.Pbucket, cx.S3Client, &wg)
                    return nil
                })
            if err != nil {
                fmt.Println(err)
            }
            //else {
            //close(c)
            //done <- "fim de envio"
            //}
        }()

        wg.Wait()

        // for cx := range c {
        //     SendFileDo(cx.Path, cx.Pbucket, cx.S3Client)
        // }

        println(<-done)
        return

    } else {

        // bucket
        pbucket := strings.Replace(pathFile, os.Getenv("HOME"), "", -1)
        p := pathFile

        wg.Add(1)
        SendFileDo(p, pbucket, s3Client, &wg) // send one file
        wg.Wait()
    }
}

func SendFileDo(pf, pbucket string, s3Client *s3.S3, wg *sync.WaitGroup) {

    t1 := time.Now()
    f, err := os.Open(pf)
    if err != nil {
        fmt.Print(err)
        return
    }
    defer f.Close()

    // size file...
    fi, err := f.Stat()
    if err != nil {
        fmt.Println(err)
        return
    }

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

    bs := string(b)

    // runer
    timer := RunerTimer()

    type Fs struct {
        Msgs3 *string
        Name  string
        Size  int64
    }

    //var msgs3 = make(chan *string)
    var cfs = make(chan Fs)

    //var wg = &sync.WaitGroup{}
    //wg.Add(1)
    // aqui deveria ter um worker
    func(pf, b, contentType string) {
        //println(pf)
        wg.Done()
        pathV := strings.Split(pf, "/")
        lastp := len(pathV)
        nameFileSpace := pathV[lastp-1]

        // Upload a file to the Space
        object := s3.PutObjectInput{
            ACL:         aws.String(ACL_AP),
            Body:        strings.NewReader(b),
            Bucket:      aws.String(BUCKET),
            Key:         aws.String(pbucket),
            ContentType: aws.String(contentType),
        }
        msgs3V, err := s3Client.PutObject(&object)
        if err != nil {
            fmt.Println(err.Error())
            return
        }
        cfs <- Fs{Msgs3: msgs3V.ETag, Name: nameFileSpace, Size: fi.Size()}
        close(cfs)
        time.Sleep(time.Millisecond * 10)
    }(pf, bs, contentType)

    //wg.Wait()
    <-timer

    csfS := <-cfs
    kb := (csfS.Size / 1024)
    t2 := time.Now()

    fmt.Print("\r")
    fmt.Print("\033[?25h")
    fmt.Println("[send success] Id["+*csfS.Msgs3+"] File["+csfS.Name+"] Size[", kb, "Kb]", "time[", t2.Sub(t1), "]")
    fmt.Print("\r")
    fmt.Print("\033[?25h")
    fmt.Print("\033[?25h")
    return
}

func DirExist(path string) bool {

    //if _, err := os.Stat(path); err == nil {
    if stat, err := os.Stat(path); err == nil && stat.IsDir() {
        return true
    }

    return false
}

func ReadKey() (endpoint, region, bucket, key, secret string, err error) {
    pathHome := os.Getenv("HOME")
    pathHome = pathHome + "/.dokeys"
    b, err := ioutil.ReadFile(pathHome) // just pass the file name
    if err != nil {
        fmt.Print("keys: ", err)
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

func RunerTimer() <-chan time.Time {

    timer := time.Tick(time.Duration(50) * time.Millisecond)

    go func() {
        sc := make(chan os.Signal, 1)
        signal.Notify(sc, os.Interrupt)

        <-sc

        fmt.Print(RedCor("\ncanceled!"))
        fmt.Print("\033[?25h")
        os.Exit(0)
    }()

    fmt.Print("\033[?25l")

    s := []rune(`|/~\`)
    //s := []rune(`-=*=`)
    //s := []rune(`◐◓◑◒`)
    i := 0

    go func() {
        for {

            <-timer
            fmt.Print("\r")
            fmt.Print(YellowCor(string(s[i])))
            i++
            if i == len(s) {
                i = 0
            }
        }
    }()

    return timer
}
