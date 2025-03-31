// Package main contains an example of the twerge CSS export utilities
// usage in a pretend real-world project.
package main

import (
	"fmt"
	"os"

	"github.com/conneroisu/twerge"
)

// This file is a simple test for our CSS export utilities
func main() {
	// Register some test classes
	twerge.RegisterClasses(map[string]string{
		"flex items-center justify-between":                          "layout-between",
		"bg-white rounded-lg shadow-md p-6":                          "card",
		"px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600": "btn-primary",
	})

	// Test default export
	fmt.Println("Testing default CSS export...")
	err := twerge.ExportCSS("test-default.css")
	if err != nil {
		fmt.Printf("Error with default export: %v\n", err)
		os.Exit(1)
	}

	// Test with options
	fmt.Println("Testing CSS export with options...")
	options := twerge.DefaultCSSExportOptions()
	options.Prefix = "test-"
	options.Format = "scss"
	err = twerge.ExportCSSWithOptions("test-options.scss", options)
	if err != nil {
		fmt.Printf("Error with options export: %v\n", err)
		os.Exit(1)
	}

	// Test optimized export
	fmt.Println("Testing optimized CSS export...")
	err = twerge.ExportOptimizedCSS("test-optimized.css")
	if err != nil {
		fmt.Printf("Error with optimized export: %v\n", err)
		os.Exit(1)
	}

	// Test custom markers
	fmt.Println("Testing custom markers...")
	customOptions := twerge.DefaultCSSExportOptions()
	customOptions.BeginMarker = "/* CUSTOM START */"
	customOptions.EndMarker = "/* CUSTOM END */"

	// Create a file with custom markers
	customContent := `/* Base styles */
body { margin: 0; }

/* CUSTOM START */
/* Will be replaced */
/* CUSTOM END */

/* More styles */
`
	err = os.WriteFile("test-markers.css", []byte(customContent), 0644)
	if err != nil {
		fmt.Printf("Error creating test file: %v\n", err)
		os.Exit(1)
	}

	err = twerge.ExportCSSWithOptions("test-markers.css", customOptions)
	if err != nil {
		fmt.Printf("Error with custom markers: %v\n", err)
		os.Exit(1)
	}

	// Print all generated files
	fmt.Println("\nGenerated files:")

	files := []string{"test-default.css", "test-options.scss", "test-optimized.css", "test-markers.css"}
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file, err)
			continue
		}

		fmt.Printf("=== %s ===\n", file)
		fmt.Println(string(content))
		fmt.Println()
	}

	// Clean up
	fmt.Println("Cleaning up...")
	for _, file := range files {
		_ = os.Remove(file)
	}

	fmt.Println("All tests completed successfully!")
}
