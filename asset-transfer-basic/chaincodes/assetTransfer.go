package main

import (
	"chaincodes/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"log"
)

func main() {

	identidad := new(chaincode.SmartContractIdentidad)
	actas := new(chaincode.SmartContractActas)
	resultados := new(chaincode.SmartContractResultados)

	cneChaincode, err := contractapi.NewChaincode(identidad, actas, resultados)
	if err != nil {
		log.Println("Error creating actas chaincode: %v", err)
	}

	if err := cneChaincode.Start(); err != nil {
		log.Println("Error starting actas chaincode: %v", err)
	}

}
