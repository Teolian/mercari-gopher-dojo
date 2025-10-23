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
**Branch**: [`module-01`](../../tree/module-01) _(in progress)_

Implementation of `cat` command using `io.Reader` and `io.Writer` interfaces.

```bash
git checkout module-01
go build
./ft_cat file.txt
```

---

### Module 02 - Concurrency
**Branch**: [`module-02`](../../tree/module-02) _(planned)_

**Exercise 00:** Typing game with 30-second timer using goroutines and channels.

**Exercise 01:** Parallel file downloader with HTTP Range requests and `errgroup`.

---

### Module 03 - HTTP API
**Branch**: [`module-03`](../../tree/module-03) _(planned)_

Fortune-telling API server with JSON responses and HTTP testing.

```bash
git checkout module-03
go build
./omikuji 8080
curl localhost:8080
```

## Progress

- ✅ Module 00: Complete
- ⏳ Module 01: In progress
- 📋 Module 02: Planned
- 📋 Module 03: Planned
