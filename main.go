package main

import (
	"flag"
	"fmt"
	"morse_converter/morse"
	"os"
)

func main() {
	inputFile := flag.String("i", "", "Input file path (required)")
	outputFile := flag.String("o", "", "Output file path (optional, default: stdout)")
	help := flag.Bool("h", false, "Show help message")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *inputFile == "" {
		fmt.Println("Error: input file is required")
		flag.Usage()
		os.Exit(1)
	}

	converter, err := morse.New(*inputFile, *outputFile)
	if err != nil {
		fmt.Printf("Error: creating converter: %v\n", err)
		os.Exit(1)
	}
	defer converter.Close()

	if err := converter.Process(); err != nil {
		fmt.Printf("Error: processing conversion: %v\n", err)
		os.Exit(1)
	}
}
