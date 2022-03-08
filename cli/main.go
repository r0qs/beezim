package main

import (
	"fmt"
	"os"

	"github.com/r0qs/beezim/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
