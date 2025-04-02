#!/bin/bash

# Find all .templ files
find . -name "*.templ" -type f | while read file; do
  # Read the content
  content=$(cat "$file")
  
  # Process with sad and capture the output
  new_content=$(echo "$content" | sad 'twerge.Generate\(' 'twerge.Generate(' | sad 'twerge.It\(' 'twerge.It(')
  
  # Write back to the same file
  echo "$new_content" > "$file"
  
  echo "Processed: $file"
done
