#!/bin/bash

export CURRENT_PATH=${PWD}

#variables de entorno
source ./envs.sh

#Creamos el canal y levantamos las autoridads
sudo ./network.sh up createChannel -c canal -ca
#Desplegamos el primer cc
sudo ./network.sh deployCC -ccn chaincode -ccp ../asset-transfer-basic/chaincodes/ -ccl go -c canal
KEYSTORE_PATH=$CURRENT_PATH/organizations/peerOrganizations/cne.com/users/Admin\@cne.com/msp/keystore

KEYSTORE_FILE=$(ls -t "$KEYSTORE_PATH"/*_sk | head -n 1)

export KEYSTORE_FILE="$KEYSTORE_FILE"
#Desplegamos el segundo cc
#sudo ./network.sh deployCC -ccn actas -ccp ../asset-transfer-basic/chaincodes/chaincode_actas/ -ccl go -c canal
