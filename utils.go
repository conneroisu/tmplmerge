package twerge

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	maps0 "maps"
	"os"
	"sort"
	"strings"
	"time"
)

// Section markers used to identify where generated CSS should be placed
const (
	TwergeBeginMarker = "<!-- twerge:begin -->"
	TwergeEndMarker   = "<!-- twerge:end -->"
)

// ExportCSS exports the CSS generated from the RuntimeClassMap to a file,
// replacing content between the twerge markers if they exist.
// If the file doesn't exist, it will be created.
// If the markers don't exist, the content will be appended to the file.
func ExportCSS(filePath string) error {
	cssContent := GetRuntimeClassHTML()
	return WriteCSSToFile(filePath, cssContent)
}

// ExportCSSWithMap exports the CSS generated from a specific class map to a file,
// replacing content between the twerge markers if they exist.
func ExportCSSWithMap(filePath string, classMap map[string]string) error {
	// Store original map
	originalMap := RuntimeClassMap

	// Replace with the provided map temporarily
	runtimeMapMutex.Lock()
	RuntimeClassMap = classMap
	runtimeMapMutex.Unlock()

	// Generate CSS
	cssContent := GetRuntimeClassHTML()

	// Restore original map
	runtimeMapMutex.Lock()
	RuntimeClassMap = originalMap
	runtimeMapMutex.Unlock()

	return WriteCSSToFile(filePath, cssContent)
}

