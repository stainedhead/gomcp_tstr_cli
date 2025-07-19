#!/bin/bash

# Script to combine all slide files into a single presentation
# This makes it easier to use with various Markdown presentation tools

OUTPUT_FILE="mcp_tstr-presentation.md"

echo "Generating combined presentation file: $OUTPUT_FILE"

# Start with an empty file
> $OUTPUT_FILE

# Add YAML front matter for Marp
cat << EOF >> $OUTPUT_FILE
---
marp: true
theme: default
paginate: true
header: "mcp_tstr - Model Context Protocol Testing Tool"
footer: "© $(date +%Y) - mcp_tstr Project"
---

EOF

# Combine all numbered markdown files in order
for file in $(ls -1 [0-9]*.md | sort -n); do
  echo "Adding $file..."
  cat "$file" >> $OUTPUT_FILE
done

echo "Presentation generated successfully!"
echo "You can now use this file with Marp or other Markdown presentation tools."
