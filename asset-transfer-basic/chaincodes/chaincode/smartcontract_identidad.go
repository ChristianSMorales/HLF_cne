package chaincode

import (
	"chaincodes/Assets/Structs"
	"chaincodes/Tools"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"time"
)

type SmartContractIdentidad struct {
	contractapi.Contract
}

func (s *SmartContractIdentidad) InitLedger(ctx contractapi.TransactionContextInterface, cdasCryptoJSON string) error {
	var assets = Tools.CreateAsset(cdasCryptoJSON)

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return err
		}
	}
	return nil

}

func (s *SmartContractIdentidad) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Structs.AssetIdentidad, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Structs.AssetIdentidad
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Structs.AssetIdentidad
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (s *SmartContractIdentidad) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]Structs.TransactionIdentidad, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var transactions []Structs.TransactionIdentidad
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		fmt.Println(response)
		if err != nil {
			return nil, err
		}
		var asset Structs.AssetIdentidad
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Structs.AssetIdentidad{
				ID:   assetID,
				CDA:  "",
				Cert: "",
			}
		}
		//timestamp := timestamppb.Timestamp(&response.Timestamp).AsTime()
		tx := Structs.TransactionIdentidad{
			Record:    &asset,
			TxId:      response.TxId,
			Timestamp: time.Now(),
			IsDelete:  response.IsDelete,
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

/*
package chaincode

import (
	"chaincode_identidad-go/Assets/Structs"
	"chaincode_identidad-go/Tools"
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContractIdentidad struct {
	contractapi.Contract
}

type Asset struct {
	AppraisedValue int    `json:"AppraisedValue"`
	Color          string `json:"Color"`
	ID             string `json:"ID"`
	Owner          string `json:"Owner"`
	Size           int    `json:"Size"`
}

func (s *SmartContractIdentidad) InitLedger(ctx contractapi.TransactionContextInterface, cdasCryptoJSON string) error {

	str := Tools.Idle()
	var assets = Tools.CreateAsset(str)

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return err
		}
	}
	return nil

}

func (s *SmartContractIdentidad) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]Structs.CDAcrypto, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var cdasCrypto []Structs.CDAcrypto
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var cdaCrypto Structs.CDAcrypto
		err = json.Unmarshal(queryResponse.Value, &cdaCrypto)
		if err != nil {
			return nil, err
		}
		cdasCrypto = append(cdasCrypto, cdaCrypto)
	}

	return cdasCrypto, nil
}
*/
