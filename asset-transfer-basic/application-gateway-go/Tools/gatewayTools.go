package Tools

import (
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"path"

	"github.com/hyperledger/fabric-gateway/pkg/identity"
)

// Generamos una nueva coneción gRPC usando el material criptográfico del endorser peer
func NewGrpcConnection(tlsCertPath string, gatewayPeer string, peerEndpoint string) *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		fmt.Println(err)
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("Failed to create gRPC connection: %w", err))
	}

	return connection
}

// Dado que los certificados se encuentran almacenados como ficheros en formato PEM,
// Creamos una función que cargue ese fichero y retorne un objeto del tipo certificado
func loadCertificate(filename string) (*x509.Certificate, error) {

	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// Devolvemos el primer fichero almacenado dentro de un directorio, el formato de retorno es un arreglo de bytes
// Dado que los directorios donde se almacena el materia criptogáfico almacena solo un archivo, no existe inconveniente.
func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

// Creamos una identidad x509 que represente el endorserpeer, para ello usamos el certificado y el ID del MSP
func NewIdentity(certPath string, mspID string) *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(err)
	}
	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}
	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}
	return id
}

// Creamos una firma con la cual identificar y firmar los paquetes enviados por el gateway hacia el endorser peer.
func NewSign(keypath string) identity.Sign {
	privatekeyPEM, err := readFirstFile(keypath)
	if err != nil {
		panic(err)
	}
	privatekey, err := identity.PrivateKeyFromPEM(privatekeyPEM)
	if err != nil {
		panic(err)
	}
	sign, err := identity.NewPrivateKeySign(privatekey)
	if err != nil {
		panic(err)
	}
	return sign
}
