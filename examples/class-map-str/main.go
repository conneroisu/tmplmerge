package main

import (
	"fmt"
	"time"

	"maps"

	"github.com/conneroisu/twerge"
)

func main() {
	// Populate the ClassMapStr with some frequently used class combinations
	twerge.ClassMapStr = map[string]string{
		"flex items-center justify-center":   "tw-header",
		"p-4 bg-blue-500 text-white rounded": "tw-button",
		"grid grid-cols-3 gap-4":             "tw-grid3",
		"text-xl font-bold text-gray-900":    "tw-title",
	}

	// Example 1: Direct lookup from ClassMapStr
	fmt.Println("Example 1: Direct lookup from ClassMapStr")
	start := time.Now()
	result1 := twerge.Merge("flex items-center justify-center")
	elapsed1 := time.Since(start)
	fmt.Printf("Input: \"flex items-center justify-center\"\n")
	fmt.Printf("Output: \"%s\"\n", result1)
	fmt.Printf("Time: %s (fast, direct lookup)\n\n", elapsed1)

	// Example 2: Class that needs merging (not in ClassMapStr)
	fmt.Println("Example 2: Class that needs merging")
	start = time.Now()
	result2 := twerge.Merge("p-4 bg-red-500 p-6") // p-6 should override p-4
	elapsed2 := time.Since(start)
	fmt.Printf("Input: \"p-4 bg-red-500 p-6\"\n")
	fmt.Printf("Output: \"%s\"\n", result2)
	fmt.Printf("Time: %s (slower, required merging)\n\n", elapsed2)

	// Example 3: Adding more entries to ClassMapStr
	additionalClasses := map[string]string{
		"text-sm text-gray-500":   "tw-subtitle",
		"flex flex-col space-y-4": "tw-colstack",
	}

	// Update the ClassMapStr map directly
	maps.Copy(twerge.ClassMapStr, additionalClasses)

	// This uses ClassMapStr for quick lookup
	fmt.Println("Example 3: Using ClassMapStr for lookups")
	result3 := twerge.Merge("text-xl font-bold text-gray-900") // From ClassMapStr
	result4 := twerge.Generate("text-sm text-gray-500")        // From ClassMapStr

	fmt.Printf("From ClassMapStr lookup: \"%s\"\n", result3)
	fmt.Printf("From Generate: \"%s\"\n", result4)

	// Example 4: Auto-generate code for ClassMapStr
	fmt.Println("\nExample 4: Auto-generate code for ClassMapStr")
	fmt.Println("Code that would be written to a file:")
	fmt.Println("----------------------------------------")
	code := twerge.GenerateClassMapCode()
	fmt.Println(code)
	fmt.Println("----------------------------------------")
	fmt.Println("This code can be written to a file with WriteClassMapFile()")
}
