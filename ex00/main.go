package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Word list for the typing game
var words = []string{
	"Water", "Bulbasaur", "Pokemon", "Flamethrower",
	"Thunder", "Pikachu", "Charizard", "Squirtle",
	"Jigglypuff", "Mewtwo", "Eevee", "Snorlax",
	"Gengar", "Dragonite", "Mew", "Alakazam",
}

func main() {
	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Create context with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Channel for user input
	inputCh := make(chan string)

	// Goroutine to read from stdin
	go readInput(inputCh)

	score := 0

	// Game loop
	for {
		// Pick random word
		word := words[rand.Intn(len(words))]
		fmt.Println(word)
		fmt.Print("-> ")

		// Wait for input or timeout
		select {
		case <-ctx.Done():
			// Timeout reached
			fmt.Printf("\nTime's up! Score: %d\n", score)
			return

		case input := <-inputCh:
			// Check if input matches
			if strings.TrimSpace(input) == word {
				score++
			}
		}
	}
}

// readInput reads from stdin and sends to channel
// Runs in a separate goroutine to avoid blocking
func readInput(ch chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ch <- scanner.Text()
	}
}
