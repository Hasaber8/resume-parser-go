package parser

import (
	"resumeparser/internal/models"
	"strings"
)

func (p *Parser) parseList(lines []string) (*models.ListContent, error) {
	content := &models.ListContent{}
	var currentCatergory *models.ListCategory

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				if currentCatergory != nil {
					content.Categories = append(content.Categories, *currentCatergory)
				}
				currentCatergory = &models.ListCategory{
					Name:  strings.TrimSpace(parts[0]),
					Items: make([]string, 0),
				}
				// handle items on same line
				items := strings.Split(strings.TrimSpace(parts[1]), ",")
				for _, item := range items {
					if item = strings.TrimSpace(item); item != "" {
						currentCatergory.Items = append(currentCatergory.Items, item)
					}
				}
			}
		} else if currentCatergory != nil && line != "" {
			if isBulletPoint(line) {
				currentCatergory.Items = append(currentCatergory.Items, removeBulletPoint(line))
			} else {
				items := strings.Split(line, ",")
				for _, item := range items {
					if item = strings.TrimSpace(item); item != "" {
						currentCatergory.Items = append(currentCatergory.Items, item)
					}
				}
			}
		}
	}
	if currentCatergory != nil {
		content.Categories = append(content.Categories, *currentCatergory)
	}

	return content, nil
}
