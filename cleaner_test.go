// cleaner_test.go
// ----------------
// Unit tests for filename cleaning and normalization functions.
// Tests cover ASCII conversion, POSIX filtering, case transformations,
// and the full CleanName pipeline.

package main

import (
	"testing"
)

// TestCleanASCII tests Unicode to ASCII conversion with special character mappings.
func TestCleanASCII(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// ASCII passthrough
		{"ascii-letters", "hello", "hello"},
		{"ascii-numbers", "file123", "file123"},
		{"ascii-mixed", "Hello123", "Hello123"},

		// Accent removal (NFKD decomposition)
		{"acute-accent", "café", "cafe"},
		{"grave-accent", "naïve", "naive"},
		{"circumflex", "château", "chateau"},
		{"tilde", "mañana", "manana"},
		{"diaeresis", "Müller", "Muller"},

		// Special character mappings
		{"eszett", "straße", "strasse"},
		{"ae-ligature-lower", "Ægis", "aegis"},
		{"ae-ligature-upper", "Ængel", "aengel"},
		{"oe-ligature-lower", "œuvre", "oeuvre"},
		{"oe-ligature-upper", "Œuvre", "oeuvre"},
		{"d-stroke-lower", "đặc", "dac"},
		{"d-stroke-upper", "Đặc", "dac"},
		{"eth", "ð", "d"},
		{"eth-upper", "Ð", "d"},
		{"l-stroke-lower", "Łódź", "lodz"},
		{"l-stroke-upper", "ŁÓDŹ", "lODZ"},
		{"dong-currency", "₫", "d"},

		// Quotation marks and dashes
		{"en-dash", "file–name", "file-name"},
		{"em-dash", "file—name", "file-name"},
		{"curly-quote-left", "hello", "hello"},
		{"curly-quote-right", "hello", "hello"},
		{"low-quote", "hello", "hello"},
		{"guillemets", "hello", "hello"},
		{"prime", "2", "2"},
		{"double-prime", "2", "2"},

		// Bullet and middle dot
		{"middle-dot", "file·name", "file-name"},
		{"bullet", "file•name", "file-name"},
		{"bullet-operator", "file∙name", "file-name"},

		// Empty and whitespace
		{"empty", "", ""},
		{"space", "hello world", "hello world"},
		{"tab", "hello\tworld", "helloworld"},

		// Unknown/replacement
		{"emoji", "file_name", "file_name"},
		{"special-char", "file_name", "file_name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanASCII(tt.input)
			if got != tt.expected {
				t.Errorf("cleanASCII(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestPosixify tests POSIX filename filtering.
func TestPosixify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic POSIX-safe names
		{"lowercase", "hello", "hello"},
		{"uppercase", "HELLO", "HELLO"},
		{"numbers", "file123", "file123"},
		{"underscore", "file_name", "file_name"},
		{"dash", "file-name", "file-name"},
		{"dot", "file.txt", "file.txt"},

		// Spaces to underscores
		{"single-space", "hello world", "hello_world"},
		{"multiple-spaces", "hello  world", "hello_world"},

		// Disallowed characters converted to underscore
		{"special-chars", "file@name#test", "file_name_test"},
		{"brackets", "file[name]", "file_name"},
		{"parens", "file(name)", "file_name"},
		{"comma", "file,name", "file_name"},

		// Collapse multiple underscores
		{"double-underscore", "file__name", "file_name"},
		{"triple-underscore", "file___name", "file_name"},

		// Collapse multiple dashes
		{"double-dash", "file--name", "file-name"},
		{"triple-dash", "file---name", "file-name"},

		// Trim leading/trailing dots, underscores, dashes
		{"leading-dot", ".filename", "filename"},
		{"trailing-dot", "filename.", "filename"},
		{"leading-underscore", "_filename", "filename"},
		{"trailing-underscore", "filename_", "filename"},
		{"leading-dash", "-filename", "filename"},
		{"trailing-dash", "filename-", "filename"},

		// Mixed trimming
		{"leading-dot-underscore", "._filename", "filename"},
		{"trailing-mixed", "filename._-", "filename"},

		// Empty result fallback
		{"only-dots", "...", "_"},
		{"only-underscores", "___", "_"},
		{"only-dashes", "---", "_"},
		{"only-special", "###", "_"},

		// Real-world examples
		{"windows-reserved-chars", "file:name", "file_name"},
		{"unicode-replaced", "file_name", "file_name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := posixify(tt.input)
			if got != tt.expected {
				t.Errorf("posixify(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestToTitle tests title case transformation.
func TestToTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase", "hello world", "Hello World"},
		{"uppercase", "HELLO WORLD", "Hello World"},
		{"mixed", "hELLo WoRLD", "Hello World"},
		{"numbers", "file 123 name", "File 123 Name"},
		{"underscores", "hello_world", "Hello_World"},
		{"dashes", "hello-world", "Hello-World"},
		{"dots", "hello.world", "Hello.World"},
		{"multiple-spaces", "hello  world", "Hello  World"},
		{"single-letter", "a", "A"},
		{"empty", "", ""},
		{"numbers-only", "123", "123"},
		{"mixed-case-reset", "heLLo", "Hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toTitle(tt.input)
			if got != tt.expected {
				t.Errorf("toTitle(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestCleanName tests the full pipeline: ASCII → POSIX → case transform → date prefix.
// Note: We skip date prefixing tests here since they depend on flagDateMode and time.
func TestCleanName(t *testing.T) {
	// Set up test flags
	oldCase := flagCase
	oldDateMode := flagDateMode
	defer func() {
		flagCase = oldCase
		flagDateMode = oldDateMode
	}()

	tests := []struct {
		name      string
		fullPath  string
		filename  string
		isDir     bool
		caseMode  string
		dateMode  string
		expected  string
		expectErr bool
	}{
		// Simple files with lowercase transform
		{
			name:      "simple-file-lower",
			fullPath:  "/tmp/test.txt",
			filename:  "MyFile.txt",
			isDir:     false,
			caseMode:  "lower",
			dateMode:  "",
			expected:  "myfile.txt",
			expectErr: false,
		},
		{
			name:      "unicode-file-lower",
			fullPath:  "/tmp/test.txt",
			filename:  "Café.txt",
			isDir:     false,
			caseMode:  "lower",
			dateMode:  "",
			expected:  "cafe.txt",
			expectErr: false,
		},

		// Uppercase transform
		{
			name:      "file-upper",
			fullPath:  "/tmp/test.txt",
			filename:  "myfile.txt",
			isDir:     false,
			caseMode:  "upper",
			dateMode:  "",
			expected:  "MYFILE.TXT",
			expectErr: false,
		},

		// Title case
		{
			name:      "file-title",
			fullPath:  "/tmp/test.txt",
			filename:  "hello_world.txt",
			isDir:     false,
			caseMode:  "title",
			dateMode:  "",
			expected:  "Hello_World.txt",
			expectErr: false,
		},

		// Directory (no extension)
		{
			name:      "directory",
			fullPath:  "/tmp/MyDir",
			filename:  "MyDir",
			isDir:     true,
			caseMode:  "lower",
			dateMode:  "",
			expected:  "mydir",
			expectErr: false,
		},

		// File with no extension
		{
			name:      "no-extension",
			fullPath:  "/tmp/README",
			filename:  "README",
			isDir:     false,
			caseMode:  "lower",
			dateMode:  "",
			expected:  "readme",
			expectErr: false,
		},

		// Multiple dots (only last is extension)
		{
			name:      "multiple-dots",
			fullPath:  "/tmp/archive.tar.gz",
			filename:  "archive.tar.gz",
			isDir:     false,
			caseMode:  "lower",
			dateMode:  "",
			expected:  "archive.tar.gz",
			expectErr: false,
		},

		// Special characters normalization
		{
			name:      "special-chars",
			fullPath:  "/tmp/test",
			filename:  "Łódź – testowy plik.txt",
			isDir:     false,
			caseMode:  "lower",
			dateMode:  "",
			expected:  "lodz_-_testowy_plik.txt",
			expectErr: false,
		},

		// Reserved names (Windows)
		{
			name:      "reserved-name-COM",
			fullPath:  "/tmp/COM",
			filename:  "COM",
			isDir:     false,
			caseMode:  "lower",
			dateMode:  "",
			expected:  "com",
			expectErr: false,
		},

		// Spaces and punctuation
		{
			name:      "spaces-and-punct",
			fullPath:  "/tmp/test",
			filename:  "My File (v2) [Final].txt",
			isDir:     false,
			caseMode:  "lower",
			dateMode:  "",
			expected:  "my_file_v2_final.txt",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flagCase = tt.caseMode
			flagDateMode = tt.dateMode

			got, err := CleanName(tt.fullPath, tt.filename, tt.isDir)
			if (err != nil) != tt.expectErr {
				t.Errorf("CleanName() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && got != tt.expected {
				t.Errorf("CleanName(%q, %q, %v) = %q, want %q", tt.fullPath, tt.filename, tt.isDir, got, tt.expected)
			}
		})
	}
}

// Benchmarks (optional, useful for performance tracking)

// BenchmarkCleanASCII benchmarks the ASCII conversion function.
func BenchmarkCleanASCII(b *testing.B) {
	input := "Café Münster – Łódź с кириллицей"
	for i := 0; i < b.N; i++ {
		cleanASCII(input)
	}
}

// BenchmarkPosixify benchmarks the POSIX filtering function.
func BenchmarkPosixify(b *testing.B) {
	input := "file___name--with___dashes...and___dots"
	for i := 0; i < b.N; i++ {
		posixify(input)
	}
}

// BenchmarkCleanName benchmarks the full pipeline.
func BenchmarkCleanName(b *testing.B) {
	flagCase = "lower"
	flagDateMode = ""
	input := "My File Café – Münster.txt"
	for i := 0; i < b.N; i++ {
		CleanName("/tmp/test.txt", input, false)
	}
}
