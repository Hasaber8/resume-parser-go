package parser

import (
	"fmt"
	"os"
	"regexp"
	"resumeparser/internal/models"
	"strings"
	"unicode"
)

// Parser represents the resume parser
type Parser struct {
	sectionDetectors map[string][]string
	preprocessor     *Preprocessor
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	p := &Parser{
		sectionDetectors: make(map[string][]string),
		preprocessor:     NewPreprocessor(),
	}

	// Initialize section detectors with common variations
	p.sectionDetectors["education"] = []string{
		"education",
		"academic background",
		"academic history",
		"educational background",
	}

	p.sectionDetectors["experience"] = []string{
		"experience",
		"work experience",
		"employment history",
		"professional experience",
	}

	p.sectionDetectors["skills"] = []string{
		"skills",
		"technical skills",
		"core competencies",
		"expertise",
	}

	p.sectionDetectors["projects"] = []string{
		"projects",
		"personal projects",
		"project experience",
	}

	p.sectionDetectors["achievements"] = []string{
		"achievements",
		"awards",
		"honors",
		"accomplishments",
	}

	p.sectionDetectors["contact"] = []string{
		"contact",
		"contact information",
		"personal information",
	}

	return p
}

// Parse the input text and returns a structured Resume
func (p *Parser) Parse(text string) (*models.Resume, error) {
	if text == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := p.preprocessor.Process(text)
	fmt.Fprintf(os.Stderr, "Preprocessed %d lines\n", len(lines))

	sections := p.identifySections(lines)
	var sectionNames []string
	for name := range sections {
		sectionNames = append(sectionNames, name)
	}
	fmt.Fprintf(os.Stderr, "Identified sections: %s\n", strings.Join(sectionNames, ", "))

	resume := &models.Resume{
		Raw:      make(map[string]string),
		Sections: make(map[string]models.Sections),
		Metadata: make(map[string]string),
	}

	// Process each section
	for name, sectionLines := range sections {
		fmt.Fprintf(os.Stderr, "Processing section: %s (%d lines)\n", name, len(sectionLines))

		sectionType := p.getSectionType(name)
		content, err := p.parseSection(name, sectionLines)
		if err != nil {
			return nil, fmt.Errorf("error parsing section %s: %w", name, err)
		}

		// Store the section
		resume.Sections[name] = models.Sections{
			Type:    sectionType,
			Content: content,
		}
	}

	return resume, nil
}

// identifySections identifies and groups lines into sections
func (p *Parser) identifySections(lines []string) map[string][]string {
	sections := make(map[string][]string)
	var currentSection string
	var currentLines []string

	for i, line := range lines {
		if section := p.detectSection(line); section != "" {
			// Store previous section if it exists
			if currentSection != "" && len(currentLines) > 0 {
				sections[currentSection] = currentLines
			}
			currentSection = section
			currentLines = nil
		} else if currentSection != "" {
			currentLines = append(currentLines, line)
		} else if i == 0 {
			// First line is likely contact information
			currentSection = "contact"
			currentLines = append(currentLines, line)
		}
	}

	// Add the last section
	if currentSection != "" && len(currentLines) > 0 {
		sections[currentSection] = currentLines
	}

	return sections
}

// detectSection tries to identify if a line is a section header
func (p *Parser) detectSection(line string) string {
	line = strings.TrimSpace(strings.ToLower(line))
	if line == "" {
		return ""
	}

	// Quick check for common section headers
	for name, patterns := range p.sectionDetectors {
		for _, pattern := range patterns {
			if line == pattern {
				return name
			}
		}
	}

	return ""
}

func (p *Parser) getSectionType(name string) models.SectionType {
	switch name {
	case "contact":
		return models.ContactSection
	case "education", "experience", "projects":
		return models.TimelineSection
	case "skills", "achievements":
		return models.ListSection
	default:
		return models.FreeformSection
	}
}

func (p *Parser) parseSection(name string, lines []string) (interface{}, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty section")
	}

	// Pre-process lines
	lines = p.cleanSectionLines(lines)

	switch p.getSectionType(name) {
	case models.ContactSection:
		return p.parseContact(lines)
	case models.TimelineSection:
		return p.parseTimeline(lines)
	case models.ListSection:
		return p.parseList(lines)
	default:
		return p.parseFreeform(lines)
	}
}

// containsAny checks if the string contains any of the given substrings
func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(strings.ToLower(s), strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

// cleanSectionLines removes empty lines and normalizes formatting
func (p *Parser) cleanSectionLines(lines []string) []string {
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return cleaned
}

// isBulletPoint checks if a line starts with any known bullet point marker
func isBulletPoint(line string) bool {
	line = strings.TrimSpace(line)
	if line == "" {
		return false
	}

	// Check for numbered bullets (e.g., "1.", "2)", "(1)", "a.", "b)")
	if matched, _ := regexp.MatchString(`^(?:\d+[\.\)]|\([a-zA-Z0-9]+\)|[a-zA-Z][\.\)])`, line); matched {
		return true
	}

	// Check for standard bullet points
	for _, bullet := range bulletPoints {
		if strings.HasPrefix(line, bullet) {
			return true
		}
	}

	return false
}

// normalizeBulletPoint converts any bullet point style to a standard bullet point
func normalizeBulletPoint(line string) string {
	line = strings.TrimSpace(line)

	// Handle numbered bullets
	if matched, _ := regexp.MatchString(`^(?:\d+[\.\)]|\([a-zA-Z0-9]+\)|[a-zA-Z][\.\)])`, line); matched {
		// Find the end of the bullet marker
		idx := strings.IndexFunc(line, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '(' && r != ')' && r != '.'
		})
		if idx == -1 {
			idx = len(line)
		}
		return "●" + strings.TrimSpace(line[idx:])
	}

	// Handle standard bullet points
	for _, bullet := range bulletPoints {
		if strings.HasPrefix(line, bullet) {
			return "●" + strings.TrimSpace(strings.TrimPrefix(line, bullet))
		}
	}

	return line
}

func removeBulletPoint(line string) string {
	bullets := strings.TrimSpace(line)
	line = strings.TrimSpace(line)
	for _, bullet := range bullets {
		if strings.HasPrefix(line, string(bullet)) {
			return strings.TrimSpace(strings.TrimPrefix(line, string(bullet)))
		}
	}
	return line
}

var bulletPoints = []string{
	"●", // Standard bullet
	"•", // Alternative bullet
	"⚬", // White bullet
	"○", // White circle
	"▪", // Black square
	"■", // Black square
	"◆", // Black diamond
	"★", // Star
	"-", // Hyphen
	"*", // Asterisk
	">", // Greater than
	"→", // Arrow
	"+", // Plus
	"◦",
}
