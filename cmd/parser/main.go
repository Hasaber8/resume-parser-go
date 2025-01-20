package main

import (
	"context"
	"fmt"
	"os"
	"resumeparser/internal/extractor"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: Please provide a PDF file path\n")
		os.Exit(1)
	}

	pdfPath := os.Args[1]
	// check if the the extension of the file is .pdf
	// we support only pdf atm
	if strings.Split(pdfPath, ".")[len(strings.Split(pdfPath, "."))-1] != "pdf" {
		fmt.Print("Error: Provided file is not a pdf\n")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ext := extractor.New()

	text, err := ext.Extract(ctx, pdfPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(text)
}
