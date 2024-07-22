package Tools

import (
	"application-gateway-go/Assets/Structs"
	"application-gateway-go/Tools/IPFS"
	"bufio"
	"bytes"
	"crypto"
	random "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/muonsoft/validation/validate"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// En esta sección se declara variables globales usadas para el retorno de funciones y la creación de diccionarios
var GRecinto []Structs.Recinto
var GCIDs []string
var provincias map[int]string
var cantones map[int]string
var parroquias map[int]string
var recintos map[int]string

// IsValidIP verifica si una dirección IP dada es IPv4
func IsValidIP(ip string) bool {
	ipState := validate.IP(ip, validate.DenyPrivateIP())
	// nil -> IP is public
	if ipState != nil {
		if strings.Contains(ipState.Error(), "prohibited") {
			return true
		}
	}
	return false
}

// IsVAlidMAC verifica si una dirección MAC dad tiene el formato correcto.
func IsValidMAC(mac string) bool {
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	//Check if format is right
	if !macRegex.MatchString(mac) {
		return false
	}
	_, macState := net.ParseMAC(mac)

	//Check if mac is valid
	if macState != nil {
		return false
	}
	return true
}

// Construye los diccionarios usados para traducir de texto a un identificado.
// Se utiliza la variable global GRencinto y los ficheros alojados en ./Assets/Dictionary
func BuildDictionary() {
	root, _ := os.Getwd()
	path := filepath.Join(root, "Assets", "Dictionary")
	provincias = make(map[int]string)
	cantones = make(map[int]string)
	parroquias = make(map[int]string)
	recintos = make(map[int]string)

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		archivo, error := os.Open(filepath.Join(path, e.Name()))

		if error != nil {
			log.Fatal(error)
		}
		switch e.Name() {
		case "provincias.txt":
			scanner := bufio.NewScanner(archivo)
			for scanner.Scan() {
				linea := scanner.Text()
				partes := strings.Split(linea, ": ")
				clave, _ := strconv.Atoi(partes[0])
				valor := partes[1]
				provincias[clave] = valor
			}
		case "cantones.txt":
			scanner := bufio.NewScanner(archivo)
			for scanner.Scan() {
				linea := scanner.Text()
				partes := strings.Split(linea, ": ")
				clave, _ := strconv.Atoi(partes[0])
				valor := partes[1]
				cantones[clave] = valor
			}
		case "parroquias.txt":
			scanner := bufio.NewScanner(archivo)
			for scanner.Scan() {
				linea := scanner.Text()
				partes := strings.Split(linea, ": ")
				clave, _ := strconv.Atoi(partes[0])
				valor := partes[1]
				parroquias[clave] = valor
			}
		case "recintoElectoral.txt":
			scanner := bufio.NewScanner(archivo)
			for scanner.Scan() {
				linea := scanner.Text()
				partes := strings.Split(linea, ": ")
				clave, _ := strconv.Atoi(partes[0])
				valor := partes[1]
				recintos[clave] = valor
			}
		}
	}
}

// GetCDA en base al fichero alojado en la estructura de carpetas retorna un arreglo de recintos.
func GetCDA(cdaFilePath string) (recintos []Structs.Recinto) {
	archivo, error := os.Open(cdaFilePath)
	if error != nil {
		log.Fatal(error)
	}
	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {
		linea := scanner.Text()
		partes := strings.Split(linea, ";")
		ip := partes[0]
		mac := partes[1]
		if IsValidIP(ip) && IsValidMAC(mac) {
			recintos = append(recintos, queryDictionaries(archivo.Name(), ip, mac))
		}
	}
	return recintos
}

// Parsing CDA name->Provincia-canton-parroquia-recinto and return an array of all recintos in cda folder
// GetCDAs realiza el parsing de todos los elementos dentro de la estructura de carpeta
// Se obtiene un arreglo de recintos. EL parametro de entrada es el path del directorio cda.
func GetCDAs(parentPath string) (recintos []Structs.Recinto) {
	err := filepath.Walk(parentPath, proccesFile)
	if err != nil {
		log.Fatal(err)
	}
	return GRecinto
}

