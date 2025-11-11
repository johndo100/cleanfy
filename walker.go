// walker.go
// ----------
// Directory traversal logic for Cleanfy.
// Supports recursive and non-recursive file walking via WalkDir.

package main

import (
	"flag"
	"io/fs"
	"os"
	"path/filepath"
)

// walkAndRename walks through all target paths (files/directories)
// and processes them with processOne().
// It supports both recursive (-r) and non-recursive modes.
func walkAndRename() []Result {
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	var results []Result
	for _, root := range args {
		info, err := os.Stat(root)
		if err != nil {
			results = append(results, Result{Path: root, Error: err.Error()})
			continue
		}

		if info.IsDir() {
			if flagRecursive {
				filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						results = append(results, Result{Path: path, Error: err.Error()})
						return nil
					}
					info, err := d.Info()
					if err != nil {
						results = append(results, Result{Path: path, Error: err.Error()})
						return nil
					}
					results = append(results, processOne(path, info))
					return nil
				})
			} else {
				entries, err := os.ReadDir(root)
				if err != nil {
					results = append(results, Result{Path: root, Error: err.Error()})
					continue
				}
				for _, e := range entries {
					p := filepath.Join(root, e.Name())
					info, err := e.Info()
					if err != nil {
						results = append(results, Result{Path: p, Error: err.Error()})
						continue
					}
					results = append(results, processOne(p, info))
				}
			}
		} else {
			results = append(results, processOne(root, info))
		}
	}
	return results
}
