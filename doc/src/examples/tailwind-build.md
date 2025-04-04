# Tailwind Build Integration

This example demonstrates how to integrate Twerge with Tailwind CSS build processes for optimal production workflows.

## Basic Tailwind Integration

This example shows how to use Twerge with the Tailwind CLI:

```go
package main

import (
    "fmt"
    "github.com/conneroisu/twerge"
    "os"
    "os/exec"
)

func main() {
    // 1. Register classes that we want to use
    registerClasses()

    // 2. Generate input CSS file for Tailwind
    err := twerge.GenerateInputCSSForTailwind("tailwind-input.css", "tailwind-output.css")
    if err != nil {
        fmt.Println("Error generating Tailwind input:", err)
        os.Exit(1)
    }

    // 3. Run Tailwind CLI
    cmd := exec.Command("npx", "tailwindcss",
        "-i", "tailwind-input.css",
        "-o", "styles.css",
        "--minify")

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err = cmd.Run()
    if err != nil {
        fmt.Println("Error running Tailwind CLI:", err)
        os.Exit(1)
    }

    // 4. Generate Go code for the mappings
    err = twerge.WriteClassMapFile("tailwind_classes_gen.go")
    if err != nil {
        fmt.Println("Error generating Go code:", err)
        os.Exit(1)
    }

    fmt.Println("Build completed successfully!")
    fmt.Println("- Input CSS: tailwind-input.css")
    fmt.Println("- Output CSS: styles.css")
    fmt.Println("- Go code: tailwind_classes_gen.go")
}

// Register all the classes we want to use
func registerClasses() {
    // Common button variants
    twerge.It("px-4 py-2 rounded font-medium text-white bg-blue-500 hover:bg-blue-600")
    twerge.It("px-4 py-2 rounded font-medium text-white bg-red-500 hover:bg-red-600")
    twerge.It("px-4 py-2 rounded font-medium text-gray-700 bg-gray-200 hover:bg-gray-300")

    // Common layout classes
    twerge.It("flex flex-col min-h-screen")
    twerge.It("flex items-center justify-between p-4")
    twerge.It("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8")

    // Common typography classes
    twerge.It("text-xl font-bold text-gray-900")
    twerge.It("text-sm text-gray-500")
    twerge.It("prose lg:prose-xl")

    // Dark mode variants
    twerge.It("bg-white dark:bg-gray-800")
    twerge.It("text-gray-900 dark:text-white")

    // and many more...
}
```

## Complete Build System Example

This example shows a more complete build system that:

1. Extracts classes from template files
2. Generates optimized CSS with Tailwind
3. Creates a Go package with the class mappings

