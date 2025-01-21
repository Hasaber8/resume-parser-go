package parser

import (
	"regexp"
	"strings"
)

func (p *Parser) preprocessLines(text string) []string {
	var lines []string

	// Helper functions
	isSection := func(s string) bool {
		s = strings.TrimSpace(s)
		return len(s) > 0 &&
			strings.ToUpper(s) == s &&
			len(s) > 3 && // Avoid single words
			!strings.Contains(s, "●") // Not a bullet point
	}

	isDate := func(s string) bool {
		datePatterns := []string{
			`(January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{4}`,
			`(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\s+\d{4}`,
			`Present`,
		}
		for _, pattern := range datePatterns {
			if matched, _ := regexp.MatchString(pattern, s); matched {
				return true
			}
		}
		return false
	}

	// Split on bullet points first
	parts := strings.Split(text, "●")

	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Split part into potential lines
		words := strings.Fields(part)
		var currentLine strings.Builder

		for _, word := range words {
			// If this word starts a new section
			if isSection(word) {
				if currentLine.Len() > 0 {
					lines = append(lines, strings.TrimSpace(currentLine.String()))
					currentLine.Reset()
				}
				currentLine.WriteString(word)
				continue
			}

			// If this is a date and we have content
			if isDate(word) && currentLine.Len() > 0 {
				currentLine.WriteString(" " + word)
				lines = append(lines, strings.TrimSpace(currentLine.String()))
				currentLine.Reset()
				continue
			}

			// Otherwise append to current line
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		}

		// Add remaining content
		if currentLine.Len() > 0 {
			line := strings.TrimSpace(currentLine.String())
			// If this isn't the first part, it's a bullet point
			if i > 0 {
				line = "●" + line
			}
			lines = append(lines, line)
		}
	}

	return lines
}
