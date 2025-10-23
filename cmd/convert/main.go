package main

import (
	"flag"
	"fmt"
	"os"

	"road-to-mercari-gopher-dojo-00/imgconv"
)

func main() {
	inFmt := flag.String("i", "jpg", "input format (jpg|png|gif)")
	outFmt := flag.String("o", "png", "output format (jpg|png|gif)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("error: invalid argument")
		os.Exit(1)
	}
	dir := flag.Arg(0)

	conv, err := imgconv.NewConverter(dir, *inFmt, *outFmt)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// Run conversion; returns number of per-file errors
	if nErr := conv.Run(); nErr > 0 {
		fmt.Fprintf(os.Stderr, "error: %d file(s) failed\n", nErr)
		os.Exit(1)
	}
}
