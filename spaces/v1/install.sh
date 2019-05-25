#!/bin/bash
##### @jeffotoni
DIR=/opt/dospace/v1
DIR2=/opt/dospace
EXEC=copyspace

echo "{
     "key": "key-digitalocean",
     "secret": "secret-digitalocean",
     "endpoint": "https://your-space.digitaloceanspaces.com",
     "region": "us-east-1",
     "bucket": "your-bucket-default"
}" > $HOME/.dokeys

sudo rm -rf $DIR
sudo mkdir -p $DIR
sudo wget -c "https://copyspace.sfo2.digitaloceanspaces.com/v1/copyspace" -P "$DIR"
echo "..."
sleep 1
sudo chmod 755 -R $DIR2
sudo rm -f /usr/bin/$EXEC
sudo ln -s $DIR/$EXEC /usr/bin/$EXEC

echo "\033[0;33m######### Thanks for Download ##########\033[0m"
echo "\033[0;33m You just need to configure your ~/.dokeys file \033[0m"
echo "comand: $EXEC -h"
echo "
  -acl string
    	permissao: public or private
  -bucket string
    	o nome do seu bucket
  -file string
    	nome do arquivo ou diretorio a ser enviado
  -worker string
    	quantidade de trabalhos concorrentes em sua máquina
        "