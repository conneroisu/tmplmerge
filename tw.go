package twerge

import (
	"bytes"
	"fmt"
	"os"
	"slices"
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
func GenerateTailwind(
	cssPath string,
) error {
	// Read base CSS content if the file exists
	var baseContent []byte
	var err error

	baseContent, err = os.ReadFile(cssPath)
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

	var builder strings.Builder
	var gendClasses []string
	for generated, merged := range sortMap(GenClassMergeStr) {
		gendClasses = append(gendClasses, generated)
		// Create a CSS rule using the generated class name and the merged Tailwind classes
		builder.WriteString(".")
		builder.WriteString(generated)
		builder.WriteString(" { \n\t@apply ")
		builder.WriteString(merged)
		builder.WriteString("; \n}\n")
	}
	for givenClasses, gendClass := range ClassMapStr {
		if slices.Contains(gendClasses, gendClass) {
			continue
		}
		builder.WriteString(".")
		builder.WriteString(gendClass)
		builder.WriteString(" { \n\t@apply ")
		builder.WriteString(Merge(givenClasses))
		builder.WriteString("; \n}\n")
	}
	cssContent := builder.String()

	// Add to file content
	newContent, err := replaceBetweenMarkers(baseContent, []byte(cssContent))
	if err != nil {
		return fmt.Errorf("error adding twerge content: %w", err)
	}

	// Write to output path
	err = os.WriteFile(cssPath, newContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	return nil
}

func sortMap(m map[string]string) map[string]string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sorted := make(map[string]string, len(m))
	for _, k := range keys {
		sorted[k] = m[k]
	}
	return sorted
}

// GenerateTempl creates a .templ file that can be used to generate a CSS file
// with the provided class map.
func GenerateTempl(
	templPath string,
) error {
	var pkgName string
	pkgEnd := strings.LastIndex(templPath, "/")
	if pkgEnd == -1 {
		pkgName = "main"
	} else {
		pkgStart := strings.LastIndex(templPath[:pkgEnd], "/")
		if pkgStart == -1 {
			pkgName = "main"
		} else {
			pkgName = templPath[pkgStart+1 : pkgEnd]
		}
	}

	var buf bytes.Buffer
	buf.WriteString("// Code generated by twerge. DO NOT EDIT.\n\n")
	buf.WriteString("package ")
	buf.WriteString(pkgName)
	buf.WriteString("\n\n")
	buf.WriteString("templ empty() {\n")
	buf.WriteString("<div class=\"")
	buf.WriteString("mb-4")
	buf.WriteString("\"></div>\n")
	for k := range sortMap(GenClassMergeStr) {
		// Create a CSS rule using the generated class name and the merged Tailwind classes
		buf.WriteString("<div class=\"")
		buf.WriteString(k)
		buf.WriteString("\"></div>\n")
	}
	for _, v := range sortMap(ClassMapStr) {
		buf.WriteString("<div class=\"")
		buf.WriteString(v)
		buf.WriteString("\"></div>\n")
	}
	buf.WriteString("}")

	err := os.WriteFile(templPath, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing .templ file: %w", err)
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
