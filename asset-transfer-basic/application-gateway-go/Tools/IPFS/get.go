package IPFS

import (
	shell "github.com/ipfs/go-ipfs-api"
)

func ExsitsFile(cid string) error {
	// Connect to the local IPFS node
	sh := shell.NewShell("localhost:5001")

	// Get the file from IPFS
	reader, err := sh.Cat(cid)
	if err != nil {
		return err
	}
	defer reader.Close()

	return nil
}
