#!/bin/bash
##### @jeffotoni
sudo mkdir -p /opt/dospace/v1/
sudo wget https://github.com/jeffotoni/s3godo/tree/master/spaces/v1/copyspace -P /opt/dospace/v1/
sudo ln -s /opt/dospace/v1/copysapce /usr/bin/copyspace
sudo chmod 775 /usr/bin/copyspace


echo "\033[0;33m######### Thanks for Download ##########\033[0m"
echo "copyspace -h"