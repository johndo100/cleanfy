// result.go
// ----------
// Defines the Result struct used to represent the outcome of each rename operation.
// Each Result entry records old/new names, errors, and flags such as Renamed and AutoRenamed.

package main

// Result represents the outcome of processing one file or directory.
type Result struct {
	Path        string `json:"path"`            // Full path to the processed file or directory
	OldName     string `json:"old_name"`        // Original name
	NewName     string `json:"new_name"`        // New (transformed) name
	IsDir       bool   `json:"is_dir"`          // True if the entry is a directory
	Renamed     bool   `json:"renamed"`         // True if a rename actually occurred
	AutoRenamed bool   `json:"auto_renamed"`    // True if a numeric suffix was auto-added to avoid conflicts
	Error       string `json:"error,omitempty"` // Error message if any
}

// HasError reports whether the result contains an error.
func (r Result) HasError() bool {
	return r.Error != ""
}

// IsChanged returns true if the name was altered or auto-renamed.
func (r Result) IsChanged() bool {
	return r.Renamed || r.AutoRenamed
}
