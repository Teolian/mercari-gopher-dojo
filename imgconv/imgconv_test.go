package imgconv_test

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	imgconv "road-to-mercari-gopher-dojo-00/imgconv"
)

func TestParseFormat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  imgconv.Format
	}{
		{name: "jpg", input: "jpg", want: imgconv.JPEG},
		{name: "jpeg alias", input: "JPEG", want: imgconv.JPEG},
		{name: "png", input: "png", want: imgconv.PNG},
		{name: "gif", input: "gif", want: imgconv.GIF},
		{name: "unknown", input: "bmp", want: imgconv.Unknown},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := imgconv.ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("ParseFormat(%q) error = %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestNewConverterValidation(t *testing.T) {
	tests := []struct {
		name    string
		dirFn   func(*testing.T) string
		inFmt   string
		outFmt  string
		wantErr bool
	}{
		{
			name: "valid directory",
			dirFn: func(t *testing.T) string {
				t.Helper()
				return copyTestdata(t)
			},
			inFmt:   "png",
			outFmt:  "jpg",
			wantErr: false,
		},
		{
			name:    "empty dir",
			dirFn:   func(t *testing.T) string { return "" },
			inFmt:   "png",
			outFmt:  "jpg",
			wantErr: true,
		},
		{
			name: "not a directory",
			dirFn: func(t *testing.T) string {
				t.Helper()
				f, err := os.CreateTemp(t.TempDir(), "file")
				if err != nil {
					t.Fatalf("CreateTemp: %v", err)
				}
				f.Close()
				return f.Name()
			},
			inFmt:   "png",
			outFmt:  "jpg",
			wantErr: true,
		},
		{
			name: "nonexistent directory",
			dirFn: func(t *testing.T) string {
				t.Helper()
				return filepath.Join(t.TempDir(), "missing")
			},
			inFmt:   "png",
			outFmt:  "jpg",
			wantErr: true,
		},
		{
			name: "unsupported input format",
			dirFn: func(t *testing.T) string {
				t.Helper()
				return copyTestdata(t)
			},
			inFmt:   "bmp",
			outFmt:  "jpg",
			wantErr: true,
		},
		{
			name: "same formats",
			dirFn: func(t *testing.T) string {
				t.Helper()
				return copyTestdata(t)
			},
			inFmt:   "jpg",
			outFmt:  "jpg",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := tc.dirFn(t)
			_, err := imgconv.NewConverter(dir, tc.inFmt, tc.outFmt)
			if tc.wantErr && err == nil {
				t.Fatalf("NewConverter(%q, %q) error = nil, want error", tc.inFmt, tc.outFmt)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("NewConverter(%q, %q) unexpected error: %v", tc.inFmt, tc.outFmt, err)
			}
		})
	}
}

func TestConverterRun(t *testing.T) {
	tests := []struct {
		name         string
		inFmt        string
		outFmt       string
		prepare      func(*testing.T, string)
		wantFailures int
		wantOutput   string
	}{
		{
			name:   "png to jpg",
			inFmt:  "png",
			outFmt: "jpg",
			prepare: func(t *testing.T, dir string) {
				t.Helper()
				err := os.Remove(filepath.Join(dir, "tiny.jpg"))
				if err != nil && !errors.Is(err, fs.ErrNotExist) {
					t.Fatalf("remove tiny.jpg: %v", err)
				}
			},
			wantFailures: 0,
			wantOutput:   "tiny.jpg",
		},
		{
			name:   "jpg to png",
			inFmt:  "jpg",
			outFmt: "png",
			prepare: func(t *testing.T, dir string) {
				t.Helper()
				err := os.Remove(filepath.Join(dir, "tiny.png"))
				if err != nil && !errors.Is(err, fs.ErrNotExist) {
					t.Fatalf("remove tiny.png: %v", err)
				}
			},
			wantFailures: 0,
			wantOutput:   "tiny.png",
		},
		{
			name:   "decode failure counts error",
			inFmt:  "jpg",
			outFmt: "png",
			prepare: func(t *testing.T, dir string) {
				t.Helper()
				err := os.Remove(filepath.Join(dir, "tiny.png"))
				if err != nil && !errors.Is(err, fs.ErrNotExist) {
					t.Fatalf("remove tiny.png: %v", err)
				}
				err = os.Rename(filepath.Join(dir, "not_image.txt"), filepath.Join(dir, "broken.jpg"))
				if err != nil {
					t.Fatalf("rename invalid file: %v", err)
				}
			},
			wantFailures: 1,
			wantOutput:   "tiny.png",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := copyTestdata(t)
			if tc.prepare != nil {
				tc.prepare(t, dir)
			}
			conv, err := imgconv.NewConverter(dir, tc.inFmt, tc.outFmt)
			if err != nil {
				t.Fatalf("NewConverter error: %v", err)
			}
			got := conv.Run()
			if got != tc.wantFailures {
				t.Fatalf("Run() failures = %d, want %d", got, tc.wantFailures)
			}
			if tc.wantOutput != "" {
				out := filepath.Join(dir, tc.wantOutput)
				if _, err := os.Stat(out); err != nil {
					t.Fatalf("expected output %q: %v", out, err)
				}
			}
		})
	}
}

// copyTestdata clones the testdata directory into a temp dir for isolated tests.
func copyTestdata(t *testing.T) string {
	t.Helper()
	candidates := []string{
		filepath.Join("..", "testdata"),
		filepath.Join("testdata"),
	}
	var src string
	for _, cand := range candidates {
		info, err := os.Stat(cand)
		if err == nil && info.IsDir() {
			src = cand
			break
		}
	}
	if src == "" {
		t.Fatalf("testdata directory not found")
	}
	dst := t.TempDir()
	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
	if err != nil {
		t.Fatalf("copy testdata: %v", err)
	}
	return dst
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
