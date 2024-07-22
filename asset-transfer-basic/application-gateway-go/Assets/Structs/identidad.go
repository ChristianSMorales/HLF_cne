package Structs

// ID:IP+MAC, CDA:CODPROVINCIA-CANTON....CDA, CERT
type AssetIdentidad struct {
	ID   string `json:"IdDevice"`
	CDA  string `json:"CDA"`
	Cert string `json:"Cert"`
}
