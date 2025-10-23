# Module 03 - HTTP API

Fortune-telling HTTP API server with JSON responses. Part of Road to Mercari Gopher Dojo training program.

## Exercise 00: Omikuji API

HTTP server that returns random fortune (omikuji) results in JSON format.

```bash
cd ex00
go build -o omikuji
./omikuji 8080
curl localhost:8080
```

**Response example**:
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

**Features**:
- 7 fortune types: Dai-kichi, Kichi, Chuu-kichi, Sho-kichi, Sue-kichi, Kyo, Dai-kyo
- Special New Year logic (Jan 1-3): Always returns Dai-kichi
- Testable time dependency using Clock interface
- HTTP handler tests using `httptest`

## Project Structure

```
ex00/
├── main.go              # CLI entry point
├── omikuji/             # Fortune logic package
│   ├── omikuji.go       # Fortune types and generation
│   ├── handler.go       # HTTP handler
│   ├── omikuji_test.go  # Unit tests
│   └── handler_test.go  # HTTP handler tests
├── go.mod
└── README.md
```

## Key Concepts

### HTTP Server

Basic HTTP server using `net/http`:

```go
handler := omikuji.NewHandler(omikuji.DefaultClock)
http.ListenAndServe(":8080", handler)
```

### JSON Encoding

Struct tags for JSON marshaling:

```go
type Response struct {
    Fortune string `json:"fortune"`
    Health  string `json:"health"`
}

json.NewEncoder(w).Encode(response)
```

### HTTP Handler Interface

Implementing `http.Handler`:

```go
type Handler struct {
    clock Clock
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(data)
}
```

### Testable Dependencies

Dependency injection for time:

```go
type Clock func() time.Time

// Production
handler := NewHandler(time.Now)

// Testing
mockClock := func() time.Time {
    return time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)
}
handler := NewHandler(mockClock)
```

### HTTP Testing

Using `httptest` for handler testing:

```go
req := httptest.NewRequest(http.MethodGet, "/", nil)
rec := httptest.NewRecorder()

handler.ServeHTTP(rec, req)

if rec.Code != http.StatusOK {
    t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
}

var response Response
json.NewDecoder(rec.Body).Decode(&response)
```

## Testing

```bash
# Run all tests with coverage
go test -cover ./...

# Expected output:
# ok  	ex00        	0.123s	coverage: 58.3% of statements
# ok  	ex00/omikuji	0.045s	coverage: 100.0% of statements
```
