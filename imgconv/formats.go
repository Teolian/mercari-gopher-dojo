package imgconv

import (
	"path/filepath"
	"strings"
)

// Format is a supported image format keyword.
// Supported: jpg, png, gif (jpeg maps to jpg).
type Format string

const (
	Unknown Format = ""
	JPEG    Format = "jpg"
	PNG     Format = "png"
	GIF     Format = "gif"
)

// ParseFormat normalizes and validates a format keyword.
func ParseFormat(s string) (Format, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "jpg", "jpeg":
		return JPEG, nil
	case "png":
		return PNG, nil
	case "gif":
		return GIF, nil
	default:
		return Unknown, nil
	}
}

// hasExt checks whether a path has the extension for a given format.
func hasExt(path string, f Format) bool {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	if f == JPEG {
		return ext == "jpg" || ext == "jpeg"
	}
	return ext == string(f)
}

// replaceExt swaps the file extension to match the output format.
func replaceExt(path string, f Format) string {
	base := strings.TrimSuffix(path, filepath.Ext(path))
	return base + "." + string(f)
}
