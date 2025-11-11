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
	case "lower":
		base = strings.ToLower(base)
		ext = strings.ToLower(ext)
	case "upper":
		base = strings.ToUpper(base)
		ext = strings.ToUpper(ext)
	case "title":
		base = toTitle(base)
	}

	// Add date prefix if requested
	prefix := ""
	if flagDateMode != "" {
		prefix = getDatePrefix(fullPath, flagDateMode, flagDateFormat)
		if prefix != "" {
			prefix += "_"
		}
	}

	// Reconstruct name
	newName := prefix + base
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
