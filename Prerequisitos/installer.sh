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
sudo rm go1.22.0.linux-amd64.tar.gz

# Instalar JQ (1.7)
sudo apt remove jq -y
sudo apt-get install jq -y


# Samples
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
./install-fabric.sh --fabric-version 2.4.0 docker samples binary
# Mover los samples
sudo cp ./fabric-samples/bin/*    /usr/local/bin

# Función para imprimir mensajes de error y salir
error_exit() {
    echo "$1" 1>&2
    exit 1
}

# Descargar Kubo (IPFS)
echo "Descargando Kubo (IPFS)..."
wget https://dist.ipfs.tech/kubo/v0.29.0/kubo_v0.29.0_linux-amd64.tar.gz -O kubo.tar.gz || error_exit "Error descargando Kubo."

# Descomprimir el archivo descargado
echo "Descomprimiendo Kubo..."
tar -xzf kubo.tar.gz || error_exit "Error descomprimiendo Kubo."

# Mover el binario a /usr/local/bin
echo "Instalando Kubo..."
sudo mv kubo/ipfs /usr/local/bin/ipfs || error_exit "Error moviendo Kubo a /usr/local/bin."

# Limpiar archivos temporales
sudo rm -rf kubo kubo.tar.gz

# Inicializar el nodo IPFS
echo "Inicializando el nodo IPFS..."
ipfs init || error_exit "Error inicializando el nodo IPFS."

# Configurar IPFS para arrancar automáticamente al iniciar el sistema (opcional)
echo "Configurando IPFS para iniciar automáticamente..."
ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["*"]' || error_exit "Error configurando CORS."
ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "GET", "POST"]' || error_exit "Error configurando métodos CORS."

# Mostrar información del nodo
ipfs id


