package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	//"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	. "github.com/jeffotoni/gcolor"
)

const (
	ACL  = "public-read-write"
	ACLp = "public-read-write"
)

var (
	ACL_AP = ACL
	BUCKET = ""
	WORKER = "500"
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

type sendS3 struct {
	Path     string
	Pbucket  string
	S3Client *s3.S3
	I        int
}

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
	workers := flag.String("worker", "", "quantidade de trabalhos concorrentes em sua máquina")
	flag.Parse()

	if len(pathFile) == 0 {
		flag.PrintDefaults()
		return
	}

	if len(*aclSend) > 0 && strings.ToLower(*aclSend) != "private" {
		ACL_AP = ACLp
	}

	if len(*fbucket) > 0 {
		BUCKET = *fbucket
	} else {
		BUCKET = bucket
	}

	if len(*workers) > 0 {
		WORKER = *workers
	}

	println(CyanCor("domain: " + endpoint))
	println(YellowCor("bucket: " + BUCKET))

	// var wg sync.WaitGroup

	if DirExist(pathFile) {

		workeri, _ := strconv.Atoi(WORKER)
		jobs := make(chan sendS3)
		results := make(chan string)
		//done := make(chan string, 1)

		dir := pathFile
		var i int
		//wg.Add(5)
		// inicia o worker
		for w := 1; w <= workeri; w++ {
			go worker(w, jobs, results)
		}

		go func(i int) {
			err := filepath.Walk(dir,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					i++
					pbucket := strings.Replace(path, os.Getenv("HOME"), "", -1)
					jobsy := sendS3{Path: path, Pbucket: pbucket, S3Client: s3Client, I: i}
					jobs <- jobsy

					//wg.Add(1)
					// go SendFileDo(path,
					//     pbucket,
					//     s3Client,
					//     &wg,
					// )

					return nil
				})
			if err != nil {
				fmt.Println(err)
			}
		}(i)

		// wg.Wait()

		defer close(jobs)
		defer close(results)

		for cx := range results {
			println(cx)
		}

		//done <- "fim de envio"
		println("fim de envio")
		//wg.Wait()
		// println(<-done)
		return

	} else {

		// bucket
		pbucket := strings.Replace(pathFile, os.Getenv("HOME"), "", -1)
		p := pathFile
		//wg.Add(1)
		fmt.Println(SendFileDo(p, pbucket, s3Client, 1)) // send one file
		//wg.Wait()
	}
}

func worker(id int, jobs <-chan sendS3, results chan<- string) {
	for j := range jobs {
		time.Sleep(time.Millisecond * 2)
		results <- SendFileDo(j.Path, j.Pbucket, j.S3Client, j.I)
	}
}

func SendFileDoTest(pf, pbucket string, s3Client *s3.S3) string {
	time.Sleep(time.Second)
	return `send success: [` + pbucket + `]`
}

func SendFileDo(pf, pbucket string, s3Client *s3.S3, I int) string {

	if DirExist(pf) {
		return ""
	}

	var msgReturn string

	// defer wg.Done()
	// time.Sleep(time.Second)

	t1 := time.Now()
	f, err := os.Open(pf)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	defer f.Close()

	// size file...
	fi, err := f.Stat()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	//// Use bufio.NewReader to get a Reader.
	// ... Then use ioutil.ReadAll to read the entire content.
	reader := bufio.NewReader(f)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("readAll:", err)
		return ""
	}

	contentType, err := GetFileContentType(f)
	if err != nil {
		fmt.Println("contentType: ", err)
		return ""
	}

	if len(string(b)) == 0 {
		println("Error file está vazio..")
		return ""
	}

	bs := string(b)

	// runer
	// timer := RunerTimer()

	type Fs struct {
		Msgs3 *string
		Name  string
		Size  int64
	}

	//var msgs3 = make(chan *string)
	var cfs = make(chan Fs, 1)
	//var wg2 = &sync.WaitGroup{}
	//wg2.Add(1)
	//<-timer

	// aqui deveria ter um worker
	go func(pf, b, contentType string) {
		//println(pf)
		//defer wg2.Done()
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
		//defer close(cfs)
		time.Sleep(time.Millisecond * 1)
	}(pf, bs, contentType)
	//wg2.Wait()

	//////////////////////////
	csfS := <-cfs
	kb := (csfS.Size / 1024)
	t2 := time.Now()

	sendMsgCor := YellowCor("\r\033[?25h" + "[send success]")
	msgReturn = sendMsgCor + " count[" + strconv.Itoa(I) + "] Id[" + *csfS.Msgs3 + "] File[" + pbucket + "/" + csfS.Name + "] Size[" + strconv.FormatInt(kb, 10) + "Kb]" + "time[" + t2.Sub(t1).String() + "]"
	msgReturn += "\r\033[?25h\033[?25h"

	//fmt.Print("\r")
	//fmt.Print("\033[?25h")
	//fmt.Println("[send success] Id["+*csfS.Msgs3+"] File["+csfS.Name+"] Size[", kb, "Kb]", "time[", t2.Sub(t1), "]")
	//fmt.Print("\r")
	//fmt.Print("\033[?25h")
	//fmt.Print("\033[?25h")
	/////////////////////////////

	// done..
	return msgReturn
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
