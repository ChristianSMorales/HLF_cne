package Structs

import "time"

type TransactionActa struct {
	Record    *AssetActa `json:"record"`
	TxId      string     `json:"txId"`
	Timestamp time.Time  `json:"timestamp"`
	IsDelete  bool       `json:"isDelete"`
}
