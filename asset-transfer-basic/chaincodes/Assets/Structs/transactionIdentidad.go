package Structs

import "time"

type TransactionIdentidad struct {
	Record    *AssetIdentidad `json:"record"`
	TxId      string          `json:"txId"`
	Timestamp time.Time       `json:"timestamp"`
	IsDelete  bool            `json:"isDelete"`
}
