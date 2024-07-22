package Structs

import (
	"crypto/rsa"
)

type CDAcrypto struct {
	CDA     Recinto         `json:"CDA"`
	PRIVKEY *rsa.PrivateKey `json:"PRIVKEY"`
	PUBKEY  []byte          `json:"PUBKEY"`
	CERT    []byte          `json:"CERT"`
}