// WriteCSSToFile writes CSS content to a file, handling the twerge markers.
func WriteCSSToFile(filePath string, cssContent string) error {
	// Check if file exists
	fileContent, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading file: %w", err)
	}

	var newContent []byte

	if os.IsNotExist(err) {
		// File doesn't exist, create new content with markers
		newContent = fmt.Appendf(nil, "%s\n%s\n%s\n", TwergeBeginMarker, cssContent, TwergeEndMarker)
	} else {
		// File exists, replace content between markers using optimized approach
		var result []byte
		var err error

		// Use our optimized replaceBetweenMarkers function
		result, err = replaceBetweenMarkers(fileContent, []byte(cssContent))
		if err != nil {
			return fmt.Errorf("error replacing content: %w", err)
		}

		newContent = result
	}

	// Write the new content back to the file
	err = os.WriteFile(filePath, newContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

// replaceBetweenMarkers replaces content between twerge markers
func replaceBetweenMarkers(content, replacement []byte) ([]byte, error) {
	// Find begin marker
	beginMarkerBytes := []byte(TwergeBeginMarker)
	beginIdx := bytes.Index(content, beginMarkerBytes)
	if beginIdx == -1 {
		// Markers don't exist, append content with markers
		suffix := append([]byte("\n\n"), beginMarkerBytes...)
		suffix = append(suffix, '\n')
		suffix = append(suffix, replacement...)
		suffix = append(suffix, '\n')
		suffix = append(suffix, []byte(TwergeEndMarker)...)
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
	endMarkerBytes := []byte(TwergeEndMarker)
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

// GenerateInputCSSForTailwind creates an input CSS file for the Tailwind CLI
// that includes all the @apply directives from Twerge's RuntimeClassMap.
// This is useful for building a production CSS file with Tailwind's CLI.
func GenerateInputCSSForTailwind(inputPath, outputPath string) error {
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

/* Custom CSS goes here */
`)
	}

	// Generate Twerge CSS content
	cssContent := GetRuntimeClassHTML()

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

// AppendClassesToFile appends class definitions to an existing CSS file
// or creates a new file if it doesn't exist.
func AppendClassesToFile(filePath string, classes map[string]string, prefix string) error {
	// Create a buffer to store the CSS content
	var buffer bytes.Buffer

	// Add a header comment
	buffer.WriteString("/* Generated by Twerge */\n\n")

	// Process each class
	for original, generated := range classes {
		merged := Merge(original)

		// Add the CSS rule
		buffer.WriteString(fmt.Sprintf(".%s%s {\n", prefix, generated))
		buffer.WriteString(fmt.Sprintf("  @apply %s;\n", merged))
		buffer.WriteString("}\n\n")
	}

	// Get the content as a string
	content := buffer.String()

	// Write to file, replacing between markers if they exist
	return WriteCSSToFile(filePath, content)
}

// MergeCSSMaps merges multiple class maps into a single map
// In case of conflicts, later maps take precedence
func MergeCSSMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)

	for _, m := range maps {
		maps0.Copy(result, m)
	}

	return result
}

// ProcessCSSTemplate reads a CSS template file containing @apply directives
// and fills in the content between the twerge markers with processed class names.
func ProcessCSSTemplate(templatePath, outputPath string) error {
	// Read the template
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("error reading template: %w", err)
	}

	// Get the CSS from RuntimeClassMap
	cssContent := GetRuntimeClassHTML()

	// Replace the markers
	newContent, err := replaceBetweenMarkers(content, []byte(cssContent))
	if err != nil {
		return fmt.Errorf("error processing template: %w", err)
	}

	// Write the result
	err = os.WriteFile(outputPath, newContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}

	return nil
}

// CSSExportOptions represents configuration options for CSS export
type CSSExportOptions struct {
	// Prefix is applied to all class names
	Prefix string `json:"prefix"`
	// Minify indicates whether the output should be minified
	Minify bool `json:"minify"`
	// Format specifies the output format (css, scss, less)
	Format string `json:"format"`
	// Comments indicates whether to include comments
	Comments bool `json:"comments"`
	// Markers are the custom begin/end markers (defaults to twerge markers if empty)
	BeginMarker string `json:"beginMarker"`
	EndMarker   string `json:"endMarker"`
}

// DefaultCSSExportOptions returns the default options for CSS export
func DefaultCSSExportOptions() CSSExportOptions {
	return CSSExportOptions{
		Prefix:      "",
		Minify:      false,
		Format:      "css",
		Comments:    true,
		BeginMarker: TwergeBeginMarker,
		EndMarker:   TwergeEndMarker,
	}
}

// ExportCSSWithOptions exports CSS with custom options
func ExportCSSWithOptions(filePath string, options CSSExportOptions) error {
	runtimeMapMutex.RLock()
	defer runtimeMapMutex.RUnlock()

	// Get all keys and sort them for consistent output
	keys := make([]string, 0, len(RuntimeClassMap))
	for k := range RuntimeClassMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var builder strings.Builder

	// Add header comment if comments are enabled
	if options.Comments {
		builder.WriteString("/* Generated by Twerge - https://github.com/conneroisu/twerge */\n")
		builder.WriteString("/* Generated at: " + time.Now().Format(time.RFC3339) + " */\n\n")
	}

	for _, k := range keys {
		original := k
		generated := RuntimeClassMap[k]

		// Apply prefix if specified
		if options.Prefix != "" {
			generated = options.Prefix + generated
		}

		merged := Merge(original)

		// Different formats based on the output format
		switch options.Format {
		case "scss":
			if options.Comments {
				builder.WriteString("// Original: " + original + "\n")
			}
			builder.WriteString(".")
			builder.WriteString(generated)
			builder.WriteString(" {\n  @apply ")
			builder.WriteString(merged)
			builder.WriteString(";\n}\n\n")

		case "less":
			if options.Comments {
				builder.WriteString("// Original: " + original + "\n")
			}
			builder.WriteString(".")
			builder.WriteString(generated)
			builder.WriteString(" {\n  .apply(")
			builder.WriteString(merged)
			builder.WriteString(");\n}\n\n")

		default: // css
			if options.Comments && !options.Minify {
				builder.WriteString("/* Original: " + original + " */\n")
			}
			builder.WriteString(".")
			builder.WriteString(generated)

			if options.Minify {
				builder.WriteString("{@apply ")
				builder.WriteString(merged)
				builder.WriteString(";}")
			} else {
				builder.WriteString(" {\n  @apply ")
				builder.WriteString(merged)
				builder.WriteString(";\n}\n\n")
			}
		}
	}

	cssContent := builder.String()

	// Use custom markers if provided
	beginMarker := options.BeginMarker
	if beginMarker == "" {
		beginMarker = TwergeBeginMarker
	}

	endMarker := options.EndMarker
	if endMarker == "" {
		endMarker = TwergeEndMarker
	}

	// Check if file exists
	fileContent, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading file: %w", err)
	}

	if os.IsNotExist(err) {
		// File doesn't exist, create new content with markers
		newContent := fmt.Appendf(nil, "%s\n%s\n%s\n", beginMarker, cssContent, endMarker)
		return os.WriteFile(filePath, newContent, 0644)
	}
	// File exists, replace content between markers using optimized approach
	// Find begin marker
	beginMarkerBytes := []byte(beginMarker)
	beginIdx := bytes.Index(fileContent, beginMarkerBytes)
	if beginIdx == -1 {
		// Markers don't exist, append content with markers
		suffix := append([]byte("\n\n"), beginMarkerBytes...)
		suffix = append(suffix, '\n')
		suffix = append(suffix, []byte(cssContent)...)
		suffix = append(suffix, '\n')
		suffix = append(suffix, []byte(endMarker)...)
		return os.WriteFile(filePath, append(fileContent, suffix...), 0644)
	}

	// Find the end of the line containing the begin marker
	beginLineEnd := beginIdx + len(beginMarkerBytes)
	for beginLineEnd < len(fileContent) && fileContent[beginLineEnd] != '\n' && fileContent[beginLineEnd] != '\r' {
		beginLineEnd++
	}
	if beginLineEnd < len(fileContent) {
		beginLineEnd++ // Include the newline character
	}

	// Find end marker
	endMarkerBytes := []byte(endMarker)
	endIdx := bytes.Index(fileContent[beginLineEnd:], endMarkerBytes)
	if endIdx == -1 {
		return fmt.Errorf("found begin marker but no end marker")
	}

	// Adjust end marker index to be relative to the whole content
	endIdx += beginLineEnd

	// Create new content with replacement
	result := make([]byte, 0, len(fileContent)-(endIdx-beginLineEnd)+len(cssContent)+1)
	result = append(result, fileContent[:beginLineEnd]...)
	result = append(result, []byte(cssContent)...)
	result = append(result, '\n')
	result = append(result, fileContent[endIdx:]...)

	// Write the new content back to the file
	return os.WriteFile(filePath, result, 0644)
}

// GeneratePostCSSConfig creates a PostCSS configuration file that includes
// Tailwind CSS and other necessary plugins to process Twerge's generated CSS
func GeneratePostCSSConfig(configPath string) error {
	// Basic PostCSS config that works with Tailwind CSS
	config := map[string]interface{}{
		"plugins": []string{
			"tailwindcss",
			"autoprefixer",
		},
	}

	// Convert to JSON
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error creating PostCSS config: %w", err)
	}

	// Write to file
	err = os.WriteFile(configPath, configJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing PostCSS config: %w", err)
	}

	return nil
}

// ExportOptimizedCSS exports the CSS with optimizations for production use
func ExportOptimizedCSS(filePath string) error {
	options := DefaultCSSExportOptions()
	options.Minify = true
	options.Comments = false

	return ExportCSSWithOptions(filePath, options)
}

// WatchAndExportCSS sets up a watcher that exports CSS whenever
// the RuntimeClassMap changes. This is a non-blocking function that returns
// a cancel function to stop the watcher.
//
// Interval is the time between checks in milliseconds.
// If filePath is empty, no file will be written but the callback will still be called.
func WatchAndExportCSS(filePath string, interval int, callback func(string)) (func(), error) {
	if interval < 100 {
		interval = 100 // Minimum 100ms interval to prevent high CPU usage
	}

	// Create a stop channel
	stopCh := make(chan struct{})

	// Store a hash of the last CSS content to detect changes
	var lastHash string

	// Start a goroutine to check for changes
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Get current CSS content
				runtimeMapMutex.RLock()
				cssContent := GetRuntimeClassHTML()
				runtimeMapMutex.RUnlock()

				// Calculate hash
				hash := fmt.Sprintf("%x", sha1.Sum([]byte(cssContent)))

				// Check if content has changed
				if hash != lastHash {
					lastHash = hash

					// Write to file if path is provided
					if filePath != "" {
						if err := WriteCSSToFile(filePath, cssContent); err != nil {
							// Log error but continue
							fmt.Printf("Error writing CSS to file: %v\n", err)
						}
					}

					// Call callback if provided
					if callback != nil {
						callback(cssContent)
					}
				}
			case <-stopCh:
				return
			}
		}
	}()

	// Return a function to stop the watcher
	return func() {
		close(stopCh)
	}, nil
}
