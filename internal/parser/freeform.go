package parser

import (
	"resumeparser/internal/models"
	"strings"
)

func (p *Parser) parseFreeform(lines []string) (*models.FreeformContent, error) {
	content := &models.FreeformContent{
		Entries: make([]models.FreeformEntry, 0),
	}

	var currentEntry *models.FreeformEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		indentation := countIndentation(line)
		if indentation == 0 && !isBulletPoint(line) {
			if currentEntry != nil {
				content.Entries = append(content.Entries, *currentEntry)
			}
			currentEntry = &models.FreeformEntry{
				Heading: line,
				Content: make([]string, 0),
			}
		} else if currentEntry != nil {
			if isBulletPoint(line) {
				line = removeBulletPoint(line)
			}
			currentEntry.Content = append(currentEntry.Content, line)
		}
	}

	if currentEntry != nil {
		content.Entries = append(content.Entries, *currentEntry)
	}

	return content, nil
}
