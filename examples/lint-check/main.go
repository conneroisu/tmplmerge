package main

import (
	"fmt"

	"github.com/conneroisu/twerge"
)

func main() {
	// Register some class combinations with known merges
	classesToMerge := []string{
		"bg-red-500 text-white p-4",
		"bg-blue-500 text-white p-4",
		"bg-green-500 hover:bg-blue-700 text-white p-4",
		"p-2 p-4 text-lg",
		"pt-4 pb-4 pl-4 pr-4",
		"text-sm text-lg",
		"text-base text-lg",
		"inline block",
		"hidden block",
		"p-4 py-4 px-4",
		"m-2 mx-2 my-2",
	}

	fmt.Println("Merging class combinations...")
	for _, classes := range classesToMerge {
		merged := twerge.Merge(classes)
		fmt.Printf("Original: %-40s -> Merged: %s\n", classes, merged)
	}
	
	fmt.Println("\nRunning lint check to find duplicates...")
	fmt.Println(twerge.LintString())
	
	reports := twerge.Lint()
	fmt.Printf("Found %d cases where different class combinations merge to the same final classes.\n", len(reports))
	
	fmt.Println("\nThis helps you identify redundant classes in your codebase that could be simplified.")
	fmt.Println("For example, you could replace:")
	fmt.Println("  'pt-4 pb-4 pl-4 pr-4' with 'p-4'")
	fmt.Println("  'inline block' or 'hidden block' with just 'block'")
}