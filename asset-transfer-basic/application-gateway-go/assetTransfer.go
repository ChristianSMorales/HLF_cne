package main

import (
	"application-gateway-go/Assets/Structs"
	"application-gateway-go/Tools"
	CryptoMaterial "application-gateway-go/Tools/Crypto"
	"application-gateway-go/Tools/IPFS"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Configuramos el material criptografico del endorser peer con el cual Fabric gateway se va a comunicar,
// en este caso el peer0 de la organización CNE
var (
	root, _      = os.Getwd()
	mspID        = "CNEMSP"
	cryptoPath   = filepath.Join(root, "..", "..", "cne-network", "organizations", "peerOrganizations", "cne.com")
	certPath     = filepath.Join(cryptoPath, "users", "User1@cne.com", "msp", "signcerts")
	keyPath      = filepath.Join(cryptoPath, "users", "User1@cne.com", "msp", "keystore")
	tlsCertPath  = filepath.Join(cryptoPath, "peers", "peer0.cne.com", "tls", "ca.crt")
	peerEndpoint = "dns:///localhost:7051"
	gatewayPeer  = "peer0.cne.com"
)

func main() {

	clientConnection := Tools.NewGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)
	defer clientConnection.Close()
	id := Tools.NewIdentity(certPath, mspID)
	sign := Tools.NewSign(keyPath)

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	chaincode := "chaincode"

	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincode = ccname
	}
	channelName := "canal"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}
	//Establecemos un solo canal, en el cual se trabaja con un chaincode y tres smartcontracts destinados al: registro
	//de dispositivos scanners de los centros de digitalización de actas, registro de actas, registro de resultados
	//electorales.
	network := gw.GetNetwork(channelName)
	contractIdentidad := network.GetContractWithName(chaincode, "SmartContractIdentidad")
	contractActas := network.GetContractWithName(chaincode, "SmartContractActas")
	contractResultados := network.GetContractWithName(chaincode, "SmartContractResultados")
	// En base a los ficheros alojados en ./Assets/Dictionary creamos maps los cuales sirvan para consulta en base
	//a un código, el nombre de provincia, canto o el dicionario que corresponda.
	Tools.BuildDictionary()

	seccion1 := "-> makedirs <cda.txt-path>  Crea un sistema de carpetas que describe la estructura Provincia.Canton.Parroquia.Recinto en base al fichero cda.txt\n" +
		"-> regscanners <cdaParentFolder-path> Crea material criptográfico y registra en la cadena de bloques cada dispositivo descrito como cda. Si <cdaParentForlde-Path> es creado con Makedirs, no se requiere el path\n" +
		"-> worldstate Devuelve los registros actualmente almacenados\n"
	seccion2 := "-> registraracta <acta-path> Almacena en IPFS y registra en hyperledger fabric la transacción, el acta entregada mediante el path" +
		" cada dispositivo descrito como cda\n" +
		"-> consultaracta <CID> En base a un CID devuelve el acta asociada\n" +
		"-> historicoacta <CID> En base a un CID devuelve todas las versiones asociadas\n" +
		"-> eliminaracta <CID> En base a un CID, elimina el acta del worlstate\n" +
		"-> verresultados Muestra la tabla de resultados actual"

	advice := "=============================================================\n" +
		"Uso: FabricApp <command> [<args>]\n" +
		"Commands:\n" +
		"*******Identidad de cdas*******\n" +
		seccion1 +
		"*******Procesamiento de Actas*******\n" +
		seccion2 +
		"*******************************"

	if len(os.Args) < 2 {
		fmt.Println("=================== Inicio Fabric application-gateway-go V1.0===================")
		aviso := "Esta aplicación está destinada para registrar los dispositivos scanners de cada centro de digitalización electoral, a fin de preservar la trazabilidad de los datos receptados en el sistema de escrutinio"
		lines := Tools.SplitText(aviso, 50)
		for _, line := range lines {
			fmt.Println(line)
		}
		fmt.Println(advice)
		os.Exit(1)
	}
	command := os.Args[1]

	switch command {
	//makedirs crea la estructura de carpetas del tipo /Assets/cda/Provincia/canton/parroquia/recinto/junta
	//al final de este path se guarda el documento donde se especifica el ID del dispositivo emisor de actas.
	// en base a este dispositivo y sus metadatos se crea par de llavers RSA y un certificado digital.
	// IMPORTANTE: La fuente de cuantos directorios deben crearse es /Assets/cda.txt, en este documento se lista los
	//cdas que deberan tener scanners
	case "makedirs":
		var cdaFilePath string
		if len(os.Args) == 2 {
			cdaFilePath = filepath.Join(root, "Assets", "cda.txt")
		} else if len(os.Args) == 3 {
			cdaFilePath = os.Args[2]
		} else {
			fmt.Println("Uso:Makedirs <cda.txt-path> ")
			os.Exit(1)
		}
		if cdaFilePath != "" {
			Tools.MakeFolders(cdaFilePath)
			fmt.Println("¡Estructura de directorios creada satisfactoriamente!")
			fmt.Println("Directorio de salida: Assets/")
		}
		// regscanners registra en el ledger de identidad el cdaida,id del dispositivo (IP+MAC), cert. Sentando que estos
		//dispositivos son los registrados y autorizados a enviar actas.
	case "regscanners":
		var cdaPath string
		if len(os.Args) == 2 {
			cdaPath = filepath.Join(root, "Assets", "cda")
		} else if len(os.Args) == 3 {
			cdaPath = os.Args[2]
		} else {
			fmt.Println("Uso:RegScanners <cdaParentFolder-path>")
			os.Exit(1)
		}

		if cdaPath != "" {
			//Generamos un arreglo de objetos Recinto.
			cdas := Tools.GetCDAs(cdaPath)
			//Generamos el par de llaves y el certificado para cada CDA y al final todo el arreglo lo hacemos json
			cdasCryptoJSON := CryptoMaterial.GenCrytoJSON(cdas)
			//Init ledger es la única función de escritura en el ledger de identidad.
			result, err := contractIdentidad.SubmitTransaction("InitLedger", string(cdasCryptoJSON))
			if err != nil {
				fmt.Println(err)
			}
			//Iniciamos el ledger de resultados con valores nulos.
			_, err = contractResultados.SubmitTransaction("InitLedger")
			if err != nil {
				fmt.Println(err)
			}
			if result != nil {
				log.Fatalf("No se pudo enviar la transaccion: %v", err)
			} else {
				//generamos un arreglo en base a un objeto que almacene los metadatos y el cripto material de cada cda
				var cdasCrypto []Structs.CDAcrypto
				err = json.Unmarshal([]byte(cdasCryptoJSON), &cdasCrypto)
				if err != nil {
					fmt.Println(err)
				}
				//Almacenamos en el directorio respectivo el cripto materia (par de llaves RSA y el certificado)
				err = Tools.SaveCripto(cdasCrypto)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Material Cryptográfico almacenado")
			}
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Registro de transacción exitosa")

		}
		//worldstate retorna todos los registros de los cda almacenados en la red blockchain.
	case "worldstate":
		fmt.Println("\nRetornando todos la lista de activos registrados")
		evaluateResult, err := contractIdentidad.EvaluateTransaction("GetAllAssets")
		if err != nil {
			fmt.Println(err)
		}
		result := Tools.FormatJSON(evaluateResult)

		fmt.Println("**RESULT**\n" + result)
		//registraracta, a partir de las actas enviadas por un cda, almacena el documento en IPFS,
		//firma y envia a la red blockchain para almacenar el registrar los datos de integridad del acta
	case "registraracta":
		var actaPath string
		var txsIdentidad []Structs.TransactionIdentidad
		if len(os.Args) == 2 {
			//Assets/Actas/1.260.5140.1035.32.M
			actaPath = filepath.Join(root, "Assets", "Actas", "14.796.6870.22.1.F.tif")
		} else if len(os.Args) == 3 {
			actaPath = os.Args[2]
		} else {
			fmt.Println("Uso:ReadActas <Acta-PATH>")
			os.Exit(1)
		}
		if actaPath != "" {
			privkey, _, certificado, err := Tools.GetCryptoCDA(actaPath)
			//Cargamos a IPFS, firmamos y obtenemos, el CID desde IPFS, idDevice desde el cert,
			//filesigned (base64(hashfile+signature)) datos de integridad a ser registrados en la red blockchain
			cid, idDevice, filesigned, err := Tools.UploadIPFS(actaPath, privkey, certificado)
			if err != nil {
				fmt.Println("IPFS ERROR: %S", err)
			}
			//Consultamos el ledger de identidad por el id del dispositivo para corroborar que ha sido registrado.
			response, err := contractIdentidad.EvaluateTransaction("GetAssetHistory", idDevice)
			if err != nil {
				fmt.Println("El dispositivo no se encuentra en el ledger: %v", err)
			} else {
				//Obtenemos la última transacción respecto a ese dispositivo, si existe y no tiene la bandera
				//de eliminado, se procede con el cálculo de los resultados
				err = json.Unmarshal(response, &txsIdentidad)
				lastTx := txsIdentidad[len(txsIdentidad)-1]
				if !lastTx.IsDelete {
					_, err = contractActas.SubmitTransaction("CreateAsset", cid, idDevice, filesigned)
					//Debo obtener los datos del acta y realizar la tabla de resultados
					resultados := Tools.CalculateResults(filepath.Base(actaPath))
					if err != nil {
						fmt.Println("CC ERROR: ", err)
					}
					resultadosJSON, err := json.Marshal(resultados)
					_, err = contractResultados.SubmitTransaction("CreateAsset", string(resultadosJSON))
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Acta registrada con exito!")
				} else {
					fmt.Println("Dispositivo %s tiene estado deleted en ledger!", idDevice)
				}
			}

		}
		//Eliminamos del worldstate un acta, la acción se registra en el ledger, dejando un registro con el CID y la
		//bandera de eliminado en true
	case "eliminaracta":
		var cid string
		if len(os.Args) == 2 {
			cid = "QmQ92U5Ag3e2CjkdBuFxqRyZyRr3aA3MYLqiVKvsZMTVKW"
		} else if len(os.Args) == 3 {
			cid = os.Args[2]
		} else {
			fmt.Println("Uso:eliminaracta <cid>")
			os.Exit(1)
		}
		if cid != "" {
			err := IPFS.ExsitsFile(cid)
			if err != nil {
				fmt.Println("No se encontro el acta en IPFS")
			}
			_, err = contractActas.SubmitTransaction("DeleteAsset", cid)
			if err != nil {
				fmt.Println("Error en transaction, no se pudo eliminar el cid")
			}
			fmt.Println("Acta eliminada!")
		}
		//Retorna el worldstate del acta en base al ID dado como argumento.
	case "consultaracta":
		var cid string
		var acta Structs.AssetActa
		if len(os.Args) == 2 {
			cid = "QmQ92U5Ag3e2CjkdBuFxqRyZyRr3aA3MYLqiVKvsZMTVKW"
		} else if len(os.Args) == 3 {
			cid = os.Args[2]
		} else {
			fmt.Println("Uso:consultaracta <cid>")
			os.Exit(1)
		}
		if cid != "" {
			err := IPFS.ExsitsFile(cid)
			if err != nil {
				fmt.Println("No se encontro el acta en IPFS")
			}
			asset, err := contractActas.SubmitTransaction("ReadAsset", cid)
			err = json.Unmarshal(asset, &acta)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("**ACTA**\n" +
				"CID: " + acta.CID + "\n" +
				"IDdevice: " + acta.IDDEVICE + "\n" +
				"FILESIGNED(B64): " + acta.FILESIGNED + "\n")
			if err != nil {
				fmt.Println("Error en transaction, no se pudo obtener el asset")
			}
		}
		//Retorna el historial de transacciones de un acta, el argumento es el CID del acta.
	case "historicoacta":
		var cid string
		if len(os.Args) == 2 {
			cid = "QmQ92U5Ag3e2CjkdBuFxqRyZyRr3aA3MYLqiVKvsZMTVKW"
		} else if len(os.Args) == 3 {
			cid = os.Args[2]
		} else {
			fmt.Println("Uso:consultaracta <cid>")
			os.Exit(1)
		}
		if cid != "" {
			err := IPFS.ExsitsFile(cid)
			if err != nil {
				fmt.Println("No se encontro el acta en IPFS")
			}
			transactions, err := contractActas.EvaluateTransaction("GetAssetHistory", cid)
			result := Tools.FormatJSON(transactions)

			fmt.Println("**RESULT**\n" + result)
			if err != nil {
				fmt.Println("Error en transaction, no se pudo recuperar el historico en base al cid")
			}
		}
		//Retorna la tabla de resultados actual en el worldstate.
	case "verresultados":
		if len(os.Args) == 2 {
			transaction, err := contractResultados.EvaluateTransaction("ReadAsset")
			if err != nil {
				fmt.Println("Error en transaction, no se pudo recuperar el ledger de resultados")
			}
			result := Tools.FormatJSON(transaction)
			fmt.Println("**RESULT**\n" + result)

		}
	default:
		fmt.Println("¡COMANDO NO ENCONTRADO!")
		fmt.Println(advice)
		fmt.Println("=================== Fin Fabric application-gateway-go V1.0===================")

	}
}
