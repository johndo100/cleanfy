// rename_test.go
// ---------------
// Unit tests for rename functionality.
// Tests cover makeUnique collision handling and file-level normalization.

package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestMakeUnique tests the unique filename generation for collision handling.
func TestMakeUnique(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		dir        string
		filename   string
		existing   []string
		wantBase   string
		wantSuffix string
	}{
		{
			name:       "no-collision",
			dir:        tmpDir,
			filename:   "newfile.txt",
			existing:   []string{},
			wantBase:   "newfile_2.txt",
			wantSuffix: "_2",
		},
		{
			name:       "one-collision",
			dir:        tmpDir,
			filename:   "file.txt",
			existing:   []string{"file.txt"},
			wantBase:   "file_2.txt",
			wantSuffix: "_2",
		},
		{
			name:       "multiple-collisions",
			dir:        tmpDir,
			filename:   "file.txt",
			existing:   []string{"file.txt", "file_2.txt", "file_3.txt"},
			wantBase:   "file_4.txt",
			wantSuffix: "_4",
		},
		{
			name:       "no-extension",
			dir:        tmpDir,
			filename:   "README",
			existing:   []string{"README"},
			wantBase:   "README_2",
			wantSuffix: "_2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create existing files
			for _, fname := range tt.existing {
				fpath := filepath.Join(tt.dir, fname)
				if err := os.WriteFile(fpath, []byte(""), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
				defer os.Remove(fpath)
			}

			fullPath, newName := makeUnique(tt.dir, tt.filename)

			// Check the returned name
			if newName != tt.wantBase {
				t.Errorf("makeUnique() newName = %q, want %q", newName, tt.wantBase)
			}

			// Verify the full path matches
			expectedFull := filepath.Join(tt.dir, tt.wantBase)
			if fullPath != expectedFull {
				t.Errorf("makeUnique() fullPath = %q, want %q", fullPath, expectedFull)
			}

			// Verify the generated file doesn't exist (yet)
			if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
				t.Errorf("makeUnique() generated path already exists: %v", err)
			}
		})
	}
}

// TestProcessOneDryRun tests the dry-run mode of processOne (no actual renaming).
func TestProcessOneDryRun(t *testing.T) {
	// Save original flags
	oldFlagDo := flagDo
	oldFlagCase := flagCase
	oldFlagDateMode := flagDateMode
	oldFlagDotfiles := flagDotfiles

	defer func() {
		flagDo = oldFlagDo
		flagCase = oldFlagCase
		flagDateMode = oldFlagDateMode
		flagDotfiles = oldFlagDotfiles
	}()

	// Set up test flags for dry-run
	flagDo = false
	flagCase = "lower"
	flagDateMode = ""
	flagDotfiles = false

	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		filename   string
		isDir      bool
		wantNormal string // expected old name (unchanged)
		wantNew    string // expected new name after cleaning
	}{
		{
			name:       "file-rename",
			filename:   "MyFile.txt",
			isDir:      false,
			wantNormal: "MyFile.txt",
			wantNew:    "myfile.txt",
		},
		{
			name:       "directory",
			filename:   "MyDir",
			isDir:      true,
			wantNormal: "MyDir",
			wantNew:    "mydir",
		},
		{
			name:       "dotfile-skipped",
			filename:   ".env",
			isDir:      false,
			wantNormal: ".env",
			wantNew:    ".env", // Should be unchanged due to skip
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test file
			fpath := filepath.Join(tmpDir, tt.filename)
			if tt.isDir {
				os.Mkdir(fpath, 0755)
				defer os.Remove(fpath)
			} else {
				os.WriteFile(fpath, []byte(""), 0644)
				defer os.Remove(fpath)
			}

			info, err := os.Stat(fpath)
			if err != nil {
				t.Fatalf("failed to stat test file: %v", err)
			}

			result := processOne(fpath, info)

			if result.OldName != tt.wantNormal {
				t.Errorf("processOne() OldName = %q, want %q", result.OldName, tt.wantNormal)
			}

			if result.NewName != tt.wantNew {
				t.Errorf("processOne() NewName = %q, want %q", result.NewName, tt.wantNew)
			}

			// In dry-run mode, Renamed should be false
			if result.Renamed {
				t.Errorf("processOne() Renamed = %v, want false (dry-run mode)", result.Renamed)
			}

			// Verify the file still exists with the original name
			if _, err := os.Stat(fpath); err != nil {
				t.Errorf("processOne() file was modified during dry-run: %v", err)
			}
		})
	}
}

// TestProcessOneWithDotfiles tests the dotfile handling flag.
func TestProcessOneWithDotfiles(t *testing.T) {
	oldFlagDo := flagDo
	oldFlagCase := flagCase
	oldFlagDateMode := flagDateMode
	oldFlagDotfiles := flagDotfiles

	defer func() {
		flagDo = oldFlagDo
		flagCase = oldFlagCase
		flagDateMode = oldFlagDateMode
		flagDotfiles = oldFlagDotfiles
	}()

	flagDo = false
	flagCase = "lower"
	flagDateMode = ""

	tmpDir := t.TempDir()

	// Create a dotfile
	fpath := filepath.Join(tmpDir, ".MyEnv")
	os.WriteFile(fpath, []byte(""), 0644)
	defer os.Remove(fpath)

	info, _ := os.Stat(fpath)

	// Test with dotfiles disabled (default)
	flagDotfiles = false
	result := processOne(fpath, info)
	if !result.WasSkipped {
		t.Errorf("processOne() with dotfiles disabled: WasSkipped = %v, want true", result.WasSkipped)
	}
	if result.NewName != ".MyEnv" {
		t.Errorf("processOne() with dotfiles disabled: NewName = %q, want %q", result.NewName, ".MyEnv")
	}

	// Test with dotfiles enabled
	flagDotfiles = true
	result = processOne(fpath, info)
	if result.WasSkipped {
		t.Errorf("processOne() with dotfiles enabled: WasSkipped = %v, want false", result.WasSkipped)
	}
	// The dot gets stripped during cleaning, so .MyEnv becomes myenv
	if result.NewName != "myenv" {
		t.Errorf("processOne() with dotfiles enabled: NewName = %q, want %q", result.NewName, "myenv")
	}
}

// Benchmarks

// BenchmarkMakeUnique benchmarks the unique filename generation.
func BenchmarkMakeUnique(b *testing.B) {
	tmpDir := b.TempDir()
	// Create some existing files
	for i := 1; i < 10; i++ {
		fname := filepath.Join(tmpDir, "file.txt")
		if i > 1 {
			fname = filepath.Join(tmpDir, "file_"+string(rune('0'+i))+".txt")
		}
		os.WriteFile(fname, []byte(""), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeUnique(tmpDir, "file.txt")
	}
}
