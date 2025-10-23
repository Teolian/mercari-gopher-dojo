package imgconv

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

// Converter performs recursive image conversions within a directory.
//
// It supports JPG/JPEG, PNG, GIF using Go's standard image packages.
// When converting to JPEG, alpha is flattened over a white background.
// Non-image files are reported as errors but do not abort the entire run.
//
// Example:
//
//	c, _ := NewConverter("./images", "jpg", "png")
//	errs := c.Run() // returns number of per-file errors
//
// Default behavior (per task) is JPG -> PNG.
// Bonus flags allow arbitrary direction via -i and -o.
//
// Error policy: continue-on-error, aggregate count.
type Converter struct {
	Dir          string
	InputFormat  Format
	OutputFormat Format
}

// NewConverter validates directory and formats, returning a Converter.
func NewConverter(dir string, inFmt, outFmt string) (*Converter, error) {
	if dir == "" {
		return nil, errors.New("invalid argument")
	}
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}
	inf, err := ParseFormat(inFmt)
	if err != nil {
		return nil, err
	}
	outf, err := ParseFormat(outFmt)
	if err != nil {
		return nil, err
	}
	if inf == Unknown || outf == Unknown {
		return nil, fmt.Errorf("unsupported format (supported: jpg|png|gif)")
	}
	if inf == outf {
		return nil, fmt.Errorf("input and output formats are the same: %s", inf)
	}
	return &Converter{Dir: dir, InputFormat: inf, OutputFormat: outf}, nil
}

// Run walks the directory tree and converts matching files.
// Returns number of files that failed to convert.
func (c *Converter) Run() int {
	var failures int
	_ = filepath.WalkDir(c.Dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("error: %s: %v\n", path, err)
			failures++
			return nil
		}
		if d.IsDir() {
			return nil
		}
		// Check extension against the chosen input format
		if !hasExt(path, c.InputFormat) {
			// As per example, report invalid files that don't match required input
			// but continue processing. Suppress when it's already the output format.
			if !hasExt(path, c.OutputFormat) {
				fmt.Printf("error: %s is not a valid file\n", path)
			}
			return nil
		}
		if err := c.convertOne(path); err != nil {
			fmt.Printf("error: %s: %v\n", path, err)
			failures++
		}
		return nil
	})
	return failures
}

func (c *Converter) convertOne(srcPath string) error {
	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer in.Close()

	img, err := decodeByFormat(in, c.InputFormat)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	dstPath := replaceExt(srcPath, c.OutputFormat)
	// If destination already exists, skip to keep idempotence.
	if _, err := os.Stat(dstPath); err == nil {
		return nil
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	return encodeByFormat(out, img, c.OutputFormat)
}

// decodeByFormat dispatches to a specific decoder based on declared format.
func decodeByFormat(r io.Reader, f Format) (image.Image, error) {
	switch f {
	case JPEG:
		return jpeg.Decode(r)
	case PNG:
		return png.Decode(r)
	case GIF:
		return gif.Decode(r)
	default:
		return nil, fmt.Errorf("unsupported input format: %s", f)
	}
}

// encodeByFormat dispatches to a specific encoder based on declared format.
// For JPEG, alpha is flattened over a white background.
func encodeByFormat(w io.Writer, img image.Image, f Format) error {
	switch f {
	case JPEG:
		// Flatten alpha if present
		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
		draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Over)
		return jpeg.Encode(w, rgba, &jpeg.Options{Quality: 90})
	case PNG:
		return png.Encode(w, img)
	case GIF:
		palettedImg := toPaletted(img)
		return gif.Encode(w, palettedImg, nil)
	default:
		return fmt.Errorf("unsupported output format: %s", f)
	}
}
