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
	"os"
	"path/filepath"
	"strings"
)

func CleanName(fullPath, name string, isDir bool) (string, error) {
	var base, ext string

	if !isDir {
		if i := strings.LastIndexByte(name, '.'); i > 0 && i < len(name)-1 {
			base, ext = name[:i], name[i+1:]
		} else {
			base = name
		}
	} else {
		base = name
	}

	base = cleanASCII(base)
	ext = cleanASCII(ext)

	base = posixify(base)
	if ext != "" {
		ext = posixify(ext)
	}

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

	prefix := ""
	if flagDateMode == "mtime" {
		info, err := os.Stat(fullPath)
		if err != nil {
			return name, err
		}
		prefix = formatDate(info.ModTime(), flagDateStyle)
		if prefix != "" {
			prefix += flagDelim
		}
	}

	newName := prefix + base
	if ext != "" {
		newName += "." + ext
	}
	if newName == "" {
		return name, errors.New("empty result name")
	}

	if isWindowsReserved(strings.TrimSuffix(newName, filepath.Ext(newName))) {
		newName = "_" + newName
	}

	return newName, nil
}
