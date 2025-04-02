package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"maps"

	"github.com/conneroisu/twerge"
)

func main() {
	// Create the directory structure
	createDirectories()

	// Step 1: Create class map
	fmt.Println("Creating class map...")
	classMap := createClassMap()

	// Add to ClassMapStr for other operations
	maps.Copy(twerge.ClassMapStr, classMap)
	twerge.GenClassMergeStr = twerge.ClassMapStr

	// Step 2: Generate the input CSS file for Tailwind CLI
	fmt.Println("Generating input CSS file...")
	inputCSSPath := filepath.Join("src", "input.css")
	err := twerge.GenerateTailwind(inputCSSPath)
	if err != nil {
		log.Fatalf("Error generating input CSS: %v", err)
	}

	// Step 3: Run the Tailwind CLI to generate the output CSS
	fmt.Println("Running Tailwind CLI...")
	runTailwindCLI()

	// Step 4: Generate a sample HTML file to demonstrate the classes
	fmt.Println("Generating sample HTML...")
	generateSampleHTML()

	fmt.Println("Done! Open dist/index.html in your browser to see the results.")
}

func createDirectories() {
	// Create necessary directories
	dirs := []string{
		"src",
		"dist",
		"dist/css",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Error creating directory %s: %v", dir, err)
		}
	}

	// Create Tailwind config if it doesn't exist
	configPath := "tailwind.config.js"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := `/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{html,js}", "./dist/**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [],
}
`
		err = os.WriteFile(configPath, []byte(config), 0644)
		if err != nil {
			log.Fatalf("Error creating Tailwind config: %v", err)
		}
	}

	// Create initial input.css if it doesn't exist
	inputCSSPath := filepath.Join("src", "input.css")
	if _, err := os.Stat(inputCSSPath); os.IsNotExist(err) {
		inputCSS := `@tailwind base;
@tailwind components;
@tailwind utilities;

/* Custom CSS goes here */

/* twerge:begin */
/* Twerge generated classes will be inserted here */
/* twerge:end */
`
		err = os.WriteFile(inputCSSPath, []byte(inputCSS), 0644)
		if err != nil {
			log.Fatalf("Error creating input CSS: %v", err)
		}
	}
}

func createClassMap() map[string]string {
	// Create a map of common utility classes
	return map[string]string{
		// Layout
		"flex items-center justify-between":                    "layout-between",
		"flex items-center justify-center":                     "layout-center",
		"grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6": "grid-responsive",

		// Components
		"bg-white rounded-lg shadow-md p-6":                               "card",
		"bg-white rounded-md shadow-sm hover:shadow-md transition-shadow": "card-hover",

		// Buttons
		"px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600":    "btn-primary",
		"px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300": "btn-secondary",
		"px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600":      "btn-danger",

		// Typography
		"text-xl font-bold text-gray-900": "text-title",
		"text-sm text-gray-500":           "text-muted",

		// Forms
		"block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500": "input",

		// Navigation
		"bg-white shadow": "navbar",
		"px-6 py-3 flex items-center justify-between": "navbar-inner",
	}
}

func runTailwindCLI() {
	// Make sure Tailwind CLI is installed
	checkTailwindInstallation()

	// Run the Tailwind CLI command
	inputPath := filepath.Join("src", "input.css")
	outputPath := filepath.Join("dist", "css", "styles.css")

	cmd := exec.Command("npx", "tailwindcss", "-i", inputPath, "-o", outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running Tailwind CLI: %v\nOutput: %s", err, output)
	}
	fmt.Printf("Tailwind CSS output: %s\n", output)
}

func checkTailwindInstallation() {
	// Check if Tailwind CLI is installed
	cmd := exec.Command("npx", "tailwindcss", "--help")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Tailwind CLI is not installed. Installing...")

		// Install Tailwind CSS
		installCmd := exec.Command("npm", "install", "-D", "tailwindcss")
		output, err := installCmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Error installing Tailwind CSS: %v\nOutput: %s", err, output)
		}
		fmt.Printf("Tailwind CSS installed: %s\n", output)
	}
}

func generateSampleHTML() {
	// Generate a sample HTML file
	htmlContent := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Twerge with Tailwind CSS Example</title>
    <link rel="stylesheet" href="css/styles.css">
</head>
<body class="bg-gray-50 min-h-screen">
    <!-- Navbar using Twerge class names -->
    <nav class="navbar">
        <div class="navbar-inner max-w-7xl mx-auto">
            <div class="text-title">Twerge + Tailwind</div>
            <div class="flex space-x-4">
                <a href="#" class="text-blue-500 hover:text-blue-700">Home</a>
                <a href="#" class="text-blue-500 hover:text-blue-700">About</a>
                <a href="#" class="text-blue-500 hover:text-blue-700">Contact</a>
            </div>
        </div>
    </nav>

    <!-- Main content -->
    <main class="max-w-7xl mx-auto py-12 px-4">
        <h1 class="text-3xl font-bold text-center mb-12">Twerge with Tailwind CSS Example</h1>
        
        <!-- Cards section using Twerge classes -->
        <div class="grid-responsive">
            <!-- Card 1 -->
            <div class="card">
                <h2 class="text-title mb-4">Utility Class Merging</h2>
                <p class="mb-4">Twerge merges Tailwind utility classes to resolve conflicts and reduce duplication.</p>
                <button class="btn-primary">Learn More</button>
            </div>
            
            <!-- Card 2 -->
            <div class="card">
                <h2 class="text-title mb-4">Class Name Generation</h2>
                <p class="mb-4">Generate short, unique class names for your Tailwind utilities to reduce CSS size.</p>
                <button class="btn-secondary">See Example</button>
            </div>
            
            <!-- Card 3 -->
            <div class="card">
                <h2 class="text-title mb-4">Build Integration</h2>
                <p class="mb-4">Seamlessly integrate with your build process to optimize CSS production.</p>
                <button class="btn-danger">View Docs</button>
            </div>
        </div>
        
        <!-- Form section using Twerge classes -->
        <div class="card mt-12">
            <h2 class="text-title mb-6">Contact Us</h2>
            <form class="space-y-4">
                <div>
                    <label for="name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
                    <input type="text" id="name" name="name" class="input">
                </div>
                <div>
                    <label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email</label>
                    <input type="email" id="email" name="email" class="input">
                </div>
                <div>
                    <label for="message" class="block text-sm font-medium text-gray-700 mb-1">Message</label>
                    <textarea id="message" name="message" rows="4" class="input"></textarea>
                </div>
                <div class="layout-between">
                    <button type="button" class="btn-secondary">Cancel</button>
                    <button type="submit" class="btn-primary">Send Message</button>
                </div>
            </form>
        </div>
    </main>
    
    <!-- Footer -->
    <footer class="bg-gray-800 text-white py-12 mt-12">
        <div class="max-w-7xl mx-auto px-4">
            <div class="layout-between flex-col md:flex-row">
                <div class="mb-8 md:mb-0">
                    <h3 class="text-lg font-semibold mb-4">Twerge + Tailwind</h3>
                    <p class="text-muted text-gray-300">Optimizing Tailwind CSS for production.</p>
                </div>
                <div class="text-center md:text-right">
                    <p class="text-muted text-gray-300">© 2025 Twerge. All rights reserved.</p>
                </div>
            </div>
        </div>
    </footer>
</body>
</html>`

	htmlPath := filepath.Join("dist", "index.html")
	err := os.WriteFile(htmlPath, []byte(htmlContent), 0644)
	if err != nil {
		log.Fatalf("Error generating sample HTML: %v", err)
	}
}
