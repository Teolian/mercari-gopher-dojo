# Exercise 00: Fortune API (おみくじ)

HTTP API server that returns random fortune (omikuji) in JSON format.

## Build and Run

```bash
cd ex00
go build -o omikuji
./omikuji 8080
```

## Usage

```bash
# Start server
./omikuji 4242 &

# Get fortune
curl localhost:4242
```

## Response Format

```json
{
  "fortune": "Kichi",
  "health": "You will fully recover, but stay attentive after you do.",
  "residence": "You will have good fortune with a new house.",
  "travel": "When traveling, you may find something to treasure.",
  "study": "Things will be better. It may be worth aiming for a school in a different area.",
  "love": "The person you are looking for is very close to you."
}
```

## Fortune Types

- `Dai-kichi` - Great blessing
- `Kichi` - Blessing
- `Chuu-kichi` - Middle blessing
- `Sho-kichi` - Small blessing
- `Sue-kichi` - Future blessing
- `Kyo` - Curse
- `Dai-kyo` - Great curse

## Special Behavior

During New Year (January 1-3), the API always returns `Dai-kichi`.

## Testing

```bash
# Run all tests with coverage
go test -cover ./...

# Run tests in specific package
cd omikuji && go test -v
```

## Implementation Details

### Package Structure

- `main.go` - CLI entry point, server setup
- `omikuji/omikuji.go` - Fortune logic and types
- `omikuji/handler.go` - HTTP handler implementation
- `omikuji/omikuji_test.go` - Unit tests for fortune generation
- `omikuji/handler_test.go` - HTTP handler tests using `httptest`

### Testable Time Dependency

The implementation uses dependency injection for time to enable testing:

```go
type Clock func() time.Time

func GenerateFortune(clock Clock) Response {
    now := clock()
    // Use now for date logic
}
```

This allows tests to mock the current time and verify New Year behavior without waiting for January 1-3.

### HTTP Testing

Uses `net/http/httptest` package for HTTP handler testing:

```go
req := httptest.NewRequest(http.MethodGet, "/", nil)
rec := httptest.NewRecorder()
handler.ServeHTTP(rec, req)

// Verify status code, headers, JSON structure
```
