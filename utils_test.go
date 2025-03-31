package twerge

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteCSSToFile(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "twerge-css-test-*.css")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tempFile.Name()) }()
	_ = tempFile.Close()

	// Test content
	cssContent := ".test-class { @apply text-red-500; }"

	// Write to the file
	err = WriteCSSToFile(tempFile.Name(), cssContent)
	assert.NoError(t, err)

	// Read back the file
	content, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	// Check that the markers and content are there
	assert.Contains(t, string(content), TwergeBeginMarker)
	assert.Contains(t, string(content), TwergeEndMarker)
	assert.Contains(t, string(content), cssContent)
}

func TestReplaceBetweenMarkers(t *testing.T) {
	// Test with existing markers
	original := []byte("Some content\n" + TwergeBeginMarker + "\nold content\n" + TwergeEndMarker + "\nMore content")
	replacement := []byte("new content")

	result, err := replaceBetweenMarkers(original, replacement)
	assert.NoError(t, err)
	assert.Contains(t, string(result), "new content")
	assert.NotContains(t, string(result), "old content")

	// Test with no markers
	original = []byte("Some content without markers")
	result, err = replaceBetweenMarkers(original, replacement)
	assert.NoError(t, err)
	assert.Contains(t, string(result), TwergeBeginMarker)
	assert.Contains(t, string(result), TwergeEndMarker)
	assert.Contains(t, string(result), "new content")
}

func TestExportCSS(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "twerge-export-test-*.css")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tempFile.Name()) }()
	_ = tempFile.Close()

	// Register some test classes
	ClearRuntimeMap()
	RegisterClasses(map[string]string{
		"text-red-500 font-bold": "tw-test1",
		"bg-blue-500 p-4":        "tw-test2",
	})

	// Export to file
	err = ExportCSS(tempFile.Name())
	assert.NoError(t, err)

	// Read back the file
	content, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	// Check content
	contentStr := string(content)
	assert.Contains(t, contentStr, ".tw-test1")
	assert.Contains(t, contentStr, ".tw-test2")
	assert.Contains(t, contentStr, "@apply")
}

func TestGenerateInputCSSForTailwind(t *testing.T) {
	// Create temporary input and output files
	inputFile, err := os.CreateTemp("", "twerge-input-*.css")
	if err != nil {
		t.Fatalf("Failed to create temp input file: %v", err)
	}
	defer func() { _ = os.Remove(inputFile.Name()) }()
	
	outputFile, err := os.CreateTemp("", "twerge-output-*.css")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer func() { _ = os.Remove(outputFile.Name()) }()
	
	// Write some content to the input file
	inputContent := `@tailwind base;
@tailwind components;
@tailwind utilities;

/* Custom styles */
.custom-class {
  color: blue;
}

` + TwergeBeginMarker + `
/* Old generated content */
` + TwergeEndMarker + `

/* More styles */
`
	err = os.WriteFile(inputFile.Name(), []byte(inputContent), 0644)
	assert.NoError(t, err)
	
	// Register some test classes
	ClearRuntimeMap()
	RegisterClasses(map[string]string{
		"text-red-500 font-bold": "tw-test1",
		"bg-blue-500 p-4":        "tw-test2",
	})
	
	// Generate input CSS
	err = GenerateInputCSSForTailwind(inputFile.Name(), outputFile.Name())
	assert.NoError(t, err)
	
	// Read the output file
	outputContent, err := os.ReadFile(outputFile.Name())
	assert.NoError(t, err)
	
	// Check content
	outputStr := string(outputContent)
	assert.Contains(t, outputStr, "@tailwind base")
	assert.Contains(t, outputStr, ".custom-class")
	assert.Contains(t, outputStr, ".tw-test1")
	assert.Contains(t, outputStr, ".tw-test2")
	assert.NotContains(t, outputStr, "Old generated content")
	assert.Contains(t, outputStr, "/* More styles */")
}

func TestMergeCSSMaps(t *testing.T) {
	map1 := map[string]string{
		"class1": "tw-1",
		"class2": "tw-2",
	}
	
	map2 := map[string]string{
		"class2": "tw-2-new", // This should override the previous value
		"class3": "tw-3",
	}
	
	merged := MergeCSSMaps(map1, map2)
	
	assert.Equal(t, 3, len(merged))
	assert.Equal(t, "tw-1", merged["class1"])
	assert.Equal(t, "tw-2-new", merged["class2"]) // Should have the updated value
	assert.Equal(t, "tw-3", merged["class3"])
}

func TestAppendClassesToFile(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "twerge-append-test-*.css")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tempFile.Name()) }()
	_ = tempFile.Close()
	
	// Test classes
	classes := map[string]string{
		"text-red-500 font-bold": "test1",
		"bg-blue-500 p-4":        "test2",
	}
	
	// Custom prefix
	prefix := "custom-"
	
	// Append to file
	err = AppendClassesToFile(tempFile.Name(), classes, prefix)
	assert.NoError(t, err)
	
	// Read back the file
	content, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)
	
	// Check content
	contentStr := string(content)
	assert.Contains(t, contentStr, ".custom-test1")
	assert.Contains(t, contentStr, ".custom-test2")
	assert.Contains(t, contentStr, "@apply")
	
	// Count occurrences of class definitions
	classDefCount := strings.Count(contentStr, "@apply")
	assert.Equal(t, 2, classDefCount)
}