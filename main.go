package main

import (
	"fmt"
	"io"
	"os"
)

// cat copies data from r to w using io.Copy.
// This function demonstrates the power of io.Reader and io.Writer interfaces,
// allowing any source and destination that implement these interfaces.
func cat(r io.Reader, w io.Writer) error {
	_, err := io.Copy(w, r)
	return err
}

func main() {
	var exitCode int

	// No arguments: read from stdin
	if len(os.Args) == 1 {
		if err := cat(os.Stdin, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "ft_cat: stdin: %v\n", err)
			exitCode = 1
		}
		os.Exit(exitCode)
	}

	// With arguments: read each file
	for _, filename := range os.Args[1:] {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ft_cat: %s: %v\n", filename, err)
			exitCode = 1
			continue
		}

		if err := cat(f, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "ft_cat: %s: %v\n", filename, err)
			exitCode = 1
		}

		f.Close()
	}

	os.Exit(exitCode)
}
