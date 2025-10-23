package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

const (
	// Number of parallel downloads (not too many to avoid excessive requests)
	numParts = 4
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: download <URL>")
		os.Exit(1)
	}

	url := os.Args[1]

	if err := downloadFile(url); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// downloadFile downloads a file using parallel range requests
func downloadFile(url string) error {
	// Get file size and check if range requests are supported
	size, supportsRange, err := getFileInfo(url)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Extract filename from URL
	filename := filepath.Base(url)

	if !supportsRange || size == 0 {
		// Fall back to simple download
		return simpleDownload(url, filename)
	}

	// Download in parallel using range requests
	return parallelDownload(url, filename, size)
}

// getFileInfo makes HEAD request to get file size and range support
func getFileInfo(url string) (int64, bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, false, err
	}
	defer resp.Body.Close()

	// Check if server supports range requests
	supportsRange := resp.Header.Get("Accept-Ranges") == "bytes"

	// Get content length
	size := resp.ContentLength

	return size, supportsRange, nil
}

// simpleDownload downloads file without range requests
func simpleDownload(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// parallelDownload downloads file in parallel using range requests
func parallelDownload(url, filename string, totalSize int64) error {
	// Calculate part size
	partSize := totalSize / int64(numParts)

	// Create temporary directory for parts
	tempDir, err := os.MkdirTemp("", "download-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use errgroup for parallel downloads with error handling
	g, ctx := errgroup.WithContext(ctx)

	// Download each part
	for i := 0; i < numParts; i++ {
		i := i // Capture loop variable
		start := int64(i) * partSize
		end := start + partSize - 1

		// Last part gets remaining bytes
		if i == numParts-1 {
			end = totalSize - 1
		}

		g.Go(func() error {
			return downloadPart(ctx, url, start, end, tempDir, i)
		})
	}

	// Wait for all downloads to complete
	if err := g.Wait(); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Merge parts into final file
	return mergeParts(tempDir, filename, numParts)
}

// downloadPart downloads a specific byte range
func downloadPart(ctx context.Context, url string, start, end int64, tempDir string, partNum int) error {
	// Create request with Range header
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	// Set Range header: "Range: bytes=start-end"
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for partial content response
	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	// Write part to temporary file
	partPath := filepath.Join(tempDir, fmt.Sprintf("part-%d", partNum))
	out, err := os.Create(partPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// mergeParts combines all parts into final file
func mergeParts(tempDir, filename string, numParts int) error {
	// Create output file
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Copy each part in order
	for i := 0; i < numParts; i++ {
		partPath := filepath.Join(tempDir, fmt.Sprintf("part-%d", i))

		partFile, err := os.Open(partPath)
		if err != nil {
			return fmt.Errorf("failed to open part %d: %w", i, err)
		}

		if _, err := io.Copy(out, partFile); err != nil {
			partFile.Close()
			return fmt.Errorf("failed to copy part %d: %w", i, err)
		}

		partFile.Close()
	}

	fmt.Printf("Downloaded: %s\n", filename)
	return nil
}
