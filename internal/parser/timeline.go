package parser

import (
	"resumeparser/internal/models"
	"strings"
)

func (p *Parser) parseTimeline(lines []string) (*models.TimelineContent, error) {
	content := &models.TimelineContent{}
	var currentEntry *models.TimelineEntry

	for _, line := range lines {
		indentation := countIndentation(line)
		line = strings.TrimSpace(line)

		if indentation == 0 && !isBulletPoint(line) {
			if currentEntry != nil {
				content.Entries = append(content.Entries, *currentEntry)
			}
			currentEntry = &models.TimelineEntry{
				Organization: line,
				Details:      make([]string, 0),
				Metadata:     make(map[string]string),
			}
		} else if currentEntry != nil {
			if isBulletPoint(line) {
				detail := removeBulletPoint(line)
				currentEntry.Details = append(currentEntry.Details, detail)
			} else if dates, ok := extractDates(line); ok {
				currentEntry.StartDate = dates.start
				currentEntry.EndDate = dates.end
			} else if indentation == 1 {
				if strings.Contains(line, ",") {
					currentEntry.Location = line
				} else {
					currentEntry.Title = line
				}
			}
		}
	}

	if currentEntry != nil {
		content.Entries = append(content.Entries, *currentEntry)
	}

	return content, nil
}
