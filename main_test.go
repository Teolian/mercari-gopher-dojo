package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   "",
			wantErr: false,
		},
		{
			name:    "single line",
			input:   "hello world",
			wantErr: false,
		},
		{
			name:    "multiple lines",
			input:   "line1\nline2\nline3\n",
			wantErr: false,
		},
		{
			name:    "binary data",
			input:   string([]byte{0x00, 0xFF, 0x42, 0x13}),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			writer := &bytes.Buffer{}

			err := cat(reader, writer)

			if (err != nil) != tt.wantErr {
				t.Errorf("cat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got := writer.String(); got != tt.input {
				t.Errorf("cat() output = %q, want %q", got, tt.input)
			}
		})
	}
}

func TestCatLargeFile(t *testing.T) {
	// Test with large input to verify efficient copying
	size := 1024 * 1024 // 1MB
	input := strings.Repeat("A", size)

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	err := cat(reader, writer)
	if err != nil {
		t.Fatalf("cat() unexpected error: %v", err)
	}

	if writer.Len() != size {
		t.Errorf("cat() copied %d bytes, want %d", writer.Len(), size)
	}
}