// Función ancla usada recursivamente dentro de cada avance de directorio.
// La función busca el fichero con la información del la identificación del dispositvo.
func proccesFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() && !strings.HasSuffix(path, ".pem") {
		archivo, error := os.Open(path)
		if error != nil {
			log.Fatal(error)
		}
		scanner := bufio.NewScanner(archivo)
		for scanner.Scan() {
			linea := scanner.Text()
			partes := strings.Split(linea, ";")
			ip := partes[0]
			mac := partes[1]
			if IsValidIP(ip) && IsValidMAC(mac) {
				GRecinto = append(GRecinto, queryDictionaries(info.Name(), ip, mac))
			}

		}
	}
	return nil
}

// queryDictionaries en base al path del cdaid y su contenido se retorna un objeto del tipo Recinto
func queryDictionaries(filename string, ip string, mac string) (recinto Structs.Recinto) {
	filename = filepath.Base(filename)
	partes := strings.Split(filename, ".")
	iprovincia, _ := strconv.Atoi(partes[0])
	icanton, _ := strconv.Atoi(partes[1])
	iparroquia, _ := strconv.Atoi(partes[2])
	irecinto, _ := strconv.Atoi(partes[3])
	recinto = Structs.Recinto{
		IdDevice:  ip + ";" + mac,
		Provincia: provincias[iprovincia],
		Canton:    cantones[icanton],
		Parroquia: parroquias[iparroquia],
		Recinto:   recintos[irecinto],
		CDAID:     filename,
	}
	return recinto
}

// randomIP genera una IP aleatoria.
func randomIP() string {
	// Define los rangos de direcciones IP privadas
	ranges := []struct {
		start [4]byte
		end   [4]byte
	}{
		{start: [4]byte{10, 0, 0, 0}, end: [4]byte{10, 255, 255, 255}},
		{start: [4]byte{172, 16, 0, 0}, end: [4]byte{172, 31, 255, 255}},
		{start: [4]byte{192, 168, 0, 0}, end: [4]byte{192, 168, 255, 255}},
	}

	// Se coloca una categoria aleatoria
	r := ranges[rand.Intn(len(ranges))]

	ip := [4]byte{}
	for i := 0; i < 4; i++ {
		if r.end[i] > r.start[i] {
			ip[i] = r.start[i] + byte(rand.Intn(int(r.end[i]-r.start[i])+1))
		} else {
			ip[i] = r.start[i]
		}
	}

	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// MakeFolders crea la estructura de carpetas en base a ./Assets/cda.txt
func MakeFolders(docPath string) {
	archivo, error := os.Open(filepath.Join(docPath))
	if error != nil {
		log.Fatal(error)
	}
	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {
		root, _ := os.Getwd()
		path := filepath.Join(root, "Assets", "cda")
		linea := scanner.Text()
		if linea != "" {
			partes := strings.Split(linea, ".")
			for _, part := range partes {
				path = filepath.Join(path, part)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					err := os.MkdirAll(path, 0755)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			pathfile := filepath.Join(path, linea)
			//---
			file, err := os.OpenFile(pathfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				fmt.Printf("Error al abrir el archivo: %v\n", err)
				return
			}
			defer file.Close()

			// Escribir la línea al archivo
			if _, err := file.Write([]byte(randomIP() + ";20:20:20:20:20:20" + "\n")); err != nil {
				fmt.Printf("Error al escribir en el archivo: %v\n", err)
				return
			}
		}
	}
}

// SAveCripto almacena en la carpeta respectiva dentro de la estructura de carpetas,
// el cripto material de cada cda
func SaveCripto(cdaCrypto []Structs.CDAcrypto) error {
	root, _ := os.Getwd()
	for _, cda_ := range cdaCrypto {
		partes := strings.Split(cda_.CDA.CDAID, ".")
		path := filepath.Join(
			root,
			"Assets",
			"cda",
			partes[0],
			partes[1],
			partes[2],
			partes[3])
		createCryptoFiles(path, cda_.PRIVKEY, cda_.PUBKEY, cda_.CERT, cda_.CDA.IdDevice)
	}
	return nil
}

// createCryptoFIles Genera los ficheros del cripto material en un path dado
func createCryptoFiles(path string, privkey *rsa.PrivateKey, pubkey []byte, cert []byte, id_device string) {
	re := regexp.MustCompile(`[:.]`)
	id_device = re.ReplaceAllString(id_device, "_") // Reemplazar caracteres no válidos por guiones bajos

	privkeyPath := filepath.Join(path, id_device+"privkey.pem")
	privkeyFile, err := os.Create(privkeyPath)
	if err != nil {
		log.Fatal(err)
	}
	defer privkeyFile.Close()
	privkeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privkey),
	})
	privkeyFile.Write(privkeyPEM)

	pubkeyPath := filepath.Join(path, id_device+"pubkey.pem")
	pubkeyFile, err := os.Create(pubkeyPath)
	if err != nil {
		log.Fatal(err)
	}
	defer pubkeyFile.Close()
	pubkeyFile.Write(pubkey)

	certPath := filepath.Join(path, id_device+"cert.pem")
	certFile, err := os.Create(certPath)
	if err != nil {
		log.Fatal(err)
	}
	defer certFile.Close()
	certFile.Write(cert)
}

