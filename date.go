// date.go
// --------
// Handles date prefixing for Cleanfy.
// Supports modes: mtime (last modified) and now (current time).
// Uses Go's time formatting layout, defaulting to ISO 8601 (2006-01-02).

package main

import (
	"os"
	"time"
)

// getDatePrefix returns a formatted date string for a given file
// based on the selected mode: "mtime" or "now".
func getDatePrefix(path string, mode string, layout string) string {
	var t time.Time

	switch mode {
	case "mtime":
		info, err := os.Stat(path)
		if err == nil {
			t = info.ModTime()
		}
	case "now":
		t = time.Now()
	default:
		return "" // no -date flag provided → skip prefix
	}

	if t.IsZero() {
		t = time.Now()
	}

	return t.Format(layout)
}
