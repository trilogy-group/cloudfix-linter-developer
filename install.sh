#! /bin/bash


Arch=$(uname -m)

if [[ "$Arch" == "x86_64" || "$Arch" == "amd64" ]]; then
    PLATFORM=linux_amd64
elif [[ "$Arch" == "aarch64" || "$Arch" == "arm64" ]]; then
    PLATFORM=linux_arm64
elif [[ "$Arch" == "i686" || "$Arch" == "i386" ]]; then
    PLATFORM=linux_386
elif [ "$Arch" = "armhf" ]; then
    PLATFORM=linux_arm
else 
    echo "Unsupported platform"
    exit 1
fi

#Installing terraform 
TERRAFORM_VERSION=1.2.6
(wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_${PLATFORM}.zip \
  && unzip terraform_${TERRAFORM_VERSION}_${PLATFORM}.zip \
  && sudo mv terraform /usr/bin \
  && rm terraform_${TERRAFORM_VERSION}_${PLATFORM}.zip) || exit 1

#Checking if terraform is installed
terraform 

#downloading unzip to unzip tflint 
sudo apt install unzip

#Installing yor_trace 
YOR_VERSION=0.1.150
wget -q -O - https://github.com/bridgecrewio/yor/releases/download/${YOR_VERSION}/yor_${YOR_VERSION}_${PLATFORM}.tar.gz | tar -xvz -C /tmp               
sudo mv /tmp/yor /usr/local/bin/yor

#Installing tflint 
# higher version have breaking changes to the plugin system and hence we can't install them without changing the plugin
export TFLINT_VERSION=v0.39.3
curl -s https://raw.githubusercontent.com/terraform-linters/tflint/master/install_linux.sh | bash

# Install cloudfix-linter
echo "Installing cloudfix-linter"
(wget https://github.com/prasheel-ti/cloudfix-linter-developer/releases/latest/download/cloudfix-linter_${PLATFORM} \
  && mv cloudfix-linter_${PLATFORM} cloudfix-linter ) || exit 1
sudo mv cloudfix-linter  /usr/local/bin/
sudo chown root:root /usr/local/bin/cloudfix-linter
sudo chmod 755 /usr/local/bin/cloudfix-linter
