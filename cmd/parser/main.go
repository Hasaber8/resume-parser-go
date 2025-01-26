package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"resumeparser/internal/extractor"
	"resumeparser/internal/models"
	"resumeparser/internal/parser"
	"strings"
	"time"
)

func main() {
	// Command line flags
	debug := flag.Bool("debug", false, "Enable debug output")
	outputFormat := flag.String("format", "json", "Output format (json or text)")
	timeout := flag.Duration("timeout", 30*time.Second, "Processing timeout")
	flag.Parse()

	// Validate arguments
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: Please provide a PDF file path\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <pdf-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	pdfPath := flag.Arg(0)

	// Validate file
	if err := validateFile(pdfPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Initialize extractor
	if *debug {
		fmt.Fprintf(os.Stderr, "Initializing PDF extractor...\n")
	}
	ext := extractor.New()

	// Extract text
	if *debug {
		fmt.Fprintf(os.Stderr, "Processing PDF: %s\n", pdfPath)
	}
	text, err := ext.Extract(ctx, pdfPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting text: %v\n", err)
		os.Exit(1)
	}

	if *debug {
		fmt.Fprintf(os.Stderr, "Extracted %d characters of text\n", len(text))
		if len(text) > 100 {
			fmt.Fprintf(os.Stderr, "First 100 chars: %q\n", text[:100])
		}
	}

	// Create and configure parser
	if *debug {
		fmt.Fprintf(os.Stderr, "Initializing parser...\n")
	}
	p := parser.NewParser()

	// Parse the resume
	if *debug {
		fmt.Fprintf(os.Stderr, "Parsing resume content...\n")
	}
	resume, err := p.Parse(text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing resume: %v\n", err)
		os.Exit(1)
	}

	// Output results
	switch strings.ToLower(*outputFormat) {
	case "json":
		outputJSON(resume)
	case "text":
		outputText(resume)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown output format %q\n", *outputFormat)
		os.Exit(1)
	}
}

func validateFile(path string) error {
	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("error accessing file: %v", err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return fmt.Errorf("not a regular file: %s", path)
	}

	// Check extension
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".pdf" {
		return fmt.Errorf("unsupported file type %q (only PDF files are supported)", ext)
	}

	return nil
}

func outputJSON(resume *models.Resume) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(resume); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding result: %v\n", err)
		os.Exit(1)
	}
}

func outputText(resume *models.Resume) {
	// Print contact information
	if contact, ok := resume.Sections["contact"]; ok {
		if contactContent, ok := contact.Content.(*models.ContactContent); ok {
			fmt.Println("Contact Information:")
			if contactContent.Name != "" {
				fmt.Printf("  Name: %s\n", contactContent.Name)
			}
			if contactContent.Location != "" {
				fmt.Printf("  Location: %s\n", contactContent.Location)
			}
			if len(contactContent.Email) > 0 {
				fmt.Printf("  Email: %s\n", strings.Join(contactContent.Email, ", "))
			}
			if len(contactContent.Number) > 0 {
				fmt.Printf("  Phone: %s\n", strings.Join(contactContent.Number, ", "))
			}
			for platform, url := range contactContent.Social {
				fmt.Printf("  %s: %s\n", strings.Title(platform), url)
			}
			fmt.Println()
		}
	}

	// Print education
	if education, ok := resume.Sections["education"]; ok {
		if educationContent, ok := education.Content.(*models.TimelineContent); ok {
			fmt.Println("Education:")
			for _, entry := range educationContent.Entries {
				fmt.Printf("  %s", entry.Organization)
				if entry.Location != "" {
					fmt.Printf(", %s", entry.Location)
				}
				fmt.Println()

				if entry.Title != "" {
					fmt.Printf("  %s\n", entry.Title)
				}
				if entry.StartDate != "" || entry.EndDate != "" {
					fmt.Printf("  %s - %s\n", entry.StartDate, entry.EndDate)
				}
				for _, detail := range entry.Details {
					fmt.Printf("    • %s\n", detail)
				}
				fmt.Println()
			}
		}
	}

	// Print experience
	if experience, ok := resume.Sections["experience"]; ok {
		if experienceContent, ok := experience.Content.(*models.TimelineContent); ok {
			fmt.Println("Experience:")
			for _, entry := range experienceContent.Entries {
				fmt.Printf("  %s", entry.Organization)
				if entry.Location != "" {
					fmt.Printf(", %s", entry.Location)
				}
				fmt.Println()

				if entry.Title != "" {
					fmt.Printf("  %s\n", entry.Title)
				}
				if entry.StartDate != "" || entry.EndDate != "" {
					fmt.Printf("  %s - %s\n", entry.StartDate, entry.EndDate)
				}
				for _, detail := range entry.Details {
					fmt.Printf("    • %s\n", detail)
				}
				fmt.Println()
			}
		}
	}

	// Print projects
	if projects, ok := resume.Sections["projects"]; ok {
		if projectsContent, ok := projects.Content.(*models.TimelineContent); ok {
			fmt.Println("Projects:")
			for _, entry := range projectsContent.Entries {
				fmt.Printf("  %s", entry.Organization)
				if entry.Location != "" {
					fmt.Printf(", %s", entry.Location)
				}
				fmt.Println()

				if entry.Title != "" {
					fmt.Printf("  %s\n", entry.Title)
				}
				if entry.StartDate != "" || entry.EndDate != "" {
					fmt.Printf("  %s - %s\n", entry.StartDate, entry.EndDate)
				}
				for _, detail := range entry.Details {
					fmt.Printf("    • %s\n", detail)
				}
				fmt.Println()
			}
		}
	}

	// Print skills
	if skills, ok := resume.Sections["skills"]; ok {
		if skillsContent, ok := skills.Content.(*models.ListContent); ok {
			fmt.Println("Skills:")
			for _, category := range skillsContent.Categories {
				if category.Name != "" {
					fmt.Printf("  %s:\n", category.Name)
					for _, item := range category.Items {
						fmt.Printf("    • %s\n", item)
					}
				} else {
					for _, item := range category.Items {
						fmt.Printf("  • %s\n", item)
					}
				}
			}
			fmt.Println()
		}
	}

	// Print achievements
	if achievements, ok := resume.Sections["achievements"]; ok {
		if achievementsContent, ok := achievements.Content.(*models.ListContent); ok {
			fmt.Println("Achievements:")
			for _, category := range achievementsContent.Categories {
				if category.Name != "" {
					fmt.Printf("  %s:\n", category.Name)
					for _, item := range category.Items {
						fmt.Printf("    • %s\n", item)
					}
				} else {
					for _, item := range category.Items {
						fmt.Printf("  • %s\n", item)
					}
				}
			}
			fmt.Println()
		}
	}

	// Print freeform sections
	for name, section := range resume.Sections {
		if section.Type == models.FreeformSection {
			if freeformContent, ok := section.Content.(*models.FreeformContent); ok {
				fmt.Printf("%s:\n", strings.Title(name))
				for _, entry := range freeformContent.Entries {
					if entry.Heading != "" {
						fmt.Printf("  %s\n", entry.Heading)
					}
					for _, content := range entry.Content {
						fmt.Printf("    %s\n", content)
					}
					fmt.Println()
				}
			}
		}
	}
}