// Separa un texto dado en reglones de una misma longitud.
func SplitText(text string, maxLength int) []string {
	var lines []string

	// Dividir el texto en líneas de igual longitud
	for len(text) > maxLength {
		// Encontrar el índice para cortar
		idx := strings.LastIndex(text[:maxLength], " ")
		if idx <= 0 {
			idx = maxLength
		}

		// Agregar la línea cortada
		lines = append(lines, "||\t"+text[:idx])
		// Eliminar la parte añadida del texto
		text = text[idx+1:]
	}

	// Agregar la última línea
	if len(text) > 0 {
		lines = append(lines, "||\t"+text)
	}

	return lines
}

// Genera un json en base a un arreglo de bytes serializado.
func FormatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}

// SignFIle, firma un documento en base a una llave privada rsa dada. La función retorna
// la firma y el hash del fichero
func SignFIle(filePath string, privateKey *rsa.PrivateKey) ([]byte, [32]byte, error) {
	document, err := os.ReadFile(filePath)
	if err != nil {
		return nil, [32]byte{}, nil
	}
	hashfile := sha256.Sum256(document)
	signature, err := rsa.SignPKCS1v15(random.Reader, privateKey, crypto.SHA256, hashfile[:])
	if err != nil {
		return nil, [32]byte{}, nil
	}
	return signature, hashfile, nil
}

// Dado que la llave privada se encuentran almacenados como ficheros en formato PEM,
// Creamos una función que cargue ese fichero y retorne un objeto del tipo rsa.Privatekey
func loadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// Dado que las llaves públicas se encuentran almacenados como ficheros en formato PEM,
// Creamos una función que cargue ese fichero y retorne un objeto del tipo rsa.publickey
func loadPubKey(keyPath string) (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}

	return rsaPub, nil
}

// combineBytes combuna dos arreglos de bytes en un solo arreglo
func combineBytes(data []byte, hash [32]byte) []byte {
	combined := append(data, hash[:]...)
	return combined
}

// separateBytes separa los bytes en posisiciones estándar, de tal forma que se pueda obtener el hash, signature y filesigned
func separateBytes(combined []byte) ([]byte, [32]byte, error) {
	if len(combined) < 32 {
		return nil, [32]byte{}, fmt.Errorf("combined byte array is too short")
	}
	data := combined[:len(combined)-32]
	var hash [32]byte
	copy(hash[:], combined[len(combined)-32:])
	return data, hash, nil
}

// UploadIPFS firma el acta, carga en IPFS y retorna los metadastos necesarios para almacenar el registro en el Ledger
func UploadIPFS(actaPath string, privkey *rsa.PrivateKey, certificado *x509.Certificate) (string, string, string, error) {

	signature, hashfile, err := SignFIle(actaPath, privkey)
	if err != nil {
		return "", "", "", err
	}
	cid, err := IPFS.UploadFile(actaPath)
	if err != nil {
		return "", "", "", err
	}
	idDevice := certificado.Issuer.CommonName

	filesigned := base64.StdEncoding.EncodeToString(combineBytes(signature, hashfile))
	return cid, idDevice, filesigned, nil

}

