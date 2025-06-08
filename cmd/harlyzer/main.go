package main

import (
	"fmt"
	"os"

	"github.com/cappstr/harlyzer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: harlyzer <har-file>")
		return
	}
	harFilePath := os.Args[1]

	har, err := harlyzer.ParseHarFile(harFilePath)
	if err != nil {
		fmt.Printf("error parsing HAR file: %v\n", err)
	}

	term := harlyzer.NewTerminal()
	term.Init()
	if err := term.Run(har); err != nil {
		fmt.Printf("application error: %v\n", err)
		os.Exit(1)
	}
}
