package twerge

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceBetweenMarkers(t *testing.T) {
	// Test with existing markers
	original := []byte("Some content\n" + twergeBeginMarker + "\nold content\n" + twergeEndMarker + "\nMore content")
	replacement := []byte("new content")

	result, err := replaceBetweenMarkers(original, replacement)
	assert.NoError(t, err)
	assert.Contains(t, string(result), "new content")
	assert.NotContains(t, string(result), "old content")

	// Test with no markers
	original = []byte("Some content without markers")
	result, err = replaceBetweenMarkers(original, replacement)
	assert.NoError(t, err)
	assert.Contains(t, string(result), twergeBeginMarker)
	assert.Contains(t, string(result), twergeEndMarker)
	assert.Contains(t, string(result), "new content")
}

func TestGenerateInputCSSForTailwind(t *testing.T) {
	// Create temporary input and output files
	inputFile, err := os.CreateTemp("", "twerge-input-*.css")
	if err != nil {
		t.Fatalf("Failed to create temp input file: %v", err)
	}
	defer func() { _ = os.Remove(inputFile.Name()) }()

	// templ file
	templFile, err := os.CreateTemp("", "twerge-templ-*.templ")
	if err != nil {
		t.Fatalf("Failed to create temp templ file: %v", err)
	}
	// defer print(templFile.Name())
	defer func() { _ = os.Remove(templFile.Name()) }()

	// Write some content to the input file
	inputContent := `@tailwind base;
@tailwind components;
@tailwind utilities;

/* Custom styles */
.custom-class {
  color: blue;
}

` + twergeBeginMarker + `
/* Old generated content */
` + twergeEndMarker + `

/* More styles */
`
	err = os.WriteFile(inputFile.Name(), []byte(inputContent), 0644)
	assert.NoError(t, err)

	// Create a test class map
	GenClassMergeStr = map[string]string{
		"tw-test1": "text-red-500",
		"tw-test2": "bg-blue-500",
	}

	// Generate input CSS
	err = GenerateTailwind(inputFile.Name())
	assert.NoError(t, err)

	// Read the output file
	outputContent, err := os.ReadFile(inputFile.Name())
	assert.NoError(t, err)

	// Check content
	outputStr := string(outputContent)
	assert.Contains(t, outputStr, "@tailwind base")
	assert.Contains(t, outputStr, ".custom-class")
	assert.Contains(t, outputStr, ".tw-test1")
	assert.Contains(t, outputStr, ".tw-test2")
	assert.NotContains(t, outputStr, "Old generated content")
	assert.Contains(t, outputStr, "/* More styles */")

	err = GenerateTempl(templFile.Name())
	assert.NoError(t, err)
}
