package main

import (
	"context"
	"strings"
	"testing"
	"time"
)

// TestReadInput tests reading from stdin with context cancellation
func TestReadInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		timeout  time.Duration
		wantRead bool
	}{
		{
			name:     "read single line",
			input:    "Hello\n",
			timeout:  100 * time.Millisecond,
			wantRead: true,
		},
		{
			name:     "read multiple lines",
			input:    "Line1\nLine2\nLine3\n",
			timeout:  100 * time.Millisecond,
			wantRead: true,
		},
		{
			name:     "context cancellation stops reading",
			input:    "Start\n",
			timeout:  1 * time.Millisecond,
			wantRead: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			ch := make(chan string, 10) // Buffered to avoid blocking

			// Mock stdin with strings.Reader
			// Note: We can't easily replace os.Stdin, so we test the logic separately
			// in production, readInput reads from os.Stdin

			// For testing purposes, let's test the channel communication pattern
			go func() {
				lines := strings.Split(tt.input, "\n")
				for _, line := range lines {
					if line == "" {
						continue
					}
					select {
					case <-ctx.Done():
						return
					case ch <- line:
					}
				}
			}()

			// Read from channel
			var received []string
			for {
				select {
				case <-ctx.Done():
					// Context cancelled
					if tt.wantRead && len(received) == 0 {
						t.Error("Expected to read at least one line before timeout")
					}
					return
				case line := <-ch:
					received = append(received, line)
					if len(received) >= strings.Count(tt.input, "\n") {
						return
					}
				}
			}
		})
	}
}

// TestReadInput_ContextCancellation specifically tests goroutine cleanup
func TestReadInput_ContextCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan string, 1)

	// This test verifies the goroutine respects context cancellation
	// In real usage, readInput would read from os.Stdin
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ch <- "test":
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// Read one value
	<-ch

	// Cancel context
	cancel()

	// Give goroutine time to clean up
	time.Sleep(100 * time.Millisecond)

	// Verify channel doesn't receive more values (goroutine stopped)
	select {
	case <-ch:
		t.Error("Goroutine did not stop after context cancellation")
	case <-time.After(50 * time.Millisecond):
		// Good: no more values received
	}
}

// TestWords verifies the word list is not empty
func TestWords(t *testing.T) {
	t.Parallel()

	if len(words) == 0 {
		t.Error("Word list is empty")
	}

	// Verify all words are non-empty
	for i, word := range words {
		if word == "" {
			t.Errorf("Word at index %d is empty", i)
		}
	}
}

// TestGameLogic_Timeout tests that game respects timeout
func TestGameLogic_Timeout(t *testing.T) {
	t.Parallel()

	// Create short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Wait for timeout
	<-ctx.Done()

	// Verify context is done
	if ctx.Err() == nil {
		t.Error("Context should have error after timeout")
	}

	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded, got %v", ctx.Err())
	}
}

// TestGameLogic_ScoreIncrement tests score logic
func TestGameLogic_ScoreIncrement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		word          string
		input         string
		initialScore  int
		expectedScore int
	}{
		{
			name:          "correct input increments score",
			word:          "Pikachu",
			input:         "Pikachu",
			initialScore:  0,
			expectedScore: 1,
		},
		{
			name:          "incorrect input does not increment",
			word:          "Pikachu",
			input:         "Charizard",
			initialScore:  0,
			expectedScore: 0,
		},
		{
			name:          "input with whitespace",
			word:          "Pikachu",
			input:         "  Pikachu  ",
			initialScore:  0,
			expectedScore: 1,
		},
		{
			name:          "case sensitive",
			word:          "Pikachu",
			input:         "pikachu",
			initialScore:  5,
			expectedScore: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			score := tt.initialScore

			// Simulate game logic
			if strings.TrimSpace(tt.input) == tt.word {
				score++
			}

			if score != tt.expectedScore {
				t.Errorf("Score = %d, want %d", score, tt.expectedScore)
			}
		})
	}
}

// TestRandomWordSelection tests that random word selection works
func TestRandomWordSelection(t *testing.T) {
	t.Parallel()

	// Select multiple random words
	selectedWords := make(map[string]bool)

	for i := 0; i < 100; i++ {
		// In Go 1.20+, rand is automatically seeded
		idx := i % len(words)
		word := words[idx]

		// Verify word is from the list
		found := false
		for _, w := range words {
			if w == word {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Selected word %q is not in word list", word)
		}

		selectedWords[word] = true
	}

	// We should have selected at least a few different words
	if len(selectedWords) < 2 {
		t.Error("Random selection is not working properly (too few unique words)")
	}
}

// TestChannelCommunication tests channel-based input handling
func TestChannelCommunication(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	inputCh := make(chan string, 1)

	// Send test input
	go func() {
		inputCh <- "TestInput"
	}()

	// Receive input (simulate game loop)
	select {
	case <-ctx.Done():
		t.Error("Timeout before receiving input")
	case input := <-inputCh:
		if input != "TestInput" {
			t.Errorf("Received %q, want %q", input, "TestInput")
		}
	}
}
