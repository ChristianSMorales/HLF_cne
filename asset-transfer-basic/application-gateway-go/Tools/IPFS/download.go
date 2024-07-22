package IPFS

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	shell "github.com/ipfs/go-ipfs-api"
)

// MIMETypeToExtension is a map that maps MIME types to file extensions
var MIMETypeToExtension = map[string]string{
	"application/pdf":              ".pdf",
	"image/jpeg":                   ".jpg",
	"image/png":                    ".png",
	"text/plain":                   ".txt",
	"application/msword":           ".doc",
	"application/vnd.ms-excel":     ".xls",
	"application/zip":              ".zip",
	"application/x-rar-compressed": ".rar",
	"application/octet-stream":     ".docx",
}

// DownloadFileFromCID downloads a file from IPFS using its CID and saves it locally with the appropriate extension
func DownloadFile(cid string, outputPath string) error {

	// Connect to the local IPFS node
	sh := shell.NewShell("localhost:5001")

	// Get the file from IPFS
	reader, err := sh.Cat(cid)
	if err != nil {
		return fmt.Errorf("failed to get file from IPFS: %v", err)
	}
	defer reader.Close()

	// Determine the file extension based on the MIME type
	mimeType, err := getFileMIMEType(reader)
	if err != nil {
		return fmt.Errorf("could not get mime type of file: %v", err)
	}

	// Get the MIME type extension
	ext, ok := MIMETypeToExtension[mimeType]
	if !ok {
		ext = ".bin"
	}

	// Create the output file with the correct extension
	outFile, err := os.Create(filepath.Join(outputPath, cid+ext))
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Copy content from IPFS to local file
	_, err = io.Copy(outFile, reader)
	if err != nil {
		return fmt.Errorf("could not copy content to output file: %v", err)
	}
	return nil
}

// getFileMIMEType gets the MIME type of the file based on its contents
func getFileMIMEType(reader io.Reader) (string, error) {
	buffer := make([]byte, 512)
	_, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("could not read the contents of the file: %v", err)
	}

	// Determine the MIME type
	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}
