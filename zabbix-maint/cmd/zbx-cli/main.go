package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// TODO: implement CLI entry point
	// 1. Parse command line arguments
	// 2. Route to interactive mode, one-shot mode, or config mode
	return nil
}
