package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/conneroisu/twerge"
)

func main() {
	// Initialize with common class combinations
	twerge.InitWithCommonClasses()
	
	// Register some custom classes
	twerge.RegisterClasses(map[string]string{
		"bg-gradient-to-r from-purple-500 to-pink-500 text-white p-4": "tw-gradient",
		"bg-black text-white p-6 rounded-lg shadow-xl": "tw-card-dark",
		"grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6": "tw-responsive-grid",
	})
	
	// Generate a new class name at runtime
	customClass := twerge.RuntimeGenerate("text-sm font-medium text-gray-700")
	
	// Print all the registered class mappings
	fmt.Println("Registered Class Mappings:")
	fmt.Println("==========================")
	
	mapping := twerge.GetRuntimeMapping()
	
	// Print common classes
	fmt.Println("\nCommon Classes:")
	printSection(mapping, "tw-layout", "tw-text", "tw-btn", "tw-card", "tw-input")
	
	// Print custom classes
	fmt.Println("\nCustom Classes:")
	printSection(mapping, "tw-gradient", "tw-card-dark", "tw-responsive-grid")
	
	// Print dynamically generated classes
	fmt.Println("\nDynamically Generated Classes:")
	fmt.Printf("text-sm font-medium text-gray-700 -> %s\n", customClass)
	
	// Generate HTML for CSS styles
	html := twerge.GetRuntimeClassHTML()
	
	// Save the HTML to a file
	err := os.WriteFile("styles.css", []byte(html), 0644)
	if err != nil {
		fmt.Printf("Error writing CSS file: %v\n", err)
	} else {
		fmt.Println("\nGenerated CSS file: styles.css")
	}
	
	// Show examples of how to use these classes in HTML
	fmt.Println("\nHTML Usage Examples:")
	fmt.Println("===================")
	fmt.Println("<div class=\"tw-layout1\">Centered flex container</div>")
	fmt.Println("<button class=\"tw-btn1\">Primary Button</button>")
	fmt.Println("<div class=\"tw-gradient\">Gradient Background</div>")
	fmt.Println("<div class=\"" + customClass + "\">Small Text</div>")
}

// printSection prints a subset of the mapping
func printSection(mapping map[string]string, prefixes ...string) {
	for k, v := range mapping {
		for _, prefix := range prefixes {
			if strings.HasPrefix(v, prefix) {
				fmt.Printf("%s -> %s\n", k, v)
				break
			}
		}
	}
}