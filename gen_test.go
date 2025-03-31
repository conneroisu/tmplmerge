package twerge

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	// Reset the global class map for testing
	mapMutex.Lock()
	globalClassMap = make(ClassMap)
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
	assert.Equal(t, 10, len(class1), "Generated class should be 10 characters long (tw- + 7 chars)")
}

func TestGetMapping(t *testing.T) {
	// Reset the global class map for testing
	mapMutex.Lock()
	globalClassMap = make(ClassMap)
	mapMutex.Unlock()

	// Generate some class names and store them directly in the global map for testing
	class1 := "tw-abcdefg"
	class2 := "tw-hijklmn"
	
	mapMutex.Lock()
	globalClassMap["text-red-500 bg-blue-500"] = class1
	globalClassMap["text-green-300 p-4"] = class2
	mapMutex.Unlock()

	// Get the mapping
	mapping := GetMapping()

	// Check that the mapping contains the expected entries
	assert.Equal(t, class1, mapping["text-red-500 bg-blue-500"], "Mapping should contain the original class string and generated class name")
	assert.Equal(t, class2, mapping["text-green-300 p-4"], "Mapping should contain the original class string and generated class name")
	assert.Equal(t, 2, len(mapping), "Mapping should contain 2 entries")
}

func TestGenerateClassMapCode(t *testing.T) {
	// Reset the global class map for testing
	mapMutex.Lock()
	globalClassMap = make(ClassMap)
	
	// Store directly in the global map for testing
	globalClassMap["text-red-500 bg-blue-500"] = "tw-abcdefg"
	globalClassMap["text-green-300 p-4"] = "tw-hijklmn"
	mapMutex.Unlock()

	// Generate the code
	code := GenerateClassMapCode()

	// Check that the code contains the expected content
	assert.True(t, strings.Contains(code, "package twerge"), "Generated code should contain package declaration")
	assert.True(t, strings.Contains(code, "ClassMapStr"), "Generated code should contain ClassMapStr variable")
	assert.True(t, strings.Contains(code, `"text-red-500 bg-blue-500"`), "Generated code should contain the original class strings")
	assert.True(t, strings.Contains(code, `"text-green-300 p-4"`), "Generated code should contain the original class strings")
}