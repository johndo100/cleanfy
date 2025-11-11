// flags.go
// --------
// Defines all command-line flags, usage text, and global configuration options for Cleanfy.
// Flags control behavior such as recursion, case transformation, date prefixing,
// JSON output, and automatic conflict resolution (--unique).
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flagDo, flagRecursive, flagQuiet, flagDotfiles, flagJSON, flagVersion bool
	flagCase, flagDateMode, flagDateFormat                                string
)

// Always enable --unique behavior
const flagUnique = true

func parseFlags() {
	// // Core behavior
	// flag.BoolVar(&flagDo, "d", false, "Perform actual renaming (default: preview only)")
	// flag.BoolVar(&flagDo, "do", false, "Perform actual renaming (default: preview only)")

	// flag.BoolVar(&flagRecursive, "r", false, "Recurse into subdirectories")
	// flag.BoolVar(&flagRecursive, "recursive", false, "Recurse into subdirectories")

	// flag.BoolVar(&flagQuiet, "q", false, "Suppress normal output (show only errors)")
	// flag.BoolVar(&flagQuiet, "quiet", false, "Suppress normal output (show only errors)")

	// flag.BoolVar(&flagDotfiles, "a", false, "Include hidden files (starting with .) when processing")
	// flag.BoolVar(&flagDotfiles, "dotfiles", false, "Include hidden files (starting with .) when processing")

	// flag.BoolVar(&flagJSON, "j", false, "Show output in JSON format")
	// flag.BoolVar(&flagJSON, "json", false, "Show output in JSON format")

	// flag.BoolVar(&flagVersion, "v", false, "Show program version and exit")
	// flag.BoolVar(&flagVersion, "version", false, "Show program version and exit")

	// // Options
	// flag.StringVar(&flagCase, "c", "lower", "Case transform: lower|upper|title (default: lower)")
	// flag.StringVar(&flagCase, "case", "lower", "Case transform: lower|upper|title (default: lower)")

	// flag.StringVar(&flagDateMode, "t", "mtime", "Add date prefix: mtime (last modified, default) | now (current date)")
	// flag.StringVar(&flagDateMode, "date", "mtime", "Add date prefix: mtime (last modified, default) | now (current date)")

	// flag.StringVar(&flagDateFormat, "f", "2006-01-02", "Date format (default ISO 8601)")
	// flag.StringVar(&flagDateFormat, "date-format", "2006-01-02", "Date format (default ISO 8601)")

	// Core
	flag.BoolVar(&flagDo, "d", false, "Perform actual renaming (default: preview only) [--do]")
	flag.BoolVar(&flagRecursive, "r", false, "Recurse into subdirectories [--recursive]")
	flag.BoolVar(&flagQuiet, "q", false, "Suppress normal output (show only errors) [--quiet]")
	flag.BoolVar(&flagJSON, "j", false, "Show output in JSON format [--json]")
	flag.BoolVar(&flagVersion, "v", false, "Show program version and exit [--version]")

	// File handling
	flag.BoolVar(&flagDotfiles, "a", false, "Include hidden files (starting with .) [--dotfiles]")

	// Options
	flag.StringVar(&flagCase, "c", "lower", "Case transform: lower|upper|title (default: lower) [--case]")
	flag.StringVar(&flagDateMode, "t", "mtime", "Add date prefix: mtime (last modified, default) | now (current date) [--date]")
	flag.StringVar(&flagDateFormat, "f", "2006-01-02", "Date format (default ISO 8601) [--date-format]")

	// Parse flags
	flag.Parse()

	// Validate -case
	switch flagCase {
	case "lower", "upper", "title":
		// valid
	default:
		fmt.Fprintln(os.Stderr, "Invalid -case value. Use one of: lower | upper | title")
		os.Exit(2)
	}

	// Validate -date
	switch flagDateMode {
	case "mtime", "now":
		// valid
	default:
		fmt.Fprintln(os.Stderr, "Invalid -date value. Use one of: mtime | now")
		os.Exit(2)
	}

}
