package parser

import (
	"resumeparser/internal/models"
	"strings"
)

func (p *Parser) parseList(lines []string) (*models.ListContent, error) {
	content := &models.ListContent{
		Categories: make([]models.ListCategory, 0),
	}
	var currentCategory *models.ListCategory

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if line is a category header
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				if currentCategory != nil {
					content.Categories = append(content.Categories, *currentCategory)
				}
				currentCategory = &models.ListCategory{
					Name:  strings.TrimSpace(parts[0]),
					Items: make([]string, 0),
				}
				// Handle items on same line as category
				items := parseItems(parts[1])
				currentCategory.Items = append(currentCategory.Items, items...)
			}
		} else if currentCategory != nil {
			// Handle items under current category
			if isBulletPoint(line) {
				item := removeBulletPoint(line)
				currentCategory.Items = append(currentCategory.Items, item)
			} else {
				items := parseItems(line)
				currentCategory.Items = append(currentCategory.Items, items...)
			}
		} else {
			// Handle items without category
			currentCategory = &models.ListCategory{
				Name:  "",
				Items: make([]string, 0),
			}
			if isBulletPoint(line) {
				item := removeBulletPoint(line)
				currentCategory.Items = append(currentCategory.Items, item)
			} else {
				items := parseItems(line)
				currentCategory.Items = append(currentCategory.Items, items...)
			}
		}
	}

	if currentCategory != nil {
		content.Categories = append(content.Categories, *currentCategory)
	}

	return content, nil
}

// parseItems splits a line into individual items, handling various delimiters
func parseItems(line string) []string {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	// Handle different types of delimiters
	delimiters := []string{",", "•", "|", ";"}
	for _, delimiter := range delimiters {
		if strings.Contains(line, delimiter) {
			var items []string
			for _, item := range strings.Split(line, delimiter) {
				if item = strings.TrimSpace(item); item != "" {
					items = append(items, item)
				}
			}
			return items
		}
	}

	// If no delimiters found, treat as single item
	return []string{line}
}
