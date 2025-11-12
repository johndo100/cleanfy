# cleanfy

A smart, fast CLI tool to batch normalize and clean filenames with Unicode support, case transformation, and intelligent conflict resolution.

## Features

- 🌍 **Unicode Support** — Converts accented characters and special letters to ASCII equivalents
- 🔧 **Case Transformation** — Lowercase, uppercase, or title case (optional, no forced defaults)
- 📅 **Date Prefixing** — Add file modification date or current date as prefix (optional)
- ⚡ **Conflict Resolution** — Automatically handles duplicate filenames with numeric suffixes (`_2`, `_3`, etc.)
- 🏃 **Recursive Processing** — Optionally process subdirectories with `-r`
- 👁️ **Dry-Run Mode** — Preview all changes before applying (default behavior)
- 📋 **JSON Output** — Machine-readable results for automation and scripting
- 🚀 **Fast & Safe** — Handles thousands of files efficiently with proper error handling

## Installation

### From Source
```bash
git clone https://github.com/johndo100/cleanfy.git
cd cleanfy
go build -o cleanfy
```

### Binary Release
Download prebuilt binaries from [Releases](https://github.com/johndo100/cleanfy/releases)

## Quick Start

### Preview Changes (Dry-Run)
```bash
cleanfy ./photos
```

### Execute Changes
```bash
cleanfy -x ./photos
```

### With Options
```bash
cleanfy -x -r --case=lower ./all_files
```

## Flags

| Flag | Long Form | Description |
|------|-----------|-------------|
| `-x` | `--execute` | **Execute renames** (default: preview only) |
| `-r` | `--recursive` | Recurse into subdirectories |
| `-q` | `--quiet` | Suppress output (errors only) |
| `-j` | `--json` | Output as JSON |
| `-v` | `--version` | Show version and exit |
| `-a` | `--dotfiles` | Include hidden files (starting with `.`) |
| `--case=` | `--case=` | Case transform: `lower`, `upper`, `title` (optional) |
| `--date=` | `--date=` | Date prefix: `mtime` (modified) or `now` (current) (optional) |
| `--date-format=` | `--date-format=` | Go time layout (default: `2006-01-02`) |

**Note:** All value flags must use the `=` form (e.g., `--case=lower`, `--date=mtime`)

## Examples

```bash
# Preview changes
cleanfy ./photos

# Apply lowercase transformation to all files
cleanfy -x --case=lower ./photos

# Recursive processing with uppercase
cleanfy -x -r --case=upper ./documents

# Add current date prefix to files
cleanfy -x --date=now ./backup

# Add file modification date as prefix with custom format
cleanfy -x --date=mtime --date-format="20060102_150405" ./archive

# Recursive with date and custom case
cleanfy -x -r --case=title --date=mtime ./library

# Process including hidden files
cleanfy -x -a ./config

# JSON output for scripting
cleanfy -j ./files | jq '.[] | select(.renamed == true)'

# Quiet mode with errors only
cleanfy -x -q -r ./large_directory

# Custom date format examples
cleanfy --date=now --date-format="2006-01-02" ./files         # 2025-11-12
cleanfy --date=now --date-format="Jan 2 2006" ./files         # Nov 12 2025
cleanfy --date=now --date-format="20060102_150405" ./files    # 20251112_165030
```

## What Gets Normalized

### Character Transformations
- **Accents** — `café` → `cafe`, `Zürich` → `Zurich`
- **Special Ligatures** — `straße` → `strasse`, `Æther` → `AEther`
- **Strokes** — `Łódź` → `Lodz`, `Đàm` → `dam`
- **Quotes & Dashes** — Smart handling of curly quotes, en-dashes, em-dashes

### Filename Cleanup
- **Spaces & Punctuation** — Converted to underscores: `My File!` → `my_file`
- **Multiple Separators** — Collapsed: `file___name` → `file_name`
- **Leading/Trailing** — Trimmed: `.file_` → `file`
- **Reserved Names** — Windows reserved names prefixed with `_`: `COM` → `_com`
- **Length** — Safely truncated while preserving UTF-8 validity

### Optional Transforms
- **Case** — Lower, upper, or title case (opt-in via `--case=`)
- **Date Prefix** — Add mtime or current date (opt-in via `--date=`)

### Conflict Resolution
- **Duplicates** — Auto-resolved with numeric suffixes: `file.txt` → `file_2.txt`
- **Always On** — Prevents overwrites automatically

## Output Format

### Preview Mode (Text Output)
```
RENAME  MyFile.txt -> myfile.txt
OK      already_clean.txt
RENAME* Duplicate.pdf -> duplicate_2.pdf   (auto-resolved)
ERR     Protected.txt : permission denied
```

### Execution Mode (Text Output)
```
RENAME  MyFile.txt -> myfile.txt
RENAME* Duplicate.pdf -> duplicate_2.pdf   (auto-resolved)
```

### JSON Output
```json
[
  {
    "path": "/path/to/myfile.txt",
    "old_name": "MyFile.txt",
    "new_name": "myfile.txt",
    "is_dir": false,
    "renamed": true,
    "auto_renamed": false,
    "skipped": false,
    "error": ""
  }
]
```

## Testing

Run the full test suite:

```bash
# Verbose output
go test ./... -v

# With coverage report
go test ./... -cover

# Run benchmarks
go test -bench=. ./...
```

**Current Coverage:** 34.5% with 100+ test cases across core functions

## Performance

- **Single directory (100 files):** ~50ms
- **Recursive (1000 files):** ~200ms  
- **Large batch (10,000 files):** ~2s

## Safety Features

✅ **Dry-run by default** — Use `-x` to apply changes  
✅ **Dotfiles skipped by default** — Use `-a` to process hidden files  
✅ **No forced transformations** — Case/date are optional  
✅ **Automatic conflict resolution** — Prevents overwrites  
✅ **Error reporting** — Clear feedback on failures  

## Platform Support

- ✅ Linux
- ✅ macOS (Intel & Apple Silicon)
- ✅ Windows

## Version

Current version is automatically detected from git tags and build metadata.

```bash
cleanfy -v
# cleanfy v0.9.0-pre (linux/amd64)
```

## CI/CD

This project uses GitHub Actions for:
- ✅ Automated testing on push/PR
- ✅ Cross-platform builds on tag
- ✅ Automatic release creation

See [CI_RELEASE_GUIDE.md](CI_RELEASE_GUIDE.md) for details.

## License

See LICENSE file for details

## Support

Found a bug or have a feature request?  
Please open an issue on [GitHub](https://github.com/johndo100/cleanfy/issues)

---

**Made with ❤️ for batch file normalization**
