// rename.go
// ----------
// Core renaming logic for Cleanfy.
// Handles per-file name normalization, conflict resolution (--unique),
// dotfile preservation, and performing actual rename operations on disk.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// processOne applies the rename process to a single file or directory.
// It returns a Result struct describing the outcome.
func processOne(path string, info os.FileInfo) Result {
	name := info.Name()
	isDir := info.IsDir()

	// 🧩 Dotfile handling
	// By default, Cleanfy skips hidden files (starting with '.').
	// Users can override this behavior using --dotfiles.
	if strings.HasPrefix(name, ".") && len(name) > 1 && !flagDotfiles {
		return Result{
			Path:       path,
			OldName:    name,
			NewName:    name,
			IsDir:      isDir,
			Renamed:    false,
			WasSkipped: true,
		}
	}

	newName, err := CleanName(path, name, isDir)
	if err != nil {
		return Result{Path: path, OldName: name, Error: err.Error(), IsDir: isDir}
	}

	// No change
	if newName == name {
		return Result{Path: path, OldName: name, NewName: newName, IsDir: isDir}
	}

	// Dry-run mode
	if !flagDo {
		return Result{Path: path, OldName: name, NewName: newName, IsDir: isDir}
	}

	newFull := filepath.Join(filepath.Dir(path), newName)

	// Handle existing destination
	if _, err := os.Stat(newFull); err == nil {
		if flagUnique {
			// Automatically generate a unique name if --unique is enabled
			newFull, newName = makeUnique(filepath.Dir(path), newName)
		} else {
			return Result{
				Path: path, OldName: name, NewName: newName, IsDir: isDir,
				Error: "destination exists",
			}
		}
	}

	// Perform the rename
	if err := os.Rename(path, newFull); err != nil {
		return Result{
			Path: path, OldName: name, NewName: newName, IsDir: isDir,
			Error: err.Error(),
		}
	}

	// Detect if suffix (_2, _3, etc.) was added
	autoRenamed := strings.Contains(newName, "_") &&
		strings.HasPrefix(filepath.Base(newName), strings.Split(filepath.Base(name), ".")[0])

	return Result{
		Path:        newFull,
		OldName:     name,
		NewName:     newName,
		IsDir:       isDir,
		Renamed:     true,
		AutoRenamed: autoRenamed,
	}
}

// makeUnique generates a non-conflicting name by appending a numeric suffix.
// Example: "file.txt" → "file_2.txt" → "file_3.txt" → ...
func makeUnique(dir, name string) (string, string) {
	base := name
	ext := ""
	if i := strings.LastIndexByte(name, '.'); i > 0 && i < len(name)-1 {
		base = name[:i]
		ext = name[i:]
	}

	for n := 2; ; n++ {
		candidate := fmt.Sprintf("%s_%d%s", base, n, ext)
		fullPath := filepath.Join(dir, candidate)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fullPath, candidate
		}
	}
}
