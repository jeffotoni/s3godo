package main

import (
	"encoding/json"

	"fmt"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
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

	// https: //s3wf.sfo2.digitaloceanspaces.com/tmp/s3densotouch.tar.bz2
	fmt.Println(len(os.Args), os.Args)

	if len(os.Args) < 4 {

		fmt.Println("Ex:")
		fmt.Println("$copyobject cp do://bucket/dir1/dir2/dir3/file.pdf .")
		fmt.Println("$copyobject cp do://bucket/dir1/dir2/dir3/file.pdf /tmp/file.pdf")
		return
	}

	if os.Args[1] != "cp" {
		fmt.Println("Erro comando não permitido!")
		return
	}

	if len(os.Args[2]) <= 0 {
		fmt.Println("Erro, a origem é obrigatório")
		return
	}

	param2 := os.Args[2]
	vt := strings.Split(param2, ":")
	fmt.Println(vt, len(vt))

	if len(vt) <= 0 || vt[0] != "do" || len(vt[1]) == 0 {
		fmt.Println("Erro, o padrao pode esta errado, confira o exemplo:")
		fmt.Println("$ copyobject do://bucket/dir/file.pdf")
		return
	}

	vt2 := strings.Split(vt[1], "/")
	if len(vt2[2]) <= 0 {
		fmt.Println("Erro, o nome do bucket está vazio..")
		return
	}

	BUCKET = vt2[2]
	fmt.Println("bucket:: ", vt2[2], " len: ", len(vt2[3:]))
	lenght := len(vt2[3:])

	dest := make([]string, lenght)
	copy(dest, vt2[3:])
	//fmt.Println("caminho: ", dest)
	item := strings.Join(dest, "/")

	if len(os.Args[3]) <= 0 {
		fmt.Println("Erro, o caminho de destino é obrigatorio!")
		return
	}

	destino := os.Args[3]
	if destino == "." {
		destino = strings.TrimLeft(item, "/")
	} else {
		destino += "/" + item
	}

	//fmt.Println(len(dest) - 1)
	lenght = len(dest[:len(dest)-1])
	pathMkdir := strings.Join(dest[0:lenght], "/")
	fmt.Println("destino local:: ", pathMkdir)
	fmt.Println("key bucket:: ", destino)

	if err := os.MkdirAll(pathMkdir, 0755); err != nil {
		fmt.Println("Erro ao criar diretorio: ", err)
		return
	}

	// config
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key.Key, key.Secret, ""),
		Endpoint:    aws.String(key.Endpoint),
		Region:      aws.String(key.Region), // This is counter intuitive, but it will fail with a non-AWS region name.
	}

	newSession := session.New(s3Config)
	downloader := s3manager.NewDownloader(newSession)

	file, err := os.Create(item)
	if err != nil {
		fmt.Printf("Unable to open file %q, %v\n", err)
	}

	defer file.Close()

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(BUCKET),
			Key:    aws.String(item),
		})
	if err != nil {
		fmt.Printf("Unable to download item %q, %v", item, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
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
