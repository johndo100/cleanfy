// posix.go
// ---------
// Ensures filenames are POSIX-safe by removing illegal characters,
// collapsing duplicates, and trimming leading/trailing dots, underscores, and dashes.

package main

import (
	"regexp"
	"strings"
)

var (
	reDisallowed  = regexp.MustCompile(`[^A-Za-z0-9._-]+`)
	reMultiUnders = regexp.MustCompile(`_+`)
	reMultiDashes = regexp.MustCompile(`-+`)
)

func posixify(s string) string {
	if s == "" {
		return "_"
	}
	s = strings.ReplaceAll(s, " ", "_")
	s = reDisallowed.ReplaceAllString(s, "_")
	s = reMultiUnders.ReplaceAllString(s, "_")
	s = reMultiDashes.ReplaceAllString(s, "-")
	s = strings.Trim(s, "._-")
	if s == "" {
		return "_"
	}
	return s
}
