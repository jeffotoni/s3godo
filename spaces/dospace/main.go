package main

import (
	//"context"
	"encoding/json"
	//"flag"
	"fmt"
	"io/ioutil"
	//"net/http"
	"os"
	"os/user"
	//"path/filepath"
	//"strings"
	//"sync"
	//"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	//"github.com/jeffotoni/gcolor"
)

var (
	BUCKET   = ""
	WORKER   = 500 // quantidade de workers trabalhando simultaneamente
	ACL_S3   = "private"
	HOME_DIR = ""
)

type sendS3 struct {
	Path     string
	Pbucket  string
	S3Client *s3.S3
	Counter  int
}

// DOKey contem dados para autenticacao na Digital Ocean(acho).
type DOKey struct {
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	Endpoint string `json:"endpoint"`
	Region   string `json:"region"`
	Bucket   string `json:"bucket"`
}

var Comands = map[string]string{
	"ls":    "",
	"do://": "",
	"cp":    "",
	"rm":    "",
	"",
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

	// Initialize a client using Spaces
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key.Key, key.Secret, ""),
		Endpoint:    aws.String(key.Endpoint),
		Region:      aws.String(key.Region), // This is counter intuitive, but it will fail with a non-AWS region name.
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	// agora capturando dados..
	// var pathFile string
	// var workers int
	// flag.StringVar(&pathFile, "file", "", "nome do arquivo ou diretorio a ser enviado")
	// aclSend := flag.String("acl", "private", "permissao: public or private")
	// fbucket := flag.String("bucket", "", "o nome do seu bucket")
	// flag.IntVar(&workers, "worker", WORKER, "quantidade de trabalhos concorrentes em sua máquina")
	// ls := flag.String("ls", "", "listar bucket")
	// flag.Parse()

	fmt.Println(s3Client)
	//fmt.Println(aclSend)
	//fmt.Println(fbucket)
	//fmt.Println(*ls)

	fmt.Println(os.Args, len(os.Args))

	fmt.Println(os.Args[1])
	fmt.Println(os.Args[2])
	return
}

func ReadKey() (*DOKey, error) {
	// user, err := user.Current()
	// if err != nil {
	// 	return nil, err
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
