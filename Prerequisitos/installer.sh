#!/bin/bash

# Git
sudo apt-get remove --auto-remove git -y
sudo apt-get install git -y

# Curl
sudo apt remove curl -y
sudo apt-get install curl -y

# Docker
sudo apt-get install ca-certificates -y
sudo install -m 0755 -d /etc/apt/keyrings -y
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update -y
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
sudo usermod -aG docker $USER

# Instalar Go (1.22.3)
wget https://storage.googleapis.com/golang/go1.22.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
sudo echo -e "\nexport GOROOT=/usr/local/go\nexport GOPATH=\$HOME/go\nexport PATH=\$GOPATH/bin:\$GOROOT/bin:\$PATH" >> ~/.profile
source ~/.profile

# Instalar JQ (1.7)
sudo apt remove jq -y
sudo apt-get install jq -y


# Samples
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
./install-fabric.sh --fabric-version 2.4.0 docker samples binary
# Mover los samples
sudo cp ./fabric-samples/bin/*    /usr/local/bin
