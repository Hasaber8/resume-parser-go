package parser

import (
	"fmt"
	"maps"
	"os"
	"resumeparser/internal/models"
	"strings"
	"unicode"
)

type Parser struct {
	sectionDetectors map[string][]string
}

func NewParser() *Parser {
	return &Parser{
		sectionDetectors: map[string][]string{
			"contact":    {"CONTACT", "PERSONAL"}, // Removed empty string default
			"education":  {"EDUCATION", "ACADEMIC BACKGROUND"},
			"experience": {"PROFESSIONAL EXPERIENCE", "EXPERIENCE"},
			"skills":     {"SKILLS", "TECHNICAL SKILLS", "SKILLS & OTHER"},
		},
	}
}

func (p *Parser) Parse(text string) (*models.Resume, error) {
	lines := p.preprocessLines(text)
	fmt.Fprintf(os.Stderr, "Preprocessed %d lines\n", len(lines))

	sections := p.identifySections(lines)
	fmt.Fprintf(os.Stderr, "Identified sections: %v\n", maps.Keys(sections))

	resume := &models.Resume{
		Raw:      make(map[string]string),
		Sections: make(map[string]models.Sections),
		Metadata: make(map[string]string),
	}

	// Process each section
	for normalizedName, sectionLines := range sections {
		fmt.Fprintf(os.Stderr, "Processing section: %s (%d lines)\n", normalizedName, len(sectionLines))

		sectionType := p.getSectionType(normalizedName)
		content, err := p.parseSection(normalizedName, sectionLines)
		if err != nil {
			return nil, fmt.Errorf("error parsing section %s: %w", normalizedName, err)
		}

		// Store the section
		resume.Sections[normalizedName] = models.Sections{
			Type:    sectionType,
			Content: content,
		}
	}

	return resume, nil
}

func (p *Parser) identifySections(lines []string) map[string][]string {
	sections := make(map[string][]string)
	currentSection := "contact" //default to contact section
	var currentLines []string

	for i, line := range lines {
		// check if this line is a section header
		if newSection := p.detectSection(line, lines, i); newSection != "" {
			if len(currentLines) > 0 {
				sections[currentSection] = currentLines
			}
			currentSection = newSection
			currentLines = nil
		} else {
			currentLines = append(currentLines, line)
		}
	}
	if len(currentLines) > 0 {
		sections[currentSection] = currentLines
	}

	return sections
}

func (p *Parser) detectSection(line string, lines []string, pos int) string {
	// Normalize the input line
	normalizedLine := strings.TrimSpace(strings.ToUpper(line))

	// Remove special characters but preserve spaces
	normalizedLine = strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) && r != '&' { // Preserve & for "SKILLS & OTHER"
			return -1
		}
		return r
	}, normalizedLine)

	// Check each section pattern
	for normalizedName, patterns := range p.sectionDetectors {
		for _, pattern := range patterns {
			normalizedPattern := strings.TrimSpace(strings.ToUpper(pattern))

			// Try exact match first
			if normalizedLine == normalizedPattern {
				return normalizedName
			}

			// Try contains for compound sections (e.g., "SKILLS & OTHER" contains "SKILLS")
			if strings.Contains(normalizedLine, normalizedPattern) {
				return normalizedName
			}
		}
	}

	// Special case for first section if it contains contact information
	if pos < 3 && (strings.Contains(line, "@") || strings.Contains(line, "linkedin")) {
		return "contact"
	}

	return ""
}

func (p *Parser) parseSection(name string, lines []string) (interface{}, error) {
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

func (p *Parser) getSectionType(name string) models.SectionType {
	switch name {
	case "contacts":
		return models.ContactSection
	case "education", "experience":
		return models.TimelineSection
	case "skills":
		return models.ListSection
	default:
		return models.FreeformSection
	}
}

// common helper stuff
func splitByDelimiters(text string, delimiters []string) []string {
	// First, replace all delimiters with a consistent one
	for _, d := range delimiters {
		text = strings.ReplaceAll(text, d, "|")
	}
	return strings.Split(text, "|")
}

func countIndentation(line string) int {
	return len(line) - len(strings.TrimLeft(line, " \t"))
}

func isBulletPoint(line string) bool {
	bullets := []string{"•", "●", "-", "*"}
	line = strings.TrimSpace(line)
	for _, bullet := range bullets {
		if strings.HasPrefix(line, bullet) {
			return true
		}
	}
	return false
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

type dateInfo struct {
	start string
	end   string
}

func extractDates(line string) (dateInfo, bool) {
	separators := []string{"-", "–", "to"}

	for _, sep := range separators {
		if strings.Contains(line, sep) {
			parts := strings.Split(line, sep)
			if len(parts) == 2 {
				start := strings.TrimSpace(parts[0])
				end := strings.TrimSpace(parts[1])
				return dateInfo{start, end}, true
			}
		}
	}
	return dateInfo{}, false
}
