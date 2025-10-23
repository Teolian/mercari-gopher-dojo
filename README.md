# Module 02 - Concurrency

Go concurrency patterns using goroutines, channels, and context. Part of Road to Mercari Gopher Dojo training program.

## Exercises

### Exercise 00: Typing Game
**Directory**: `ex00/`

Terminal typing game with 30-second timer.

```bash
cd ex00
go build -o typing_game
./typing_game
```

**Concepts**:
- Goroutines for non-blocking I/O
- Channels for communication
- Context with timeout
- Select statement

---

### Exercise 01: Parallel Downloader
**Directory**: `ex01/`

HTTP file downloader using parallel range requests.

```bash
cd ex01
go build -o download
./download https://example.com/file.zip
```

**Concepts**:
- HTTP Range requests
- errgroup for error handling
- Parallel goroutines
- Context cancellation

## Project Structure

```
├── ex00/              # Typing game
│   ├── main.go
│   ├── go.mod
│   └── README.md
└── ex01/              # Parallel downloader
    ├── main.go
    ├── go.mod
    └── README.md
```

## Key Concepts

### Goroutines
Lightweight threads managed by Go runtime.

```go
go func() {
    // Runs concurrently
}()
```

### Channels
Communication between goroutines.

```go
ch := make(chan string)
go func() { ch <- "message" }()
msg := <-ch
```

### Select
Multiplexing over channels.

```go
select {
case msg := <-ch1:
    // Handle ch1
case <-time.After(timeout):
    // Timeout
}
```

### Context
Cancellation and deadlines.

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

select {
case <-ctx.Done():
    // Timeout or cancelled
}
```

### errgroup
Managing multiple goroutines with error handling.

```go
g, ctx := errgroup.WithContext(ctx)

g.Go(func() error {
    // Concurrent work
    return nil
})

if err := g.Wait(); err != nil {
    // Handle first error
}
```
