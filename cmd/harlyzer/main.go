package main

import (
	"fmt"
	"os"

	"github.com/cap79/harlyzer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: harlyzer <har-file>")
		return
	}
	harFilePath := os.Args[1]

	har, err := harlyzer.ParseHarFile(harFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing HAR file: %v\n", err)
	}

	emulator := os.Getenv("TERM")
	if emulator == "" {
		emulator = "xterm-256color"
		err := os.Setenv("TERM", emulator)
		if err != nil {
			fmt.Println("Error setting TERM environment variable:", err)
			return
		}
	}

	term := harlyzer.NewTerminal()
	term.Init()
	if err := term.Run(har); err != nil {
		fmt.Fprintf(os.Stderr, "application error: %v\n", err)
		os.Exit(1)
	}
}
