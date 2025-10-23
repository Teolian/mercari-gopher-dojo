package omikuji

import (
	"testing"
	"time"
)

func TestGenerateFortune_NewYear(t *testing.T) {
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

			resp := GenerateFortune(mockClock)

			if resp.Fortune != tt.want {
				t.Errorf("GenerateFortune() fortune = %v, want %v", resp.Fortune, tt.want)
			}

			// Verify required fields are present
			if resp.Health == "" {
				t.Error("GenerateFortune() health field is empty")
			}
			if resp.Residence == "" {
				t.Error("GenerateFortune() residence field is empty")
			}
			if resp.Travel == "" {
				t.Error("GenerateFortune() travel field is empty")
			}
			if resp.Study == "" {
				t.Error("GenerateFortune() study field is empty")
			}
			if resp.Love == "" {
				t.Error("GenerateFortune() love field is empty")
			}
		})
	}
}

func TestGenerateFortune_NonNewYear(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mockTime time.Time
	}{
		{
			name:     "January 4th can return any fortune",
			mockTime: time.Date(2025, time.January, 4, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "December 31st can return any fortune",
			mockTime: time.Date(2025, time.December, 31, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "July 15th can return any fortune",
			mockTime: time.Date(2025, time.July, 15, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClock := func() time.Time {
				return tt.mockTime
			}

			resp := GenerateFortune(mockClock)

			// Verify fortune is one of the valid values
			validFortune := false
			for _, f := range AllFortunes {
				if resp.Fortune == f {
					validFortune = true
					break
				}
			}

			if !validFortune {
				t.Errorf("GenerateFortune() fortune = %v, not in AllFortunes", resp.Fortune)
			}

			// Verify required fields are present
			if resp.Health == "" {
				t.Error("GenerateFortune() health field is empty")
			}
		})
	}
}

func TestAllFortunes_Complete(t *testing.T) {
	t.Parallel()

	expected := []Fortune{DaiKichi, Kichi, ChuuKichi, ShoKichi, SueKichi, Kyo, DaiKyo}

	if len(AllFortunes) != len(expected) {
		t.Errorf("AllFortunes length = %d, want %d", len(AllFortunes), len(expected))
	}

	for i, fortune := range expected {
		if AllFortunes[i] != fortune {
			t.Errorf("AllFortunes[%d] = %v, want %v", i, AllFortunes[i], fortune)
		}
	}
}
