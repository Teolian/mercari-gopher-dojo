package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestGetFileInfo tests HEAD request to get file size and range support
func TestGetFileInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		acceptRanges  string
		contentLength int64
		wantSize      int64
		wantRange     bool
	}{
		{
			name:          "supports range requests",
			acceptRanges:  "bytes",
			contentLength: 1024,
			wantSize:      1024,
			wantRange:     true,
		},
		{
			name:          "does not support range requests",
			acceptRanges:  "none",
			contentLength: 2048,
			wantSize:      2048,
			wantRange:     false,
		},
		{
			name:          "no accept-ranges header",
			acceptRanges:  "",
			contentLength: 512,
			wantSize:      512,
			wantRange:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodHead {
					t.Errorf("Expected HEAD request, got %s", r.Method)
				}

				if tt.acceptRanges != "" {
					w.Header().Set("Accept-Ranges", tt.acceptRanges)
				}
				w.Header().Set("Content-Length", fmt.Sprintf("%d", tt.contentLength))
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			size, supportsRange, err := getFileInfo(server.URL)
			if err != nil {
				t.Fatalf("getFileInfo() error = %v", err)
			}

			if size != tt.wantSize {
				t.Errorf("getFileInfo() size = %d, want %d", size, tt.wantSize)
			}

			if supportsRange != tt.wantRange {
				t.Errorf("getFileInfo() supportsRange = %v, want %v", supportsRange, tt.wantRange)
			}
		})
	}
}

// TestSimpleDownload tests fallback to simple download
func TestSimpleDownload(t *testing.T) {
	t.Parallel()

	testData := []byte("Hello, World! This is test data for simple download.")

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(testData)
	}))
	defer server.Close()

	// Create temp directory for test
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test-file.txt")

	// Download file
	if err := simpleDownload(server.URL, filename); err != nil {
		t.Fatalf("simpleDownload() error = %v", err)
	}

	// Verify file contents
	got, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(got) != string(testData) {
		t.Errorf("Downloaded content = %q, want %q", got, testData)
	}
}

// TestDownloadPart tests downloading a byte range
func TestDownloadPart(t *testing.T) {
	t.Parallel()

	// Full content: "0123456789"
	fullContent := []byte("0123456789")

	tests := []struct {
		name    string
		start   int64
		end     int64
		want    string
		wantErr bool
	}{
		{
			name:  "first half",
			start: 0,
			end:   4,
			want:  "01234",
		},
		{
			name:  "second half",
			start: 5,
			end:   9,
			want:  "56789",
		},
		{
			name:  "middle portion",
			start: 3,
			end:   6,
			want:  "3456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock server with range support
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				rangeHeader := r.Header.Get("Range")
				if rangeHeader == "" {
					t.Error("Missing Range header")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				// Parse range and return partial content
				var start, end int64
				_, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
				if err != nil {
					t.Errorf("Failed to parse Range header: %v", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(fullContent)))
				w.WriteHeader(http.StatusPartialContent)
				w.Write(fullContent[start : end+1])
			}))
			defer server.Close()

			// Create temp directory
			tempDir := t.TempDir()

			// Download part
			err := downloadPart(t.Context(), server.URL, tt.start, tt.end, tempDir, 0)
			if (err != nil) != tt.wantErr {
				t.Fatalf("downloadPart() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			// Verify part file contents
			partPath := filepath.Join(tempDir, "part-0")
			got, err := os.ReadFile(partPath)
			if err != nil {
				t.Fatalf("Failed to read part file: %v", err)
			}

			if string(got) != tt.want {
				t.Errorf("Part content = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestMergeParts tests combining parts into final file
func TestMergeParts(t *testing.T) {
	t.Parallel()

	// Create temp directory with test parts
	tempDir := t.TempDir()

	parts := []string{"Hello", " ", "World", "!"}
	for i, content := range parts {
		partPath := filepath.Join(tempDir, fmt.Sprintf("part-%d", i))
		if err := os.WriteFile(partPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create part %d: %v", i, err)
		}
	}

	// Merge parts
	outputFile := filepath.Join(tempDir, "output.txt")
	if err := mergeParts(tempDir, outputFile, len(parts)); err != nil {
		t.Fatalf("mergeParts() error = %v", err)
	}

	// Verify merged content
	got, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	want := "Hello World!"
	if string(got) != want {
		t.Errorf("Merged content = %q, want %q", got, want)
	}
}

// TestParallelDownload tests parallel download with range requests
func TestParallelDownload(t *testing.T) {
	t.Parallel()

	// Create test data (1KB)
	testData := make([]byte, 1024)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	// Create mock server with range support
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rangeHeader := r.Header.Get("Range")

		if rangeHeader == "" {
			// Full content request
			w.WriteHeader(http.StatusOK)
			w.Write(testData)
			return
		}

		// Parse range
		var start, end int64
		_, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Return partial content
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(testData)))
		w.Header().Set("Accept-Ranges", "bytes")
		w.WriteHeader(http.StatusPartialContent)
		w.Write(testData[start : end+1])
	}))
	defer server.Close()

	// Create temp directory for download
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test-download.bin")

	// Download file
	if err := parallelDownload(server.URL, filename, int64(len(testData))); err != nil {
		t.Fatalf("parallelDownload() error = %v", err)
	}

	// Verify downloaded content matches original
	got, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if len(got) != len(testData) {
		t.Errorf("Downloaded size = %d, want %d", len(got), len(testData))
	}

	for i := range testData {
		if got[i] != testData[i] {
			t.Errorf("Byte mismatch at position %d: got %d, want %d", i, got[i], testData[i])
			break
		}
	}
}

// TestDownloadFile_WithRangeSupport tests full download flow with range support
func TestDownloadFile_WithRangeSupport(t *testing.T) {
	t.Parallel()

	testData := []byte("Test data for range-supported download")

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodHead:
			w.Header().Set("Accept-Ranges", "bytes")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(testData)))
			w.WriteHeader(http.StatusOK)

		case http.MethodGet:
			rangeHeader := r.Header.Get("Range")
			if rangeHeader == "" {
				w.WriteHeader(http.StatusOK)
				w.Write(testData)
				return
			}

			var start, end int64
			fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(testData)))
			w.WriteHeader(http.StatusPartialContent)
			w.Write(testData[start : end+1])
		}
	}))
	defer server.Close()

	// Save original directory and change to temp
	originalDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	// Download file
	if err := downloadFile(server.URL + "/test.txt"); err != nil {
		t.Fatalf("downloadFile() error = %v", err)
	}

	// Verify file exists and content matches
	got, err := os.ReadFile("test.txt")
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(got) != string(testData) {
		t.Errorf("Downloaded content = %q, want %q", got, testData)
	}
}

// TestDownloadFile_WithoutRangeSupport tests fallback to simple download
func TestDownloadFile_WithoutRangeSupport(t *testing.T) {
	t.Parallel()

	testData := []byte("Test data without range support")

	// Create mock server without range support
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodHead:
			w.Header().Set("Accept-Ranges", "none")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(testData)))
			w.WriteHeader(http.StatusOK)

		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			w.Write(testData)
		}
	}))
	defer server.Close()

	// Save original directory and change to temp
	originalDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	// Download file
	if err := downloadFile(server.URL + "/test.txt"); err != nil {
		t.Fatalf("downloadFile() error = %v", err)
	}

	// Verify file exists and content matches
	got, err := os.ReadFile("test.txt")
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(got) != string(testData) {
		t.Errorf("Downloaded content = %q, want %q", got, testData)
	}
}
