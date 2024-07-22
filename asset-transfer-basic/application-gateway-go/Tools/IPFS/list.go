package IPFS

import (
	"context"
	"fmt"
	"path/filepath"

	shell "github.com/ipfs/go-ipfs-api"
)

// ListCID lists all CIDs within a given MFS path
func ListCID(mfsPath string) ([]string, error) {

	// Convert file paths from (/) to (\)
	mfsPath = filepath.ToSlash(mfsPath)

	// Connect to the local IPFS node
	sh := shell.NewShell("localhost:5001")

	// Verify the connection
	if !sh.IsUp() {
		return nil, fmt.Errorf("IPFS node is not running")
	}

	// Create a context
	ctx := context.Background()

	// List the directory in MFS
	entries, err := sh.FilesLs(ctx, mfsPath)
	if err != nil {
		return nil, fmt.Errorf("could not list directory in MFS: %v", err)
	}

	// Extract CIDs from the directory entries
	var cids []string
	for _, entry := range entries {
		entryPath := filepath.ToSlash(filepath.Join(mfsPath, entry.Name))
		stat, err := sh.FilesStat(ctx, entryPath)
		if err != nil {
			return nil, fmt.Errorf("could not stat file %s: %v", entry.Name, err)
		}
		fmt.Printf("%s, CID: %s\n", entry.Name, stat.Hash)
		cids = append(cids, stat.Hash)
	}

	return cids, nil
}
