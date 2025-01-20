package extractor

import (
	"context"
	"testing"
)

func TestPDFExtractor(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "non existent file",
			path:    "testdata/nonexistent.pdf",
			wantErr: true,
		},
		// Add more test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			_, err := e.Extract(context.Background(), tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
