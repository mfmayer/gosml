package main

import (
	"bufio"
	"fmt"
	"os"

	sml "github.com/mfmayer/gosml"
)

func printUsage() {
	fmt.Printf("Usage: %s [FILE]...\n", os.Args[0])
	fmt.Println("  Reads FILE(s) and outputs found electricity meter readings (obis 1.8.0)")
}

func main() {
	// Check if at least one argument is given
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// define a handling callback function that get's called for matching obis list entries
	handleFunc := func(message *sml.ListEntry) {
		fmt.Printf("%s %f\n", message.ObjectName(), message.Float())
	}

	// Go through all arguments
	for _, filePath := range os.Args[1:] {
		// check if argument is a valid file path
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Error: File '%s' does not exist\n\n", filePath)
			printUsage()
			os.Exit(1)
		}
		// try to open the file
		f, err := os.OpenFile(filePath, os.O_RDONLY|256, 0666)
		if err != nil {
			panic(err)
		}
		// create a buffered reader for the file
		r := bufio.NewReader(f)
		// read the file using gosml module with the option
		// to call handleFunc for list entries that match the 1-1:1.8.0 obis code
		sml.Read(r, sml.WithObisCallback(sml.OctetString{1, 0, 1, 8, 0}, handleFunc))
	}
}
