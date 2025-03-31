# CSS Integration

Twerge provides tools for integrating with CSS workflows, allowing you to export class mappings to CSS files and work with Tailwind CLI.

## CSS Export Functions

Twerge offers several functions for working with CSS files:

### Exporting CSS with Markers

```go
import "github.com/conneroisu/twerge"

// Export CSS to a file between twerge markers
err := twerge.ExportCSS("styles.css")
if err != nil {
    // Handle error
}
```

This will write the CSS for all registered classes between markers:

```css
/* styles.css */

/* Content before the markers is preserved */

<!-- twerge:begin -->
.tw-a1b2c3d4 {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem;
}
.tw-e5f6g7h8 {
    font-size: 1.125rem;
    font-weight: 700;
    color: #1f2937;
}
<!-- twerge:end -->

/* Content after the markers is preserved */
```

### Exporting with a Specific Class Map

```go
// Export CSS from a specific map
componentMap := map[string]string{
    "flex items-center p-2": "component-base",
    "bg-blue-500 text-white": "component-primary",
}

err := twerge.ExportCSSWithMap("component-styles.css", componentMap)
if err != nil {
    // Handle error
}
```

### Appending Classes to Files

```go
// Append classes to an existing CSS file
newClasses := map[string]string{
    "grid grid-cols-3 gap-4": "grid-layout",
}

err := twerge.AppendClassesToFile("styles.css", newClasses, "/* Grid Components */")
if err != nil {
    // Handle error
}
```

## Tailwind CLI Integration

Twerge can generate input files for Tailwind CLI, allowing you to leverage Tailwind's processing:

### Generating Input CSS for Tailwind

```go
// Creates a Tailwind CLI input file
err := twerge.GenerateInputCSSForTailwind("input.css", "output.css")
if err != nil {
    // Handle error
}
```

This generates a file that includes all the class selectors needed for Tailwind CLI to process:

```css
/* input.css */

@tailwind base;
@tailwind components;
@tailwind utilities;

/* twerge-generated class selectors */
.tw-a1b2c3d4 {
}
.tw-e5f6g7h8 {
}
/* end twerge-generated class selectors */
```

You can then run Tailwind CLI:

```bash
npx tailwindcss -i input.css -o output.css
```

### Processing CSS Templates

```go
// Process a CSS template with twerge markers
err := twerge.ProcessCSSTemplate("template.css", "processed.css")
if err != nil {
    // Handle error
}
```

This allows you to have a template CSS file with markers that Twerge will fill in.

## Integration Examples

### Server-Side Rendering with Runtime CSS

```go
package main

import (
    "fmt"
    "github.com/conneroisu/twerge"
    "net/http"
)

func main() {
    // Pre-register some common classes
    twerge.InitWithCommonClasses()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Generate HTML with runtime class generation
        html := `
        <!DOCTYPE html>
        <html>
        <head>
            <style>
            ` + twerge.GetRuntimeClassHTML() + `
            </style>
        </head>
        <body>
            <div class="` + twerge.RuntimeGenerate("flex items-center justify-between p-4") + `">
                <h1 class="` + twerge.RuntimeGenerate("text-2xl font-bold") + `">Hello World</h1>
                <button class="` + twerge.RuntimeGenerate("bg-blue-500 text-white px-4 py-2 rounded") + `">Click Me</button>
            </div>
        </body>
        </html>
        `

        fmt.Fprintf(w, html)
    })

    http.ListenAndServe(":8080", nil)
}
```

### Build-Time CSS Generation

```go
package main

import (
    "github.com/conneroisu/twerge"
    "log"
)

func main() {
    // Register all the classes used in your application
    twerge.Generate("flex items-center justify-between p-4")
    twerge.Generate("text-2xl font-bold")
    twerge.Generate("bg-blue-500 text-white px-4 py-2 rounded")

    // Export the CSS for Tailwind processing
    err := twerge.GenerateInputCSSForTailwind("tailwind-input.css", "tailwind-output.css")
    if err != nil {
        log.Fatal(err)
    }

    // Generate Go code for the class mappings
    err = twerge.WriteClassMapFile("tailwind_classes_gen.go")
    if err != nil {
        log.Fatal(err)
    }

    log.Println("CSS and Go code generated successfully")
}
```

## Advanced CSS Integration

### Custom CSS Markers

You can customize the markers used in CSS files:

```go
// Set custom CSS markers
twerge.SetCSSMarkers("/* TWERGE-START */", "/* TWERGE-END */")
```

### Combining Multiple CSS Files

```go
// Generate component-specific CSS files
twerge.ExportCSSWithMap("buttons.css", buttonClasses)
twerge.ExportCSSWithMap("layout.css", layoutClasses)
twerge.ExportCSSWithMap("typography.css", typographyClasses)

// Merge the CSS files (using standard tools)
// cat buttons.css layout.css typography.css > combined.css
```
