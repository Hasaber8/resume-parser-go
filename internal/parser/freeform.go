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
	var buffer []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			// Empty line could indicate a section break
			if len(buffer) > 0 {
				if currentEntry != nil {
					currentEntry.Content = append(currentEntry.Content, strings.Join(buffer, " "))
					content.Entries = append(content.Entries, *currentEntry)
				} else {
					// Create entry without explicit heading
					content.Entries = append(content.Entries, models.FreeformEntry{
						Content: []string{strings.Join(buffer, " ")},
					})
				}
				buffer = nil
				currentEntry = nil
			}
			continue
		}

		indentation := countIndentation(line)
		if indentation == 0 && !isBulletPoint(line) && !strings.HasPrefix(line, "-") {
			// This might be a new heading
			if currentEntry != nil {
				if len(buffer) > 0 {
					currentEntry.Content = append(currentEntry.Content, strings.Join(buffer, " "))
				}
				content.Entries = append(content.Entries, *currentEntry)
				buffer = nil
			}
			currentEntry = &models.FreeformEntry{
				Heading: line,
				Content: make([]string, 0),
			}
		} else {
			// This is content
			if isBulletPoint(line) {
				// Handle bullet points
				if len(buffer) > 0 {
					if currentEntry != nil {
						currentEntry.Content = append(currentEntry.Content, strings.Join(buffer, " "))
					}
					buffer = nil
				}
				line = removeBulletPoint(line)
				if currentEntry != nil {
					currentEntry.Content = append(currentEntry.Content, line)
				} else {
					content.Entries = append(content.Entries, models.FreeformEntry{
						Content: []string{line},
					})
				}
			} else {
				// Accumulate text in buffer
				buffer = append(buffer, line)
			}
		}
	}

	// Handle any remaining content
	if len(buffer) > 0 {
		if currentEntry != nil {
			currentEntry.Content = append(currentEntry.Content, strings.Join(buffer, " "))
			content.Entries = append(content.Entries, *currentEntry)
		} else {
			content.Entries = append(content.Entries, models.FreeformEntry{
				Content: []string{strings.Join(buffer, " ")},
			})
		}
	} else if currentEntry != nil {
		content.Entries = append(content.Entries, *currentEntry)
	}

	return content, nil
}

func countIndentation(line string) int {
	indentation := 0
	for _, char := range line {
		if char == ' ' || char == '\t' {
			indentation++
		} else {
			break
		}
	}
	return indentation
}