// GetCryptoCDA en base a un acta, se busca en la estructura de carpetas el criptomaterial asociado al cda
// que registró esa acta. Em esta función se trata los datos de entrada y en la función encadenada se realiza
// la obtención del criptomaterial
func GetCryptoCDA(actaPath string) (*rsa.PrivateKey, *rsa.PublicKey, *x509.Certificate, error) {
	filename := filepath.Base(actaPath)
	partes := strings.Split(filename, ".")
	root, _ := os.Getwd()
	iprovincia, _ := strconv.Atoi(partes[0])
	icanton, _ := strconv.Atoi(partes[1])
	iparroquia, _ := strconv.Atoi(partes[2])
	irecinto, _ := strconv.Atoi(partes[3])
	workingPath := filepath.Join(root, "Assets", "cda", strconv.Itoa(iprovincia), strconv.Itoa(icanton), strconv.Itoa(iparroquia), strconv.Itoa(irecinto))
	cdaID := strings.Join([]string{strconv.Itoa(iprovincia), strconv.Itoa(icanton), strconv.Itoa(iparroquia), strconv.Itoa(irecinto)}, ".")
	privkey, pubkey, cert, err := getCrypto(workingPath, cdaID)
	if err != nil {
		return nil, nil, nil, err
	}
	return privkey, pubkey, cert, nil
}

// getCrypto carga el cryotomaterial de un path especifico
func getCrypto(workingPath string, cdaID string) (*rsa.PrivateKey, *rsa.PublicKey, *x509.Certificate, error) {
	cdaFilePath := filepath.Join(workingPath, cdaID)
	deviceid, err := getDeviceIDFromCDA(cdaFilePath)
	if err != nil {
		return nil, nil, nil, err
	}
	privKey, err := loadPrivateKey(filepath.Join(workingPath, deviceid+"privkey.pem"))
	if err != nil {
		return nil, nil, nil, err
	}
	pubKey, err := loadPubKey(filepath.Join(workingPath, deviceid+"pubkey.pem"))
	if err != nil {
		return nil, nil, nil, err
	}
	cert, err := loadCertificate(filepath.Join(workingPath, deviceid+"cert.pem"))
	if err != nil {
		return nil, nil, nil, err
	}
	return privKey, pubKey, cert, nil
}

// getDeviceIDfromCA obtiene el código de un cda mediante la lectura de lo
func getDeviceIDFromCDA(cdafilepath string) (string, error) {
	archivo, error := os.Open(cdafilepath)
	if error != nil {
		log.Fatal(error)
	}
	scanner := bufio.NewScanner(archivo)
	id_device := ""
	for scanner.Scan() {
		linea := scanner.Text()
		partes := strings.Split(linea, ";")
		ip := partes[0]
		mac := partes[1]
		id_device = ip + ";" + mac
	}
	re := regexp.MustCompile(`[:.]`)
	id_device = re.ReplaceAllString(id_device, "_") // Reemplazar caracteres no válidos por guiones bajos

	return id_device, nil
}

// Quema los resultados en base a una serie de actas establecidads
func CalculateResults(actaid string) Structs.Resultados {
	partes := strings.Split(actaid, ".")
	acta := partes[0] + "." + partes[1] + "." + partes[2] + "." + partes[3] + "." + partes[4] + "." + partes[5]
	/*
		1.260.320.23.8.M
		14.796.6870.22.1.F
		17.60.5220.120.6.M
		5.110.2490.1240.9.M
	*/
	var salida Structs.Resultados
	acta8m := Structs.Resultados{
		VotosCandidatoA: 170,
		VotosCandidatoB: 85,
		TotalVotantes:   276,
		VotosNulos:      20,
		VotosBlancos:    1,
		IdActa:          acta,
	}
	acta1f := Structs.Resultados{
		VotosCandidatoA: 31,
		VotosCandidatoB: 21,
		TotalVotantes:   54,
		VotosNulos:      2,
		VotosBlancos:    0,
		IdActa:          acta,
	}
	acta6m := Structs.Resultados{
		VotosCandidatoA: 161,
		VotosCandidatoB: 85,
		TotalVotantes:   264,
		VotosNulos:      16,
		VotosBlancos:    2,
		IdActa:          acta,
	}
	acta9m := Structs.Resultados{
		VotosCandidatoA: 67,
		VotosCandidatoB: 24,
		TotalVotantes:   101,
		VotosNulos:      8,
		VotosBlancos:    2,
		IdActa:          acta,
	}
	switch acta {
	case "1.260.320.23.8.M":
		salida = acta8m
	case "14.796.6870.22.1.F":
		salida = acta1f
	case "17.60.5220.120.6.M":
		salida = acta6m
	case "5.110.2490.1240.9.M":
		salida = acta9m
	}
	return salida
}
