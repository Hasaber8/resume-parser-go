package parser

import (
	"regexp"
	"resumeparser/internal/models"
	"strings"
)

type dateInfo struct {
	start string
	end   string
}

// extractDates attempts to extract start and end dates from a line of text
func extractDates(line string) (dateInfo, bool) {
	// Common date formats
	datePatterns := []string{
		`(January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{4}`,
		`(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\.?\s+\d{4}`,
		`\d{2}/\d{4}`,
		`\d{2}-\d{4}`,
	}

	var dates []string
	for _, pattern := range datePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(line, -1)
		dates = append(dates, matches...)
	}

	// Check for "Present" or "Current"
	if strings.Contains(line, "Present") || strings.Contains(line, "Current") {
		if len(dates) > 0 {
			return dateInfo{
				start: dates[0],
				end:   "Present",
			}, true
		}
	}

	if len(dates) >= 2 {
		return dateInfo{
			start: dates[0],
			end:   dates[1],
		}, true
	} else if len(dates) == 1 {
		return dateInfo{
			start: dates[0],
			end:   dates[0],
		}, true
	}

	return dateInfo{}, false
}

func (p *Parser) parseTimeline(lines []string) (*models.TimelineContent, error) {
	content := &models.TimelineContent{
		Entries: make([]models.TimelineEntry, 0),
	}
	var currentEntry *models.TimelineEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		indentation := countIndentation(line)
		if indentation == 0 && !isBulletPoint(line) {
			// New entry
			if currentEntry != nil {
				content.Entries = append(content.Entries, *currentEntry)
			}
			currentEntry = &models.TimelineEntry{
				Organization: line,
				Details:      make([]string, 0),
				Metadata:     make(map[string]string),
			}
		} else if currentEntry != nil {
			// Process entry details
			if isBulletPoint(line) {
				detail := removeBulletPoint(line)
				currentEntry.Details = append(currentEntry.Details, detail)
			} else if dates, ok := extractDates(line); ok {
				currentEntry.StartDate = dates.start
				currentEntry.EndDate = dates.end
			} else if indentation == 1 {
				// This could be title or location
				if strings.Contains(line, ",") {
					parts := strings.Split(line, ",")
					if len(parts) >= 2 {
						currentEntry.Title = strings.TrimSpace(parts[0])
						currentEntry.Location = strings.TrimSpace(parts[1])
					} else {
						currentEntry.Location = line
					}
				} else {
					currentEntry.Title = line
				}
			}
		}
	}

	// Add the last entry if exists
	if currentEntry != nil {
		content.Entries = append(content.Entries, *currentEntry)
	}

	return content, nil
}
