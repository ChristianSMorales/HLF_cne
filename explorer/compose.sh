#!/bin/bash

#export KEYSTORE_FILE=$(sudo ls -t "$PWD"/../cne-network/organizations/peerOrganizations/cne.com/users/Admin\@cne.com/msp/keystore/*_sk | head -n 1)

sudo -E docker-compose up -d

echo "Explorador levantado con exito!"
