# cleanfy

A command-line tool to normalize and clean filenames in bulk.

## Features

- 🌍 **Unicode Support** — Converts accented characters and special letters to ASCII equivalents
- 🔧 **Case Transformation** — Lowercase, uppercase, or title case
- 📅 **Date Prefixing** — Add creation or modification date to filenames
- ⚡ **Conflict Resolution** — Automatically handles duplicate filenames with numeric suffixes
- 🏃 **Recursive Processing** — Optionally process subdirectories
- 👁️ **Dry-Run Mode** — Preview changes before applying them
- 📋 **JSON Output** — Machine-readable results for automation

## Installation

```bash
go build -o cleanfy
```

## Usage

### Basic Preview (Dry-Run)
```bash
cleanfy /path/to/directory
```

### Apply Changes
```bash
cleanfy -d /path/to/directory
```

### Recursive with Lowercase
```bash
cleanfy -r -c lower /path/to/directory
```

### Add Date Prefix (Modified Date)
```bash
cleanfy -d -t mtime -f "2006-01-02" /path/to/directory
```

### JSON Output
```bash
cleanfy -j /path/to/directory
```

## Flags

| Flag | Long Form | Description |
|------|-----------|-------------|
| `-d` | `--do` | Apply changes (default: dry-run preview only) |
| `-r` | `--recursive` | Recurse into subdirectories |
| `-q` | `--quiet` | Suppress output (errors only) |
| `-j` | `--json` | Output as JSON |
| `-v` | `--version` | Show version info |
| `-a` | `--dotfiles` | Include hidden files (starting with `.`) |
| `-c` | `--case` | Case transform: `lower`, `upper`, `title` (default: `lower`) |
| `-t` | `--date` | Date prefix mode: `mtime` (file modified), `now` (current time) |
| `-f` | `--date-format` | Date format using Go time layout (default: `2006-01-02`) |

## Examples

```bash
# Preview all changes
cleanfy ./photos

# Apply lowercase transformation
cleanfy -d -c lower ./photos

# Add current date to filenames
cleanfy -d -t now ./photos

# Recursive with date prefix
cleanfy -d -r -t mtime ./documents

# JSON output for scripting
cleanfy -j ./files | jq '.[] | select(.renamed == true)'

# Include hidden files
cleanfy -a ./config

# Custom date format (YYYY-MM-DD HH:MM)
cleanfy -t mtime -f "2006-01-02_15-04" ./backup
```

## What Gets Normalized

- **Accents** — `café` → `cafe`
- **Special Characters** — `Łódź` → `Lodz`, `straße` → `strasse`
- **Spaces & Punctuation** — `My File!` → `my_file`
- **Reserved Names** — Windows reserved names are prefixed with `_`
- **Case** — Based on `-c` flag (default: lowercase)
- **Length** — Long filenames are truncated safely
- **Conflicts** — Duplicates get numeric suffixes (`_2`, `_3`, etc.)

## Output Format

### Dry-Run (Text)
```
RENAME  OldFile.txt -> oldfile.txt
OK      already_normalized.txt
RENAME* Duplicate.txt -> duplicate_2.txt   (auto-resolved)
ERR     Problem.txt : permission denied
```

### JSON
```json
[
  {
    "path": "/path/to/file.txt",
    "old_name": "OldFile.txt",
    "new_name": "oldfile.txt",
    "is_dir": false,
    "renamed": true,
    "auto_renamed": false,
    "skipped": false,
    "error": ""
  }
]
```

## Testing

Run the test suite:

```bash
go test ./... -v           # Verbose output
go test ./... -cover       # With coverage report
go test -bench=. ./...     # Run benchmarks
```

Test coverage: **43.1%** with 103 test cases

## Performance

- Single directory: < 100ms
- Recursive (1000 files): ~500ms
- Benchmarks included for critical paths

## Notes

- **Dry-run by default** — Use `-d` to actually rename files
- **Dotfiles skipped by default** — Use `-a` to process hidden files
- **Conflict resolution enabled** — Duplicates automatically get suffixes
- **Case-sensitive on Unix** — File systems may differ in case handling
- **Recursive is opt-in** — Use `-r` to process subdirectories

## License

See LICENSE file for details

## Support

For issues or feature requests, please open an issue on GitHub.
