package main

import (
	"fmt"
	"os"

	"github.com/rishad1234/term-video-transcoder/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
