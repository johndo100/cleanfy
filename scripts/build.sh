#!/usr/bin/env bash
# =========================================================
# Cleanfy build script
# Cross-compiles Cleanfy for multiple platforms.
# Output binaries go into ./bin/<os>-<arch>/
# =========================================================

set -e
cd "$(dirname "$0")/.."

APP_NAME="cleanfy"
BIN_DIR="bin"
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")

# Supported platforms for build
PLATFORMS=(
  "linux-amd64"
  "windows-amd64"
  "darwin-arm64"
)

echo "🏗️  Building $APP_NAME (version: $VERSION)"
echo

for platform in "${PLATFORMS[@]}"; do
  GOOS="${platform%-*}"
  GOARCH="${platform#*-}"
  OUTDIR="$BIN_DIR/$platform"
  mkdir -p "$OUTDIR"

  BINNAME="$APP_NAME"
  [[ "$GOOS" == "windows" ]] && BINNAME="$APP_NAME.exe"

  echo "  → Building for $platform"
  GOOS=$GOOS GOARCH=$GOARCH go build \
    -ldflags="-s -w -X main.version=$VERSION" \
    -o "$OUTDIR/$BINNAME" .

done

echo
echo "✅ Builds complete!"
if command -v tree >/dev/null 2>&1; then
  tree "$BIN_DIR"
else
  ls -R "$BIN_DIR"
fi
