# Image Converter

CLI tool for recursive image format conversion. Part of Road to Mercari Gopher Dojo training program (Module 00).

## Description

Converts images between JPG, PNG, and GIF formats in a directory tree. Supports customizable input/output formats via command-line flags.

## Build

```bash
go build -o convert ./cmd/convert
```

Or using Makefile:
```bash
make build
```

## Usage

Convert JPG to PNG (default):
```bash
./convert images/
```

Convert PNG to JPG:
```bash
./convert -i=png -o=jpg images/
```

Convert GIF to PNG:
```bash
./convert -i=gif -o=png images/
```

## Options

- `-i` - Input format (jpg, png, gif). Default: jpg
- `-o` - Output format (jpg, png, gif). Default: png

## Testing

```bash
go test ./...
```

With coverage:
```bash
go test -cover ./...
```

## Project Structure

```
├── cmd/convert/    # CLI entry point
├── imgconv/        # Image conversion library
│   ├── converter.go
│   ├── formats.go
│   ├── util.go
│   └── imgconv_test.go
└── testdata/       # Test fixtures
```

## Implementation Notes

- Uses standard library only (no external dependencies)
- Recursive directory traversal via `filepath.WalkDir`
- Alpha channel flattening for JPEG output (transparent → white)
- Error handling: continues processing on individual file errors
- Exit code: 0 if all conversions succeed, 1 if any fail