```go
package main

import (
    "fmt"
    "github.com/conneroisu/twerge"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strings"
)

func main() {
    // 1. Extract classes from templates
    classes := extractClassesFromTemplates("./templates")
    fmt.Printf("Found %d unique class combinations in templates\n", len(classes))

    // 2. Register all classes with Twerge
    for _, class := range classes {
        twerge.It(class)
    }

    // 3. Create build directory
    err := os.MkdirAll("./build", 0755)
    if err != nil {
        fmt.Println("Error creating build directory:", err)
        os.Exit(1)
    }

    // 4. Generate Tailwind input file
    err = twerge.GenerateInputCSSForTailwind("./build/tailwind-input.css", "./build/styles.css")
    if err != nil {
        fmt.Println("Error generating Tailwind input:", err)
        os.Exit(1)
    }

    // 5. Create tailwind.config.js if it doesn't exist
    if _, err := os.Stat("tailwind.config.js"); os.IsNotExist(err) {
        createTailwindConfig()
    }

    // 6. Run Tailwind CLI
    cmd := exec.Command("npx", "tailwindcss",
        "-i", "./build/tailwind-input.css",
        "-o", "./public/css/styles.css",
        "--minify",
        "--config", "./tailwind.config.js")

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err = cmd.Run()
    if err != nil {
        fmt.Println("Error running Tailwind CLI:", err)
        os.Exit(1)
    }

    // 7. Generate Go package with class mappings
    err = os.MkdirAll("./internal/ui/styles", 0755)
    if err != nil {
        fmt.Println("Error creating package directory:", err)
        os.Exit(1)
    }

    // Set package name for the generated code
    twerge.SetGeneratedPackage("styles")
    twerge.SetGeneratedMapName("TailwindClasses")

    // Write the class map file
    err = twerge.WriteClassMapFile("./internal/ui/styles/tailwind_gen.go")
    if err != nil {
        fmt.Println("Error generating Go code:", err)
        os.Exit(1)
    }

    fmt.Println("\nBuild completed successfully!")
    fmt.Println("- Tailwind input: ./build/tailwind-input.css")
    fmt.Println("- CSS output: ./public/css/styles.css")
    fmt.Println("- Go package: ./internal/ui/styles/tailwind_gen.go")
}

// Extract all class combinations from template files
func extractClassesFromTemplates(dir string) []string {
    var classes []string
    classMap := make(map[string]bool)

    // Regex to find class attributes in templates
    // This handles both standard HTML and templ syntax
    re := regexp.MustCompile(`class="([^"]+)"|class=\{[^{]*"([^"]+)"[^}]*\}`)

    // Walk through template directory
    filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Only process template files
        ext := filepath.Ext(path)
        if ext != ".html" && ext != ".tmpl" && ext != ".templ" && ext != ".gohtml" {
            return nil
        }

        // Read file
        content, err := ioutil.ReadFile(path)
        if err != nil {
            return err
        }

        // Find all class attributes
        matches := re.FindAllSubmatch(content, -1)
        for _, match := range matches {
            // Get class string from either capture group
            var classString string
            if len(match[1]) > 0 {
                classString = string(match[1])
            } else if len(match[2]) > 0 {
                classString = string(match[2])
            }

            // Skip if empty
            if classString == "" {
                continue
            }

            // Add to map to ensure uniqueness
            if !classMap[classString] {
                classMap[classString] = true
                classes = append(classes, classString)
            }
        }

        return nil
    })

    return classes
}

// Create a basic Tailwind config file
func createTailwindConfig() {
    config := `module.exports = {
  content: [
    './build/tailwind-input.css',
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
  ],
}`

    ioutil.WriteFile("tailwind.config.js", []byte(config), 0644)
}
```

## Integrating with Package.json

For a complete setup, you can create a `package.json` file to manage Tailwind dependencies:

```json
{
  "name": "my-twerge-project",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "build:css": "go run ./cmd/build/main.go",
    "dev:css": "tailwindcss -i ./build/tailwind-input.css -o ./public/css/styles.css --watch",
    "build": "npm run build:css && go build -o ./bin/app ./cmd/app",
    "dev": "npm run build:css && go run ./cmd/app/main.go"
  },
  "dependencies": {
    "@tailwindcss/forms": "^0.5.3",
    "@tailwindcss/typography": "^0.5.9",
    "tailwindcss": "^3.3.2"
  }
}
```

## Full Project Structure

A typical project structure might look like:

```
project/
├── cmd/
│   ├── app/
│   │   └── main.go         # Main application
│   └── build/
│       └── main.go         # Build script (example above)
├── internal/
│   └── ui/
│       ├── components/     # UI components
│       └── styles/
│           └── tailwind_gen.go  # Generated class mappings
├── templates/              # HTML/templ templates
├── public/
│   └── css/
│       └── styles.css      # Generated CSS
├── build/                  # Build artifacts
│   └── tailwind-input.css  # Generated Tailwind input
├── package.json           # NPM dependencies
└── tailwind.config.js     # Tailwind configuration
```

## CI/CD Integration

You can integrate this build process into CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
name: Build

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.19

      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: 16

      - name: Install dependencies
        run: npm ci

      - name: Build CSS and code
        run: npm run build:css

      - name: Build application
        run: go build -o ./bin/app ./cmd/app

      - name: Run tests
        run: go test ./...
```

See the [tailwind-build example](https://github.com/conneroisu/twerge/tree/main/examples/tailwind-build) in the GitHub repository for a complete working example.
