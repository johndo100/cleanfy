// Package cleanfy
// ----------------
// Main entry point of the Cleanfy CLI tool.
// Handles version display, flag parsing, directory walking, and output control.

package main

import (
	"fmt"
	"runtime"
)

// version is injected at build time via -ldflags
// Default value if not set during build
var version = "dev"

func main() {
	parseFlags()

	if flagVersion {
		fmt.Printf("cleanfy %s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
		return
	}

	results := walkAndRename()

	emitResults(results)
}
