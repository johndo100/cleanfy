// cleaner.go
// -----------
// Filename normalization pipeline for Cleanfy.
// Steps:
// 1. Split extension
// 2. Normalize to ASCII (NFKD)
// 3. POSIX filtering
// 4. Case transform
// 5. Optional date prefix
// 6. Reserved name protection

package main

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

func CleanName(fullPath, name string, isDir bool) (string, error) {
	// Split name into base and extension
	var base, ext string

	// Split extension
	if !isDir {
		if i := strings.LastIndexByte(name, '.'); i > 0 && i < len(name)-1 {
			base, ext = name[:i], name[i+1:]
		} else {
			base = name
		}
	} else {
		base = name
	}

	// Normalize to ASCII
	base = cleanASCII(base)
	ext = cleanASCII(ext)

	// POSIX filtering
	base = posixify(base)
	if ext != "" {
		ext = posixify(ext)
	}

	// Apply case transformation
	switch strings.ToLower(flagCase) {
	case "", "none":
	// keep original case
	case "lower":
		base = strings.ToLower(base)
		ext = strings.ToLower(ext)
	case "upper":
		base = strings.ToUpper(base)
		ext = strings.ToUpper(ext)
	case "title":
		base = toTitle(base)
	}

	// precompile the regex once (top of file or as a package-level var)
	var datePrefixRegex = regexp.MustCompile(`^(?:\d{4}[-_.\/]?\d{2}[-_.\/]?\d{2}|\d{6})[_\-\.]`)

	// 🧩 Add date prefix only when explicitly requested (-t or -date)
	if flagDateMode != "" && !datePrefixRegex.MatchString(base) {
		if prefix := getDatePrefix(fullPath, flagDateMode, flagDateFormat); prefix != "" {
			base = prefix + "_" + base
		}
	}

	// Reconstruct name
	newName := base
	if ext != "" {
		newName += "." + ext
	}
	if newName == "" {
		return name, errors.New("empty result name")
	}

	// Prevent Windows reserved names
	if isWindowsReserved(strings.TrimSuffix(newName, filepath.Ext(newName))) {
		newName = "_" + newName
	}

	// Return final name
	return newName, nil
}
