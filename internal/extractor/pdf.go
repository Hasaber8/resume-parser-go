package extractor

import (
	"context"
	"fmt"
	"strings"

	"github.com/dslipak/pdf"
)

type PdfExtractor interface {
	Extract(ctx context.Context, path string) (string, error)
}

type pdfExtractor struct{}

func New() PdfExtractor {
	return &pdfExtractor{}
}

func (e *pdfExtractor) Extract(ctx context.Context, path string) (string, error) {
	f, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}

	var text strings.Builder
	for i := 1; i <= f.NumPage(); i++ {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			p := f.Page(i)
			content, err := p.GetPlainText(nil)
			if err != nil {
				return "", fmt.Errorf("failed to extract page %d: %w", i, err)
			}
			text.WriteString(content)
		}
	}
	return text.String(), nil
}
