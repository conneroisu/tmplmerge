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
	// First check if a class name exists in ClassMapStr for quick lookup
	if className, exists := ClassMapStr[classes]; exists {
		return className
	}

	// Then check if we already have a class name for this string in RuntimeClassMap
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
// Also updates ClassMapStr for quick lookups
func RegisterClasses(classMap map[string]string) {
	runtimeMapMutex.Lock()
	defer runtimeMapMutex.Unlock()

	maps.Copy(RuntimeClassMap, classMap)

	// Also update ClassMapStr for quick lookups
	maps.Copy(ClassMapStr, classMap)
}

// ClearRuntimeMap clears the RuntimeClassMap and ClassMapStr
// Useful for testing or resetting the maps
func ClearRuntimeMap() {
	runtimeMapMutex.Lock()
	defer runtimeMapMutex.Unlock()

	RuntimeClassMap = make(map[string]string)
	ClassMapStr = make(map[string]string)
}

// GetRuntimeMapping returns a copy of the current RuntimeClassMap
func GetRuntimeMapping() map[string]string {
	runtimeMapMutex.RLock()
	defer runtimeMapMutex.RUnlock()

	// Create a copy to avoid concurrent map access issues
	mapping := make(map[string]string, len(RuntimeClassMap))
	maps.Copy(mapping, RuntimeClassMap)

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
