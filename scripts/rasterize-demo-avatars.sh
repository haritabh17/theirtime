#!/usr/bin/env bash
# Regenerate demo avatar PNGs from SVG sources (macOS sips).
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
demo="$root/internal/menubar/demo"
size=128

for name in bob ann; do
  sips -s format png "$demo/$name.svg" --out "$demo/$name.png" >/dev/null
  sips -z "$size" "$size" "$demo/$name.png" >/dev/null
done

echo "Wrote $demo/bob.png and $demo/ann.png (${size}x${size})"
