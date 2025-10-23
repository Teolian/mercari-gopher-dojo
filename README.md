# ft_cat

Implementation of `cat` command using `io.Reader` and `io.Writer` interfaces. Part of Road to Mercari Gopher Dojo training program (Module 01).

## Description

Demonstrates Go interface composition by implementing a `cat` function that accepts `io.Reader` and `io.Writer`, enabling abstraction over data sources and destinations.

## Build

```bash
go build -o ft_cat
```

## Usage

Read from stdin:
```bash
echo "hello" | ./ft_cat
```

Read from file:
```bash
./ft_cat testdata/simple.txt
```

Concatenate multiple files:
```bash
./ft_cat testdata/simple.txt testdata/multiline.txt
```

## Testing

```bash
go test
```

With coverage:
```bash
go test -cover
```

## Implementation

**Core function:**
```go
func cat(r io.Reader, w io.Writer) error
```

This function signature demonstrates the power of Go interfaces:
- Works with any `io.Reader`: files, stdin, strings, bytes, network connections
- Works with any `io.Writer`: files, stdout, buffers, network connections
- Testable without file I/O using `strings.Reader` and `bytes.Buffer`

**Key benefits of io.Reader/Writer:**
- Composability: combine readers/writers with `io.MultiReader`, `io.TeeReader`
- Testability: mock I/O operations without touching filesystem
- Flexibility: same function works for files, network, memory
- Efficiency: `io.Copy` uses optimized system calls when available

## Project Structure

```
├── main.go          # CLI implementation
├── main_test.go     # Unit tests
└── testdata/        # Test fixtures
    ├── simple.txt
    └── multiline.txt
```
