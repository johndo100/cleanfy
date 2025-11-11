#!/usr/bin/env bash
# normfn-go sandbox builder
# Creates test files/folders listed in scripts/sandbox_files.txt

set -e          # exit on any error
set -f          # disable globbing and command substitution (important for $(...) etc.)
cd "$(dirname "$0")/.."

SANDBOX="sandbox"
LIST_FILE="scripts/sandbox_files.txt"
CLEAN=true
FORCE=false
FLAT=false

# ─────────────────────────────────────────────
# Parse command-line flags
# ─────────────────────────────────────────────
for arg in "$@"; do
  case "$arg" in
    --clean) CLEAN=true ;;
    --force) FORCE=true ;;
    --flat)  FLAT=true ;;
    --list=*)
      LIST_FILE="${arg#--list=}"
      ;;
  esac
done

# ─────────────────────────────────────────────
# Check list file
# ─────────────────────────────────────────────
if [[ ! -f "$LIST_FILE" ]]; then
  echo "❌ File list not found: $LIST_FILE"
  exit 1
fi

# ─────────────────────────────────────────────
# Clean sandbox if requested
# ─────────────────────────────────────────────
if $CLEAN; then
  echo "🧹 Cleaning sandbox..."
  rm -rf "$SANDBOX"
fi
mkdir -p "$SANDBOX"

echo "📁 Using file list: $LIST_FILE"

# ─────────────────────────────────────────────
# Main creation loop
# Reads every line literally, even without trailing newline
# ─────────────────────────────────────────────
while IFS= read -r raw || [[ -n "$raw" ]]; do
  line="${raw%$'\r'}"  # strip Windows CR if present
  [[ -z "$line" ]] && continue  # skip empty lines

  # Handle quoted filenames
  if [[ "$line" == \"*\" && "$line" == *\" ]]; then
    filename="${line:1:${#line}-2}"
  elif [[ "$line" == \'*\' && "$line" == *\' ]]; then
    filename="${line:1:${#line}-2}"
  else
    # Skip comment lines
    [[ "${line:0:1}" == "#" ]] && continue
    filename="$line"
  fi

  # If flat mode, ignore folder structure
  if $FLAT; then
    filename="$(basename "$filename")"
  fi

  target="$SANDBOX/$filename"

  # Skip existing files unless --force
  if [[ -e "$target" && $FORCE == false ]]; then
    echo "⚠️  Skipped (exists): $filename"
    continue
  fi

  # Ensure directory structure if not flat
  mkdir -p "$(dirname "$target")"

  # Create empty file
  : > "$target"
  echo "✅ Created: $filename"
done < "$LIST_FILE"

echo
echo "🎯 Sandbox ready: $SANDBOX/"
ls -1A "$SANDBOX"
