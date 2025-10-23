# Parallel Downloader

HTTP file downloader using parallel range requests. Part of Road to Mercari Gopher Dojo Module 02, Exercise 01.

## Description

Downloads files faster by splitting them into parts and downloading in parallel using HTTP Range requests.

## Build

```bash
go mod tidy
go build -o download
```

## Usage

```bash
./download <URL>
```

Example:
```bash
./download https://example.com/file.zip
```

File will be saved with original filename from URL.

## Implementation

**Concurrency concepts:**

1. **HTTP Range Requests** - `Range: bytes=0-1023` header for partial content
2. **errgroup** - manages parallel goroutines with error propagation
3. **Context** - cancellation propagates to all goroutines on first error
4. **Goroutines** - one per file part (default: 4 parts)

**Algorithm:**

1. HEAD request checks `Accept-Ranges: bytes` and gets file size
2. Calculate part ranges: `[0-N], [N+1-2N], ...`
3. Launch goroutines (via errgroup) to download each part
4. If any part fails, context cancels all others
5. Merge parts sequentially into final file
6. Clean up temporary files

**Key pattern:**
```go
g, ctx := errgroup.WithContext(ctx)

for i := 0; i < numParts; i++ {
    g.Go(func() error {
        return downloadPart(ctx, ...)
    })
}

// Wait for all, stop on first error
if err := g.Wait(); err != nil {
    return err
}
```

## Features

- Parallel downloads (4 parts by default)
- Automatic fallback to simple download if range not supported
- Error handling with automatic cancellation
- Temporary file cleanup
- Progress indication

## Project Structure

```
├── main.go    # Downloader implementation
└── go.mod     # Module with errgroup dependency
```
