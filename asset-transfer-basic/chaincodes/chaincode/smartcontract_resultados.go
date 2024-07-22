package chaincode

import (
	"chaincodes/Assets/Structs"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContractResultados struct {
	contractapi.Contract
}

func (s *SmartContractResultados) InitLedger(ctx contractapi.TransactionContextInterface) error {
	resultados := Structs.Resultados{
		VotosCandidatoA: 0,
		VotosCandidatoB: 0,
		TotalVotantes:   0,
		VotosNulos:      0,
		VotosBlancos:    0,
		IdActa:          "",
	}
	assetJSON, err := json.Marshal(resultados)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("1", assetJSON)
}
func (s *SmartContractResultados) CreateAsset(ctx contractapi.TransactionContextInterface, resultadosJSON string) error {
	var resultados Structs.Resultados
	err := json.Unmarshal([]byte(resultadosJSON), &resultados)
	if err != nil {
		return err
	}
	response, err := s.ReadAsset(ctx)
	if err != nil {
		return err
	}

	asset := Structs.Resultados{
		VotosCandidatoA: response.VotosCandidatoA + resultados.VotosCandidatoA,
		VotosCandidatoB: response.VotosCandidatoB + resultados.VotosCandidatoB,
		TotalVotantes:   response.TotalVotantes + resultados.TotalVotantes,
		VotosNulos:      response.VotosNulos + resultados.VotosNulos,
		VotosBlancos:    response.VotosBlancos + resultados.VotosBlancos,
		IdActa:          resultados.IdActa,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("1", assetJSON)
}
func (s *SmartContractResultados) ReadAsset(ctx contractapi.TransactionContextInterface) (*Structs.Resultados, error) {
	assetJSON, err := ctx.GetStub().GetState("1")
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", "1")
	}

	var asset Structs.Resultados
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}
