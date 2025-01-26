package parser

import (
	"regexp"
	"resumeparser/internal/models"
	"strings"
)

// parseContact extracts contact information from the given lines
func (p *Parser) parseContact(lines []string) (*models.ContactContent, error) {
	content := &models.ContactContent{
		Email:  make([]string, 0),
		Number: make([]string, 0),
		Social: make(map[string]string),
	}

	// Regular expressions for different contact information
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	phoneRegex := regexp.MustCompile(`(?:(?:\+?\d{1,3}[-.]?\s*)?(?:\(?\d{3}\)?[-.]?\s*)?\d{3}[-.]?\s*\d{4})`)
	linkedinRegex := regexp.MustCompile(`(?i)linkedin\.com/(?:in|profile)/[a-zA-Z0-9_-]+`)
	githubRegex := regexp.MustCompile(`(?i)github\.com/[a-zA-Z0-9_-]+`)

	for _, line := range lines {
		// Extract email addresses
		emails := emailRegex.FindAllString(line, -1)
		content.Email = append(content.Email, emails...)

		// Extract phone numbers
		phones := phoneRegex.FindAllString(line, -1)
		content.Number = append(content.Number, phones...)

		// Extract LinkedIn profile
		if linkedin := linkedinRegex.FindString(line); linkedin != "" {
			content.Social["linkedin"] = linkedin
		}

		// Extract GitHub profile
		if github := githubRegex.FindString(line); github != "" {
			content.Social["github"] = github
		}

		// Try to identify name and location
		if content.Name == "" && !containsAny(line, "@", "http", "linkedin", "github") {
			parts := strings.Split(line, "•")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part != "" {
					if content.Name == "" {
						content.Name = part
					} else if content.Location == "" && !containsAny(part, "@", "http") {
						content.Location = part
					}
				}
			}
		}
	}

	return content, nil
}
