package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/conneroisu/twerge"
)

func main() {
	// Generate some class names
	class1 := twerge.Merge("text-red-500 bg-blue-500")
	class2 := twerge.Merge("text-green-300 p-4")
	class3 := twerge.Merge("flex items-center justify-between")
	fmt.Println("Generated class names:")
	fmt.Printf("text-red-500 bg-blue-500 -> %s\n", class1)
	fmt.Printf("text-green-300 p-4 -> %s\n", class2)
	fmt.Printf("flex items-center justify-between -> %s\n", class3)

	// Test class merging functionality
	fmt.Println("\nMerged classes:")
	merged := twerge.Merge("text-red-500 text-blue-700")
	fmt.Printf("text-red-500 text-blue-700 -> %s\n", merged)

	// Class conflict resolution
	merged = twerge.Merge("p-4 p-8")
	fmt.Printf("p-4 p-8 -> %s\n", merged)

	// Generate and save the class map code
	outPath := filepath.Join(os.TempDir(), "class_map_generated.go")
	fmt.Printf("\nGenerating class map code to %s\n", outPath)

	// Reset ClassMapStr for a clean start
	twerge.ClassMapStr = make(map[string]string)
	code := twerge.GenerateClassMapCode()
	previewLen := min(len(code), 200)
	fmt.Printf("\nGenerated code preview:\n%s...\n", code[:previewLen])

	err := os.WriteFile(outPath, []byte(code), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	} else {
		fmt.Printf("Class map code written to %s\n", outPath)
	}
}
