package parser

import (
	"regexp"
	"resumeparser/internal/models"
	"strings"
)

func (p *Parser) parseContact(lines []string) (*models.ContactContent, error) {
	contact := &models.ContactContent{
		Email:  make([]string, 0),
		Number: make([]string, 0),
		Social: make(map[string]string),
	}

	// first line is usually name
	if len(lines) > 0 {
		contact.Name = strings.TrimSpace(lines[0])
		lines = lines[1:]
	}

	emailRegex := regexp.MustCompile(`[^@]+@[^@]+\.[^@]+`)
	phoneRegex := regexp.MustCompile(`[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}`)

	for _, line := range lines {
		parts := splitByDelimiters(line, []string{".", "|", ","})

		for _, part := range parts {
			part = strings.TrimSpace(part)
			switch {
			case emailRegex.MatchString(part):
				contact.Email = append(contact.Email, part)
			case phoneRegex.MatchString(part):
				contact.Number = append(contact.Number, part)
			case strings.Contains(strings.ToLower(part), "linkedin.com"):
				contact.Social["linkedin"] = part
			case strings.Contains(strings.ToLower(part), "github.com"):
				contact.Social["github"] = part
			default:
				// If part contains city/state pattern, treat as location
				if strings.Contains(part, ",") && len(strings.Split(part, ",")) == 2 {
					contact.Location = part
				}
			}
		}
	}
	return contact, nil
}
