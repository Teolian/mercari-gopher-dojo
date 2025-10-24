# Road to Mercari Gopher Dojo

Training program with 4 modules. Each module is in a separate branch for independent demonstration.

## Structure

Each module exists in its own branch with complete, runnable code in the repository root.

```bash
# Clone repository
git clone https://github.com/Teolian/mercari-gopher-dojo.git
cd mercari-gopher-dojo

# Switch to any module
git checkout module-00
```

## Modules

### Module 00 - Image Converter
**Branch**: [`module-00`](../../tree/module-00)

CLI tool for recursive image format conversion (JPG/PNG/GIF).

```bash
git checkout module-00
go build -o convert ./cmd/convert
./convert -i=jpg -o=png images/
```

**Features:**
- Recursive directory traversal
- Customizable formats via `-i` and `-o` flags
- Alpha flattening for JPEG output
- Table-driven tests with parallel execution

---

### Module 01 - I/O and Testing
**Branch**: [`module-01`](../../tree/module-01)

Implementation of `cat` command using `io.Reader` and `io.Writer` interfaces.

```bash
git checkout module-01
go build -o ft_cat
./ft_cat testdata/simple.txt
echo "hello" | ./ft_cat
```

**Features:**
- Core function using io.Reader/Writer abstraction
- Works with files, stdin, and multiple sources
- Table-driven tests with various input types
- Demonstrates Go interface composition

---

### Module 02 - Concurrency
**Branch**: [`module-02`](../../tree/module-02)

Go concurrency patterns using goroutines, channels, and context.

```bash
git checkout module-02
cd ex00 && go build -o typing_game && ./typing_game
cd ex01 && go build -o download && ./download https://example.com/file.zip
```

**Exercise 00:** Typing game with 30-second timer using goroutines and channels.

**Exercise 01:** Parallel file downloader with HTTP Range requests and `errgroup`.

**Concepts:**
- Goroutines for non-blocking I/O
- Channels for communication
- Context with timeout and cancellation
- Select statement for multiplexing
- errgroup for parallel error handling

---

### Module 03 - HTTP API
**Branch**: [`module-03`](../../tree/module-03)

Fortune-telling HTTP API server with JSON responses.

```bash
git checkout module-03
cd ex00
go build -o omikuji
./omikuji 8080
curl localhost:8080
```

**Exercise 00:** Omikuji (fortune-telling) API that returns random fortune in JSON format.

**Features:**
- 7 fortune types: Dai-kichi, Kichi, Chuu-kichi, Sho-kichi, Sue-kichi, Kyo, Dai-kyo
- New Year special logic (Jan 1-3): Always returns Dai-kichi
- Testable time dependency using Clock function type
- HTTP handler tests using `httptest`
- JSON encoding with struct tags
- 86.7% test coverage

