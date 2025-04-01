package twerge

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	// Section markers used to identify where generated CSS should be placed

	// twergeBeginMarker is the beginning of the section where the generated CSS will be placed
	twergeBeginMarker = "/* twerge:begin */"
	// twergeEndMarker is the end of the section where the generated CSS will be placed
	twergeEndMarker = "/* twerge:end */"
)

// GenerateTailwind creates an input CSS file for the Tailwind CLI
// that includes all the @apply directives from the provided class map.
//
// This is useful for building a production CSS file with Tailwind's CLI.
//
// The marker is used to identify the start and end of the @apply directives generated
// by Twerge.
func GenerateTailwind(inputPath, outputPath string, classMap map[string]string) error {
	// Read base CSS content if the file exists
	var baseContent []byte
	var err error

	baseContent, err = os.ReadFile(inputPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading input file: %w", err)
	}

	// If file doesn't exist, create minimal Tailwind directives
	if os.IsNotExist(err) {
		baseContent = []byte(`@tailwind base;
@tailwind components;
@tailwind utilities;

/* twerge:start */
/* twerge:end */
`)
	}

	// Generate Twerge CSS content
	// Get all keys and sort them for consistent output
	keys := make([]string, 0, len(classMap))
	for k := range classMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, k := range keys {
		original := k
		generated := classMap[k]
		merged := Merge(original)

		// Create a CSS rule using the generated class name and the merged Tailwind classes
		builder.WriteString(".")
		builder.WriteString(generated)
		builder.WriteString(" { \n\t@apply ")
		builder.WriteString(merged)
		builder.WriteString("; \n}\n")
	}
	cssContent := builder.String()

	// Add to file content
	newContent, err := replaceBetweenMarkers(baseContent, []byte(cssContent))
	if err != nil {
		return fmt.Errorf("error adding twerge content: %w", err)
	}

	// Write to output path
	err = os.WriteFile(outputPath, newContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	return nil
}

// replaceBetweenMarkers replaces content between twerge markers
func replaceBetweenMarkers(content, replacement []byte) ([]byte, error) {
	// Find begin marker
	beginMarkerBytes := []byte(twergeBeginMarker)
	beginIdx := bytes.Index(content, beginMarkerBytes)
	if beginIdx == -1 {
		// Markers don't exist, append content with markers
		suffix := append([]byte("\n\n"), beginMarkerBytes...)
		suffix = append(suffix, '\n')
		suffix = append(suffix, replacement...)
		suffix = append(suffix, '\n')
		suffix = append(suffix, []byte(twergeEndMarker)...)
		return append(content, suffix...), nil
	}

	// Find the end of the line containing the begin marker
	beginLineEnd := beginIdx + len(beginMarkerBytes)
	for beginLineEnd < len(content) && content[beginLineEnd] != '\n' && content[beginLineEnd] != '\r' {
		beginLineEnd++
	}
	if beginLineEnd < len(content) {
		beginLineEnd++ // Include the newline character
	}

	// Find end marker
	endMarkerBytes := []byte(twergeEndMarker)
	endIdx := bytes.Index(content[beginLineEnd:], endMarkerBytes)
	if endIdx == -1 {
		return nil, fmt.Errorf("found begin marker but no end marker")
	}

	// Adjust end marker index to be relative to the whole content
	endIdx += beginLineEnd

	// Create new content with replacement
	result := make([]byte, 0, len(content)-(endIdx-beginLineEnd)+len(replacement)+1)
	result = append(result, content[:beginLineEnd]...)
	result = append(result, replacement...)
	result = append(result, '\n')
	result = append(result, content[endIdx:]...)

	return result, nil
}
