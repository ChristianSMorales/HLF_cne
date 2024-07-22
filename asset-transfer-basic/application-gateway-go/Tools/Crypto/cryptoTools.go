package CryptoMaterial

import (
	"application-gateway-go/Assets/Structs"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

// Generates a pub and priv key pair
func genKeys() (*rsa.PrivateKey, []byte, error) {
	// Generar clave privada RSA
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("error al generar la clave privada: %v", err)
	}

	// Generar clave pública
	publicKey := privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("error al codificar la clave pública: %v", err)
	}

	// Codificar la clave pública en formato PEM
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return privateKey, publicKeyPEM, nil
}

// Generates a cert for a cda and writes cert, key pair files in a pathFolder
func genCrypto(cda Structs.Recinto) (*rsa.PrivateKey, []byte, []byte, error) {
	privKey, pubKeyPem, _ := genKeys()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            []string{"Ecuador"},
			Organization:       []string{"Concejo Nacional Electoral"},
			OrganizationalUnit: []string{cda.Recinto},
			Locality:           []string{cda.Canton},
			Province:           []string{cda.Provincia},
			StreetAddress:      []string{cda.Parroquia},
			CommonName:         cda.IdDevice,
			SerialNumber:       cda.CDAID,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // válido por 1 año
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error al crear la certificato: %v", err)
	}

	certPEMBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})

	return privKey, pubKeyPem, certPEMBytes, err
}
func GenCrytoJSON(cdas []Structs.Recinto) string {
	var cdascrypto []Structs.CDAcrypto

	for _, cda_ := range cdas {
		privkey, pubkey, cert, err := genCrypto(cda_)
		if err != nil {
			fmt.Println(err)
		}
		var cdacrypto = Structs.CDAcrypto{
			CDA:     cda_,
			PRIVKEY: privkey,
			PUBKEY:  pubkey,
			CERT:    cert,
		}
		cdascrypto = append(cdascrypto, cdacrypto)
	}
	cdasCryptoJSON, err := json.Marshal(cdascrypto)
	if err != nil {
		fmt.Println(err)
	}
	return string(cdasCryptoJSON)
}
