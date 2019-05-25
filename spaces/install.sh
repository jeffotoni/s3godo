#!/bin/bash
sudo mkdir /opt/c0d3s-generator/v1/
sudo wget https://github.com/RafaelGomides/c0d3s-generator/releases/download/v1.0/c0d3s-generator -P /opt/c0d3s-generator/v1/
sudo ln -s /opt/c0d3s-generator/v1/c0d3s-generator /usr/bin/codegen
sudo chmod 775 /usr/bin/codegen

echo 'Obrigado por baixar! Para usar, basta executar: codegen -h'
echo 'Thanks for download! Use command: codegen -h'