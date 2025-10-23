# Typing Game

Terminal-based typing game with 30-second timer. Part of Road to Mercari Gopher Dojo Module 02, Exercise 00.

## Description

Tests typing speed and accuracy. Random words appear and the player must type them correctly within 30 seconds.

## Build

```bash
go build -o typing_game
```

## Usage

```bash
./typing_game
```

Game will display random words. Type each word exactly and press Enter.

Example session:
```
Water
-> Water
Bulbasaur
-> Bulbasaur
Pokemon
-> Pokemon
Flamethrower
-> Thunder
Time's up! Score: 3
```

Score counts only correct matches.

## Implementation

**Concurrency concepts:**

1. **Goroutine** - separate goroutine reads stdin to avoid blocking main game loop
2. **Channel** - `inputCh` transfers user input from reader goroutine to main
3. **Context** - `context.WithTimeout` creates 30-second deadline
4. **Select** - multiplexes between user input and timeout

**Key pattern:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

select {
case <-ctx.Done():
    // Timeout
case input := <-inputCh:
    // User input
}
```

This demonstrates:
- Non-blocking I/O with goroutines
- Channel-based communication
- Timeout handling with context
- Select statement for multiplexing

## Project Structure

```
├── main.go         # Game implementation
└── go.mod          # Module definition
```
