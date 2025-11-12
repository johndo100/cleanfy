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
	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "cleanfy — smart batch file renamer\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  cleanfy [options] [targets...]\n\n")

		fmt.Fprintf(os.Stderr, "Core:\n")
		fmt.Fprintf(os.Stderr, "  -x, --execute              Perform actual renaming (default: preview only)\n")
		fmt.Fprintf(os.Stderr, "  -r, --recursive            Recurse into subdirectories\n")
		fmt.Fprintf(os.Stderr, "  -q, --quiet                Suppress normal output\n")
		fmt.Fprintf(os.Stderr, "  -j, --json                 Show output in JSON format\n")
		fmt.Fprintf(os.Stderr, "  -v, --version              Show program version and exit\n")
		fmt.Fprintf(os.Stderr, "  -a, --dotfiles             Include hidden files (starting with .)\n\n")

		fmt.Fprintf(os.Stderr, "Rename modifiers:\n")
		fmt.Fprintf(os.Stderr, "  --case=value               Case transform: none|lower|upper|title\n")
		fmt.Fprintf(os.Stderr, "  --date=value               Add date prefix: mtime|now\n")
		fmt.Fprintf(os.Stderr, "  --date-format=value        Go time layout, e.g. 20060102 (with --date)\n\n")

		fmt.Fprintf(os.Stderr, "Notes:\n")
		fmt.Fprintf(os.Stderr, "  • All value flags must use the = form (e.g. --date=now or -c=lower)\n")
		fmt.Fprintf(os.Stderr, "  • Use '--' to separate flags from filenames starting with '-'\n")
		fmt.Fprintf(os.Stderr, "  • No implicit defaults for --case or --date — must be explicitly set\n")
	}

	// Core
	flag.BoolVar(&flagDo, "x", false, "Perform actual renaming (default: preview only)")
	flag.BoolVar(&flagDo, "execute", false, "Alias for -x")
	flag.BoolVar(&flagRecursive, "r", false, "Recurse into subdirectories")
	flag.BoolVar(&flagRecursive, "recursive", false, "Alias for -r")
	flag.BoolVar(&flagQuiet, "q", false, "Suppress normal output (show only errors)")
	flag.BoolVar(&flagQuiet, "quiet", false, "Alias for -q")
	flag.BoolVar(&flagJSON, "j", false, "Show output in JSON format")
	flag.BoolVar(&flagJSON, "json", false, "Alias for -j")
	flag.BoolVar(&flagVersion, "v", false, "Show program version and exit")
	flag.BoolVar(&flagVersion, "version", false, "Alias for -v")

	// File handling
	flag.BoolVar(&flagDotfiles, "a", false, "Include hidden files (starting with .)")
	flag.BoolVar(&flagDotfiles, "dotfiles", false, "Alias for -a")

	// Options
	flag.StringVar(&flagCase, "c", "", "Case transform: none|lower|upper|title")
	flag.StringVar(&flagCase, "case", "", "Alias for -c")
	flag.StringVar(&flagDateMode, "d", "", "Date prefix mode: mtime|now")
	flag.StringVar(&flagDateMode, "date", "", "Alias for -d")
	flag.StringVar(&flagDateFormat, "f", "2006-01-02", "Date format (default: 2006-01-02)")
	flag.StringVar(&flagDateFormat, "date-format", "2006-01-02", "Alias for -f")

	// Parse flags
	flag.Parse()

	// Validate target position — all flags must come before targets
	for i, arg := range os.Args[1:] {
		if len(arg) > 0 && arg[0] != '-' {
			// found first target
			for _, rest := range os.Args[i+2:] {
				if len(rest) > 0 && rest[0] == '-' {
					fmt.Fprintf(os.Stderr, "❌ Error: flags must appear before targets (got '%s' after '%s')\n\n", rest, arg)
					flag.Usage()
					os.Exit(2)
				}
			}
			break
		}
	}

	// Ensure at least one target is specified
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "❌ Error: no target specified.")
		fmt.Fprintln(os.Stderr, "Hint: use '.' to scan the current directory.")
		fmt.Fprintln(os.Stderr)
		flag.Usage()
		os.Exit(2)
	}

	// Validate --case (only if provided)
	if flagCase != "" {
		switch flagCase {
		case "none", "lower", "upper", "title":
			// valid
		default:
			fmt.Fprintln(os.Stderr, "❌ Invalid --case value. Use one of: none | lower | upper | title")
			fmt.Fprintln(os.Stderr)
			flag.Usage()
			os.Exit(2)
		}
	}

	// Validate --date (only if provided)
	if flagDateMode != "" {
		switch flagDateMode {
		case "mtime", "now":
			// valid
		default:
			fmt.Fprintln(os.Stderr, "❌ Invalid --date value. Use one of: mtime | now")
			fmt.Fprintln(os.Stderr)
			flag.Usage()
			os.Exit(2)
		}
	}
}
