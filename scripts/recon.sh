#!/bin/bash

# Recon Script - Generate fresh file tree of the assigned subfolder
# This script is used by the AI agent to ensure the service map is current

set -e

# Default to current directory if no argument provided
TARGET_DIR="${1:-.}"
OUTPUT_FILE="${2:-recon_output.txt}"

echo "Generating file tree for: $TARGET_DIR"
echo "Output file: $OUTPUT_FILE"

# Generate file tree using find
# Exclude common build/cache directories
find "$TARGET_DIR" -type f \
  -not -path '*/node_modules/*' \
  -not -path '*/dist/*' \
  -not -path '*/__pycache__/*' \
  -not -path '*/.git/*' \
  -not -path '*/venv/*' \
  -not -path '*/.pytest_cache/*' \
  -not -path '*/coverage/*' \
  | sort > "$OUTPUT_FILE"

# Also generate a tree-like view
echo -e "\n=== File Tree (Tree View) ===" >> "$OUTPUT_FILE"
if command -v tree &> /dev/null; then
  tree -I 'node_modules|dist|__pycache__|.git|venv|.pytest_cache|coverage' "$TARGET_DIR" >> "$OUTPUT_FILE" 2>/dev/null || true
else
  # Fallback to find with formatting
  find "$TARGET_DIR" -type f \
    -not -path '*/node_modules/*' \
    -not -path '*/dist/*' \
    -not -path '*/__pycache__/*' \
    -not -path '*/.git/*' \
    -not -path '*/venv/*' \
    -not -path '*/.pytest_cache/*' \
    -not -path '*/coverage/*' \
    | sort | sed 's|[^/]*/|  |g' >> "$OUTPUT_FILE"
fi

echo "File tree generated successfully!"
echo "Contents:"
cat "$OUTPUT_FILE"

# Return success
exit 0
