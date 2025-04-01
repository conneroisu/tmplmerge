package twerge

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuntimeGenerate(t *testing.T) {
	// Clear the runtime map for testing
	ClearRuntimeMap()

	// Test that RuntimeGenerate creates a consistent class name for the same input
	class1 := RuntimeGenerate("text-red-500 bg-blue-500")
	class2 := RuntimeGenerate("text-red-500 bg-blue-500")
	assert.Equal(t, class1, class2, "RuntimeGenerate should return the same class name for the same input")

	// Test that RuntimeGenerate handles class merging correctly
	class3 := RuntimeGenerate("text-red-500 text-blue-700")
	assert.NotEqual(t, class1, class3, "RuntimeGenerate should return different class names for different inputs")

	// Test that the generated class name format is correct
	assert.True(t, strings.HasPrefix(class1, "tw-"), "Generated class should start with 'tw-'")
	assert.Equal(t, 10, len(class1), "Generated class should be 10 characters long (tw- + 7 chars)")

	// Check the runtime map has been updated
	mapping := GetRuntimeMapping()
	assert.Equal(t, class1, mapping["text-red-500 bg-blue-500"], "RuntimeClassMap should contain the original class string")
	assert.Equal(t, 2, len(mapping), "RuntimeClassMap should contain 2 entries")
}

func TestRegisterClasses(t *testing.T) {
	// Clear the runtime map for testing
	ClearRuntimeMap()

	// Create some test mappings
	testMap := map[string]string{
		"text-xl font-bold": "tw-custom1",
		"flex items-center": "tw-custom2",
	}

	// Register the classes
	RegisterClasses(testMap)

	// Check the classes were registered
	mapping := GetRuntimeMapping()
	assert.Equal(t, "tw-custom1", mapping["text-xl font-bold"], "RegisterClasses should add the classes to RuntimeClassMap")
	assert.Equal(t, "tw-custom2", mapping["flex items-center"], "RegisterClasses should add the classes to RuntimeClassMap")
	assert.Equal(t, 2, len(mapping), "RuntimeClassMap should contain 2 entries")
}

func TestGetRuntimeClassHTML(t *testing.T) {
	// Clear the runtime map for testing
	ClearRuntimeMap()

	// Register some test classes
	RegisterClasses(map[string]string{
		"text-red-500 bg-blue-500": "tw-test1",
		"p-4 m-2":                  "tw-test2",
	})

	// Get the HTML
	html := GetRuntimeClassHTML()

	// Check the HTML contains the expected content
	assert.Contains(t, html, ".tw-test1 { @apply ", "HTML should contain the generated class name")
	assert.Contains(t, html, ".tw-test2 { @apply ", "HTML should contain the generated class name")

	// Since we're using the actual class names rather than the original classes in ClassMapStr,
	// we should check for the class names themselves
	assert.Contains(t, html, "tw-test1", "HTML should contain the generated class name")
	assert.Contains(t, html, "tw-test2", "HTML should contain the generated class name")

	// Check that the HTML has two lines (one for each class)
	assert.Equal(t, 2, len(strings.Split(strings.TrimSpace(html), "\n")),
		"HTML should contain two lines (one for each class)")
}
