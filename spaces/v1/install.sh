#!/bin/bash
##### @jeffotoni
sudo rm -rf /opt/dospace/v1
sudo mkdir -p /opt/dospace/v1
sudo wget https://github.com/jeffotoni/s3godo/tree/master/spaces/v1/copyspace -P /opt/dospace/v1/
sudo chmod 755 -R /opt/dospace
sudo ln -s /opt/dospace/v1/copyspace /usr/bin/copyspace
echo "\033[0;33m######### Thanks for Download ##########\033[0m"
echo "comand: copyspace -h"
copyspace