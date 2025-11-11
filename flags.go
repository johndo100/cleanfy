// flags.go
// --------
// Defines all command-line flags, usage text, and global configuration options for Cleanfy.
// Flags control behavior such as recursion, case transformation, date prefixing,
// JSON output, and automatic conflict resolution (--unique).
package main

import "flag"

var (
	flagDo, flagRecursive, flagQuiet, flagJSON, flagPretty, flagVersion bool
	flagUnique                                                          bool
	flagCase, flagDateMode, flagDateStyle, flagDelim                    string
)

func parseFlags() {
	flag.BoolVar(&flagDo, "d", false, "Perform actual renaming (default: preview only)")
	flag.BoolVar(&flagRecursive, "r", false, "Recurse into subdirectories")
	flag.BoolVar(&flagQuiet, "q", false, "Suppress normal output (show only errors)")
	flag.BoolVar(&flagJSON, "j", false, "Print results as JSON")
	flag.BoolVar(&flagPretty, "p", false, "Pretty-print JSON output (implies -j)")
	flag.BoolVar(&flagVersion, "v", false, "Show program version and exit")

	flag.StringVar(&flagCase, "case", "none", "Case transform: none|lower|upper|title")
	flag.StringVar(&flagDateMode, "date", "none", "Date prefix: none|mtime")
	flag.StringVar(&flagDateStyle, "date-style", "iso", "Date style: iso|compact|month|short|withtime")
	flag.StringVar(&flagDelim, "delim", "_", "Delimiter between date and name")

	// 🆕 Automatically handle name conflicts by adding suffixes (e.g., _2, _3)
	flag.BoolVar(&flagUnique, "unique", true, "Auto-resolve filename conflicts by adding numeric suffixes")

	flag.Parse()
}
