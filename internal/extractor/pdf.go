package extractor

import (
	"bytes"
	"context"
	"fmt"
	"os"

	_ "embed"
	"io/ioutil"
	"os/exec"
)

type PdfExtractor interface {
	Extract(ctx context.Context, path string) (string, error)
}

type pdfExtractor struct {
	jarPath string
}

func New() PdfExtractor {
	// Look for PDFBox JAR in standard locations
	jarPath := "assets/pdfbox-app-3.0.3.jar"
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		// Try alternate location
		jarPath = "internal/extractor/assets/pdfbox-app-3.0.3.jar"
	}
	return &pdfExtractor{
		jarPath: jarPath,
	}
}

func (e *pdfExtractor) Extract(ctx context.Context, path string) (string, error) {
	// Create a temporary output file for text extraction
	outputFile, err := os.CreateTemp("", "pdf-extract-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp output file: %w", err)
	}
	defer os.Remove(outputFile.Name())
	outputFile.Close()

	// Prepare the Java command to extract text with comprehensive options
	cmd := exec.CommandContext(ctx,
		"java",
		"-jar", e.jarPath,
		"export:text",
		"-i="+path,
		"-o="+outputFile.Name(),
		"-encoding=UTF-8",
		"-sort",
		"-startPage=1",
	)

	// Capture potential error output
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command with verbose logging
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("PDF extraction failed: %w\nDetailed error: %s\nCommand details: %v",
			err,
			stderr.String(),
			cmd.Args,
		)
	}

	// Read the extracted text
	textBytes, err := ioutil.ReadFile(outputFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read extracted text: %w", err)
	}

	return string(textBytes), nil
}
