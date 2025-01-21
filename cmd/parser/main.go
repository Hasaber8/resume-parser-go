package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"resumeparser/internal/extractor"
	"resumeparser/internal/parser"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: Please provide a PDF file path\n")
		fmt.Fprintf(os.Stderr, "Usage: %s <pdf-file>\n", os.Args[0])
		os.Exit(1)
	}

	pdfPath := os.Args[1]
	// check if the the extension of the file is .pdf
	// we support only pdf atm
	if strings.Split(pdfPath, ".")[len(strings.Split(pdfPath, "."))-1] != "pdf" {
		fmt.Print("Error: Provided file is not a pdf\n")
		os.Exit(1)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ext := extractor.New()

	text, err := ext.Extract(ctx, pdfPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Add debug logging
	fmt.Fprintf(os.Stderr, "Reading PDF: %s\n", os.Args[1])

	fmt.Fprintf(os.Stderr, "Extracted text length: %d\n", len(text))
	if len(text) > 100 {
		fmt.Fprintf(os.Stderr, "First 100 chars: %q\n", text[:100])
	}

	// Create and configure parser
	p := parser.NewParser()

	// Parse the resume
	resume, err := p.Parse(text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing resume: %v\n", err)
		os.Exit(1)
	}

	// For now, just output as JSON
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(resume); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding result: %v\n", err)
		os.Exit(1)
	}
}
