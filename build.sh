#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:?Usage : ./build.sh <version>}"
DIST="dist"
ENTRYPOINT="./cmd/TP5"

rm -rf "$DIST"
mkdir -p "$DIST"

# GOOS/GOARCH/extension
PLATFORMS=(
  "linux/amd64/"
  "windows/amd64/.exe"
  "linux/arm64/"
)

for entry in "${PLATFORMS[@]}"; do
  IFS='/' read -r GOOS_VAL GOARCH_VAL EXT <<< "$entry"
  OUT="$DIST/gopack-${GOOS_VAL}-${GOARCH_VAL}${EXT}"
  echo "Build $OUT..."
  GOOS="$GOOS_VAL" GOARCH="$GOARCH_VAL" go build \
    -ldflags "-s -w -X main.version=${VERSION}" \
    -o "$OUT" "$ENTRYPOINT"
done

echo "Génération des checksums..."
(cd "$DIST" && sha256sum gopack-* > SHA256SUMS)

echo "Terminé. Contenu de $DIST :"
ls -lh "$DIST"
