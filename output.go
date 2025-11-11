// output.go
// ----------
// Handles all output formatting for Cleanfy.
// Supports JSON and plain-text output, with options for quiet, pretty, and error-only modes.
// Automatically highlights auto-renamed files (via --unique).

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// emitResults outputs a list of results in either JSON or plain-text format,
// depending on the selected flags.
func emitResults(results []Result) {
	if flagJSON || flagPretty {
		emitJSON(results, flagPretty)
		return
	}

	if flagQuiet {
		return
	}

	w := bufio.NewWriter(os.Stdout)
	for _, r := range results {
		printResult(w, r)
	}
	w.Flush()
}

// emitJSON prints results as JSON, optionally pretty-printed.
func emitJSON(v any, pretty bool) {
	var b []byte
	var err error
	if pretty {
		b, err = json.MarshalIndent(v, "", "  ")
	} else {
		b, err = json.Marshal(v)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "json error: %v\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(b)
	os.Stdout.Write([]byte("\n"))
}

// printResult prints one result entry in human-readable text form.
// Handles renamed, auto-renamed, and error cases.
func printResult(w *bufio.Writer, r Result) {
	if r.Error != "" {
		fmt.Fprintf(w, "ERR     %s : %s\n", r.Path, r.Error)
		return
	}
	if r.NewName == "" || r.NewName == r.OldName {
		fmt.Fprintf(w, "OK      %s\n", r.OldName)
		return
	}

	// Differentiate between normal rename and auto-rename
	if r.AutoRenamed {
		fmt.Fprintf(w, "RENAME* %s -> %s   (auto-resolved)\n", r.OldName, r.NewName)
	} else {
		fmt.Fprintf(w, "RENAME  %s -> %s\n", r.OldName, r.NewName)
	}
}
