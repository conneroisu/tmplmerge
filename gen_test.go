package twerge

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	// Reset the class map for testing
	mapMutex.Lock()
	ClassMapStr = make(map[string]string)
	mapMutex.Unlock()

	// Test that Generate creates a consistent class name for the same input
	class1 := Generate("text-red-500 bg-blue-500")
	class2 := Generate("text-red-500 bg-blue-500")
	assert.Equal(t, class1, class2, "Generate should return the same class name for the same input")

	// Test that Generate handles class merging correctly
	class3 := Generate("text-red-500 text-blue-700")
	assert.NotEqual(t, class1, class3, "Generate should return different class names for different inputs")

	// Test that the generated class name format is correct
	assert.True(t, strings.HasPrefix(class1, "tw-"), "Generated class should start with 'tw-'")
}

func TestGetMapping(t *testing.T) {
	// Reset the class map for testing
	mapMutex.Lock()
	ClassMapStr = make(map[string]string)
	mapMutex.Unlock()

	// Generate some class names and store them directly in the map for testing
	class1 := "tw-abcdefg"
	class2 := "tw-hijklmn"

	mapMutex.Lock()
	ClassMapStr["text-red-500 bg-blue-500"] = class1
	ClassMapStr["text-green-300 p-4"] = class2
	mapMutex.Unlock()

	// Get the mapping
	mapping := getMapping()

	// Check that the mapping contains the expected entries
	assert.Equal(t, class1, mapping["text-red-500 bg-blue-500"], "Mapping should contain the original class string and generated class name")
	assert.Equal(t, class2, mapping["text-green-300 p-4"], "Mapping should contain the original class string and generated class name")
	assert.Equal(t, 2, len(mapping), "Mapping should contain 2 entries")
}

func TestGenerateClassMapCode(t *testing.T) {
	// Reset the class map for testing
	mapMutex.Lock()
	ClassMapStr = make(map[string]string)

	// Store directly in the map for testing
	ClassMapStr["text-red-500 bg-blue-500"] = "tw-abcdefg"
	ClassMapStr["text-green-300 p-4"] = "tw-hijklmn"
	mapMutex.Unlock()

	// Generate the code
	code := GenerateClassMapCode("twerge")

	// Check that the code contains the expected content
	assert.True(t, strings.Contains(code, "package twerge"), "Generated code should contain package declaration")
	assert.True(t, strings.Contains(code, "ClassMapStr"), "Generated code should contain ClassMapStr variable")
	assert.True(t, strings.Contains(code, `"text-red-500 bg-blue-500"`), "Generated code should contain the original class strings")
	assert.True(t, strings.Contains(code, `"text-green-300 p-4"`), "Generated code should contain the original class strings")
}
