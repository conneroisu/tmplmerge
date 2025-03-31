package twerge

import (
	"crypto/sha1"
	"encoding/base64"
	"maps"
	"sort"
	"strings"
	"sync"
)

// RuntimeClassMap is a global runtime map of original classes to generated class names
var RuntimeClassMap = make(map[string]string)

// runtimeMapMutex protects the RuntimeClassMap
var runtimeMapMutex sync.RWMutex

// RuntimeGenerate creates a unique class name and stores it in the RuntimeClassMap
// Similar to Generate but updates the RuntimeClassMap instead of globalClassMap
func RuntimeGenerate(classes string) string {
	// First check if we already have a class name for this string
	runtimeMapMutex.RLock()
	if className, exists := RuntimeClassMap[classes]; exists {
		runtimeMapMutex.RUnlock()
		return className
	}
	runtimeMapMutex.RUnlock()

	// Merge the classes first
	merged := Merge(classes)

	// Generate a hash of the merged classes
	hash := sha1.Sum([]byte(merged))

	// Use URL-safe base64 encoding and trim to 7 characters for brevity
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	className := "tw-" + encoded[:7]

	// Store in the runtime map
	runtimeMapMutex.Lock()
	RuntimeClassMap[classes] = className
	runtimeMapMutex.Unlock()

	return className
}

// RegisterClasses pre-registers a batch of classes with their generated class names
// This is useful for initializing the RuntimeClassMap with known values
func RegisterClasses(classMap map[string]string) {
	runtimeMapMutex.Lock()
	defer runtimeMapMutex.Unlock()

	maps.Copy(RuntimeClassMap, classMap)
}

// ClearRuntimeMap clears the RuntimeClassMap
// Useful for testing or resetting the map
func ClearRuntimeMap() {
	runtimeMapMutex.Lock()
	defer runtimeMapMutex.Unlock()

	RuntimeClassMap = make(map[string]string)
}

// GetRuntimeMapping returns a copy of the current RuntimeClassMap
func GetRuntimeMapping() map[string]string {
	runtimeMapMutex.RLock()
	defer runtimeMapMutex.RUnlock()

	// Create a copy to avoid concurrent map access issues
	mapping := make(map[string]string, len(RuntimeClassMap))
	for k, v := range RuntimeClassMap {
		mapping[k] = v
	}

	return mapping
}

// GetRuntimeClassHTML generates HTML for the RuntimeClassMap
// Useful for injecting CSS classes directly into a template
func GetRuntimeClassHTML() string {
	runtimeMapMutex.RLock()
	defer runtimeMapMutex.RUnlock()

	// Get all keys and sort them for consistent output
	keys := make([]string, 0, len(RuntimeClassMap))
	for k := range RuntimeClassMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, k := range keys {
		original := k
		generated := RuntimeClassMap[k]
		merged := Merge(original)

		// Create a CSS rule using the generated class name and the merged Tailwind classes
		builder.WriteString(".")
		builder.WriteString(generated)
		builder.WriteString(" { @apply ")
		builder.WriteString(merged)
		builder.WriteString("; }\n")
	}

	return builder.String()
}

// InitWithCommonClasses pre-populates the RuntimeClassMap with commonly used class combinations
func InitWithCommonClasses() {
	commonClasses := map[string]string{
		// Layout
		"flex items-center justify-center": "tw-layout1",
		"flex flex-col space-y-4":          "tw-layout2",
		"grid grid-cols-3 gap-4":           "tw-layout3",
		"container mx-auto px-4":           "tw-layout4",

		// Typography
		"text-xl font-bold text-gray-900":                   "tw-text1",
		"text-sm text-gray-500":                             "tw-text2",
		"font-medium text-indigo-600 hover:text-indigo-500": "tw-text3",

		// Buttons
		"px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600":    "tw-btn1",
		"px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300": "tw-btn2",
		"px-4 py-2 border border-gray-300 rounded shadow-sm":            "tw-btn3",

		// Cards
		"bg-white rounded-lg shadow-md p-6": "tw-card1",
		"bg-gray-100 rounded-lg p-4":        "tw-card2",

		// Form elements
		"block w-full rounded-md border-gray-300 shadow-sm":      "tw-input1",
		"mt-1 block w-full rounded-md border-gray-300 shadow-sm": "tw-input2",
	}

	RegisterClasses(commonClasses)
}

