package parser

import (
	"regexp"
	"strings"
	"unicode"
)

type Preprocessor struct{}

func NewPreprocessor() *Preprocessor {
	return &Preprocessor{}
}

// Process preprocesses the input text
func (p *Preprocessor) Process(text string) []string {
	// Split text into lines
	lines := strings.Split(text, "\n")
	var processed []string

	// Process each line
	for _, line := range lines {
		// Clean the line
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle bullet points
		if isBulletPoint(line) {
			processed = append(processed, normalizeBulletPoint(line))
			continue
		}

		// Handle section headers (all caps)
		if isSectionHeader(line) {
			if len(processed) > 0 && processed[len(processed)-1] != "" {
				processed = append(processed, "")
			}
			processed = append(processed, line)
			continue
		}

		// Handle dates and locations (often in parentheses or after commas)
		if strings.Contains(line, ",") || strings.Contains(line, "(") {
			processed = append(processed, line)
			continue
		}

		// Handle contact information (emails, phones, links)
		if strings.Contains(line, "@") || strings.Contains(line, "http") || containsPhoneNumber(line) {
			processed = append(processed, line)
			continue
		}

		// Add other non-empty lines
		processed = append(processed, line)
	}

	return processed
}

func isSectionHeader(line string) bool {
	// Must be relatively short
	if len(line) > 50 {
		return false
	}

	// Should be uppercase and not contain certain characters
	isUpper := true
	hasLetter := false
	for _, r := range line {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				isUpper = false
				break
			}
		}
	}

	return isUpper && hasLetter &&
		!strings.Contains(line, "●") && // Not a bullet point
		!strings.Contains(line, "@") && // Not an email
		!strings.Contains(line, "http") // Not a URL
}

func containsPhoneNumber(line string) bool {
	phonePatterns := []string{
		`\+?\d{1,3}[-.]?\s*\(?\d{3}\)?[-.]?\s*\d{3}[-.]?\s*\d{4}`, // International
		`\(?\d{3}\)?[-.]?\s*\d{3}[-.]?\s*\d{4}`,                   // US/Canada
	}

	for _, pattern := range phonePatterns {
		if matched, _ := regexp.MatchString(pattern, line); matched {
			return true
		}
	}
	return false
}
