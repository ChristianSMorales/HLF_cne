#!/bin/bash
# Git
sudo apt-get remove --auto-remove git -y
sudo apt-get install git

# Curl
sudo apt remove curl
sudo apt-get install curl

# Docker
sudo apt install docker-ce
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo usermod -aG docker $USER

# Instalar Go (1.22.3)
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
sudo cp ./fabric-samples/bin/*    /usr/local/bin# Git
sudo apt-get remove --auto-remove git -y
sudo apt-get install git

# Curl
sudo apt remove curl
sudo apt-get install curl

# Docker
sudo apt install docker-ce
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo usermod -aG docker $USER

# Instalar Go (1.22.3)
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
