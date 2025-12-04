package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/skiupace/gitnav/app"
)

func main() {
	repoPath := "."

	if len(os.Args) > 1 && os.Args[1] != "." {
		repoPath = os.Args[1]
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		fmt.Println("Invalid path:", err)
		return
	}

	fmt.Println("Opening repo at:", absPath)

	if err := app.App.Run(absPath); err != nil {
		log.Fatal(err)
	}
}
