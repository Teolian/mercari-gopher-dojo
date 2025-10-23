package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ex00/omikuji"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <port>\n", os.Args[0])
		os.Exit(1)
	}

	port := os.Args[1]
	addr := ":" + port

	handler := omikuji.NewHandler(omikuji.DefaultClock)

	fmt.Printf("Starting omikuji server on port %s\n", port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
