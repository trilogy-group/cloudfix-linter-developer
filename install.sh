#! /bin/bash

if [[ "$OSTYPE" =~ ^darwin ]]; then
  OS=darwin
  brew install wget
else
  OS=linux
fi

Arch=$(uname -m)

if [[ "$Arch" == "x86_64" || "$Arch" == "amd64" ]]; then
    ARCH=amd64
elif [[ "$Arch" == "aarch64" || "$Arch" == "arm64" ]]; then
    ARCH=arm64
elif [[ "$Arch" == "i686" || "$Arch" == "i386" ]]; then
    ARCH=386
elif [ "$Arch" = "armhf" ]; then
    ARCH=arm
else 
    echo "Unsupported platform"
    exit 1
fi
PLATFORM=$OS
PLATFORM+="_"
PLATFORM+=$Arch

rm -r cloudfix-linter
mkdir cloudfix-linter
cd cloudfix-linter

#Installing terraform 
TERRAFORM_VERSION=1.2.6
( wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_${PLATFORM}.zip --no-check-certificate \
  && unzip terraform_${TERRAFORM_VERSION}_${PLATFORM}.zip \
  && rm terraform_${TERRAFORM_VERSION}_${PLATFORM}.zip)


#Installing yor_trace 
YOR_VERSION=0.1.150
wget https://github.com/bridgecrewio/yor/releases/download/${YOR_VERSION}/yor_${YOR_VERSION}_${PLATFORM}.tar.gz --no-check-certificate \
&& tar -xvzf yor_${YOR_VERSION}_${PLATFORM}.tar.gz \
&& rm yor_${YOR_VERSION}_${PLATFORM}.tar.gz               

#Installing tflint 
# higher version have breaking changes to the plugin system and hence we can't install them without changing the plugin
export TFLINT_VERSION=v0.39.3
(wget https://github.com/terraform-linters/tflint/releases/download/${TFLINT_VERSION}/tflint_${PLATFORM}.zip --no-check-certificate  \
  && unzip tflint_${PLATFORM}.zip \
  && rm tflint_${PLATFORM}.zip)

# Install cloudfix-linter
echo "Installing cloudfix-linter"
(wget https://github.com/trilogy-group/cloudfix-linter-developer/releases/latest/download/cloudfix-linter-developer_${PLATFORM} --no-check-certificate \
  && mv cloudfix-linter-developer_${PLATFORM} cloudfix-linter)
