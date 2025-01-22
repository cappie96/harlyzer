package main

import (
	"fmt"
	"os"

	"github.com/cap79/harlyzer/harlyzer"
	"github.com/cap79/harlyzer/service"
)

func main() {
	harFilePath := "/home/cap52/Development/harlyzer/www.wireshark.org.har"

	har, err := harlyzer.ParseHarFile(harFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing HAR file: %v\n", err)
	}

	//fmt.Println("HAR file Parsed Successfully")
	//fmt.Printf("Log Creator: %s v%s\n", har.Log.Creator.Name, har.Log.Creator.Version)
	//for _, entry := range har.Log.Entries {
	//	fmt.Printf("Request URL: %s, Method: %s, Status: %d\n", entry.Request.URL,
	//		entry.Request.Method, entry.Response.Status)
	//}

	emulator := os.Getenv("TERM")
	if emulator == "" {
		emulator = "xterm-256color"
		err := os.Setenv("TERM", emulator)
		if err != nil {
			fmt.Println("Error setting TERM environment variable:", err)
			return
		}
	}

	term := service.NewTerminal()
	term.Init()
	if err := term.Run(har); err != nil {
		fmt.Fprintf(os.Stderr, "application error: %v\n", err)
		os.Exit(1)
	}
}
