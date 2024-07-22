package Structs

// filesigner -> base64(filehash+sign)
type AssetActa struct {
	CID        string `json:"CID"`
	IDDEVICE   string `json:"IdDevice"`
	FILESIGNED string `json:"FileSigned"`
}
