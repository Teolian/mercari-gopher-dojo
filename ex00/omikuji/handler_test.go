package omikuji

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_ServeHTTP_StatusCode(t *testing.T) {
	t.Parallel()

	mockClock := func() time.Time {
		return time.Date(2025, time.July, 15, 12, 0, 0, 0, time.UTC)
	}

	handler := NewHandler(mockClock)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("ServeHTTP() status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestHandler_ServeHTTP_ContentType(t *testing.T) {
	t.Parallel()

	mockClock := func() time.Time {
		return time.Date(2025, time.July, 15, 12, 0, 0, 0, time.UTC)
	}

	handler := NewHandler(mockClock)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("ServeHTTP() Content-Type = %q, want %q", contentType, "application/json")
	}
}

func TestHandler_ServeHTTP_JSONStructure(t *testing.T) {
	t.Parallel()

	mockClock := func() time.Time {
		return time.Date(2025, time.July, 15, 12, 0, 0, 0, time.UTC)
	}

	handler := NewHandler(mockClock)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var resp Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Verify fortune is one of the valid values
	validFortune := false
	for _, f := range AllFortunes {
		if resp.Fortune == f {
			validFortune = true
			break
		}
	}
	if !validFortune {
		t.Errorf("ServeHTTP() fortune = %v, not in AllFortunes", resp.Fortune)
	}

	// Verify additional fields are present (at least 1 is required by spec)
	if resp.Health == "" {
		t.Error("ServeHTTP() health field is empty")
	}
	if resp.Residence == "" {
		t.Error("ServeHTTP() residence field is empty")
	}
	if resp.Travel == "" {
		t.Error("ServeHTTP() travel field is empty")
	}
	if resp.Study == "" {
		t.Error("ServeHTTP() study field is empty")
	}
	if resp.Love == "" {
		t.Error("ServeHTTP() love field is empty")
	}
}

func TestHandler_ServeHTTP_NewYearLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mockTime time.Time
		want     Fortune
	}{
		{
			name:     "January 1st returns Dai-kichi",
			mockTime: time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC),
			want:     DaiKichi,
		},
		{
			name:     "January 2nd returns Dai-kichi",
			mockTime: time.Date(2025, time.January, 2, 12, 0, 0, 0, time.UTC),
			want:     DaiKichi,
		},
		{
			name:     "January 3rd returns Dai-kichi",
			mockTime: time.Date(2025, time.January, 3, 12, 0, 0, 0, time.UTC),
			want:     DaiKichi,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClock := func() time.Time {
				return tt.mockTime
			}

			handler := NewHandler(mockClock)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			var resp Response
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			if resp.Fortune != tt.want {
				t.Errorf("ServeHTTP() fortune = %v, want %v", resp.Fortune, tt.want)
			}
		})
	}
}
