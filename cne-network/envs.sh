#!/bin/bash

# Configuración de variables de entorno para Hyperledger Fabric

export CURRENT_PATH=${PWD}

#Compose DIR
export DOCKER_COMPOSE_DIR=$PWD/../explorer/


#ENABLE TLS
export CORE_PEER_TLS_ENABLED=true

# MSP ID de la organización
export CORE_PEER_LOCALMSPID="CNEMSP"

# Certificado raíz TLS del peer
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/cne.com/peers/peer0.cne.com/tls/ca.crt

# Ruta del MSP del usuario admin
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/cne.com/users/Admin@cne.com/msp

# Dirección del peer
export CORE_PEER_ADDRESS=localhost:7051

# Variables adicionales
export FABRIC_CFG_PATH=${PWD}/../config

#Chaincode path

echo "Variables de entorno configuradas."
