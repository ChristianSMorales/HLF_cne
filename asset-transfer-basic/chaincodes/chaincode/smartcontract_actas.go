package chaincode

import (
	"chaincodes/Assets/Structs"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"time"
)

type SmartContractActas struct {
	contractapi.Contract
}

func (s *SmartContractActas) InitLedger(ctx contractapi.TransactionContextInterface) error {
	asset := Structs.AssetActa{
		CID:        "000",
		IDDEVICE:   "000",
		FILESIGNED: "000",
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(asset.CID, assetJSON)
	return nil
}
func (s *SmartContractActas) CreateAsset(ctx contractapi.TransactionContextInterface, CID string, IdDevice string, FileSigned string) error {
	asset := Structs.AssetActa{
		CID:        CID,
		IDDEVICE:   IdDevice,
		FILESIGNED: FileSigned,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(CID, assetJSON)
}
func (s *SmartContractActas) ReadAsset(ctx contractapi.TransactionContextInterface, CID string) (*Structs.AssetActa, error) {
	assetJSON, err := ctx.GetStub().GetState(CID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", CID)
	}

	var asset Structs.AssetActa
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}
func (s *SmartContractActas) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Structs.AssetActa, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Structs.AssetActa
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Structs.AssetActa
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
func (s *SmartContractActas) AssetExists(ctx contractapi.TransactionContextInterface, CID string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(CID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}
func (s *SmartContractActas) DeleteAsset(ctx contractapi.TransactionContextInterface, CID string) error {
	exists, err := s.AssetExists(ctx, CID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Asset %s does not exist", CID)
	}
	return ctx.GetStub().DelState(CID)

}
func (s *SmartContractActas) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]Structs.TransactionActa, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var transactions []Structs.TransactionActa
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		fmt.Println(response)
		if err != nil {
			return nil, err
		}
		var asset Structs.AssetActa
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Structs.AssetActa{
				CID:        assetID,
				IDDEVICE:   "",
				FILESIGNED: "",
			}
		}
		//timestamp := timestamppb.Timestamp(&response.Timestamp).AsTime()
		tx := Structs.TransactionActa{
			Record:    &asset,
			TxId:      response.TxId,
			Timestamp: time.Now(),
			IsDelete:  response.IsDelete,
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}
