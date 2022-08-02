#! /bin/sh

#Installing terraform 
sudo apt-get update && sudo apt-get install -y gnupg software-properties-common curl 
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
sudo apt-add-repository "deb [arch=$(dpkg --print-architecture)] https://apt.releases.hashicorp.com $(lsb_release -cs) main" 
sudo apt update
sudo apt install terraform

#downloading unzip to unzip tflint 
sudo apt install unzip

#Installing yor_trace 
brew tap bridgecrewio/tap
brew install bridgecrewio/tap/yor

#Installing tflint 
curl -s https://raw.githubusercontent.com/terraform-linters/tflint/master/install_linux.sh | bash


#Checking if terraform is installed
terraform 

Arch=$(uname -m)

if [[ "$Arch" == "x86_64" || "$Arch" == "amd64" ]]; then
    wget https://github.com/trilogy-group/cloudfix-linter/releases/latest/download/cloudfix-linter_linux_amd64
    mv cloudfix-linter_linux_amd64 cloudfixlinter
elif [[ "$Arch" == "aarch64" || "$Arch" == "arm64" ]]; then
    wget https://github.com/trilogy-group/cloudfix-linter/releases/latest/download/cloudfix-linter_linux_arm64
    mv cloudfix-linter_linux_arm64 cloudfixlinter
elif [[ "$Arch" == "i686" || "$Arch" == "i386" ]]; then
    wget https://github.com/trilogy-group/cloudfix-linter/releases/latest/download/cloudfix-linter_linux_386
    mv cloudfix-linter_linux_386 cloudfixlinter
elif [ "$Arch" = "armhf" ]; then
    wget https://github.com/trilogy-group/cloudfix-linter/releases/latest/download/cloudfix-linter_linux_arm
    mv cloudfix-linter_linux_arm cloudfixlinter
fi

sudo mv cloudfixlinter  /usr/local/bin/
sudo chown root:root /usr/local/bin/cloudfixlinter
sudo chmod 755 /usr/local/bin/cloudfixlinter