# Basic Examples

This page provides basic examples of how to use Twerge in your Go applications.

## Class Merging Example

```go
package main

import (
    "fmt"
    "github.com/conneroisu/twerge"
)

func main() {
    // Merge conflicting Tailwind classes
    classes := "text-red-500 bg-blue-300 text-xl"
    mergedClasses := twerge.Merge(classes)

    fmt.Println("Original:", classes)
    fmt.Println("Merged:", mergedClasses)

    // Conflict resolution
    classes = "p-4 m-2 p-8"
    mergedClasses = twerge.Merge(classes)

    fmt.Println("\nOriginal:", classes)
    fmt.Println("Merged:", mergedClasses) // p-8 should win
}
```

Output:

```
Original: text-red-500 bg-blue-300 text-xl
Merged: bg-blue-300 text-xl text-red-500

Original: p-4 m-2 p-8
Merged: m-2 p-8
```

## Class Generation Example

```go
package main

import (
    "fmt"
    "github.com/conneroisu/twerge"
)

func main() {
    // Generate short class names
    classes := "flex items-center justify-between p-4 bg-white"
    shortName := twerge.Generate(classes)

    fmt.Println("Original:", classes)
    fmt.Println("Generated:", shortName)

    // Get the mapping
    mapping := twerge.GetMapping()
    fmt.Println("\nClass Mapping:")
    for original, generated := range mapping {
        fmt.Printf("%s => %s\n", original, generated)
    }
}
```

Output:

```
Original: flex items-center justify-between p-4 bg-white
Generated: tw-a1b2c3d4

Class Mapping:
flex items-center justify-between p-4 bg-white => tw-a1b2c3d4
```

## Runtime Mapping Example

```go
package main

import (
    "fmt"
    "github.com/conneroisu/twerge"
)

func main() {
    // Register custom classes
    customClasses := map[string]string{
        "flex items-center justify-between": "tw-header",
        "text-sm font-medium text-gray-500": "tw-label",
    }
    twerge.RegisterClasses(customClasses)

    // Use the registered classes
    class1 := twerge.RuntimeGenerate("flex items-center justify-between")
    class2 := twerge.RuntimeGenerate("text-sm font-medium text-gray-500")
    class3 := twerge.RuntimeGenerate("p-4 bg-white") // Not pre-registered

    fmt.Println("Custom Header:", class1)
    fmt.Println("Custom Label:", class2)
    fmt.Println("Generated:", class3)

    // Get runtime HTML/CSS
    html := twerge.GetRuntimeClassHTML()
    fmt.Println("\nGenerated HTML/CSS:")
    fmt.Println(html)
}
```

Output:

```
Custom Header: tw-header
Custom Label: tw-label
Generated: tw-e5f6g7h8

Generated HTML/CSS:
<style>
.tw-header { display: flex; align-items: center; justify-content: space-between; }
.tw-label { font-size: 0.875rem; font-weight: 500; color: #6b7280; }
.tw-e5f6g7h8 { padding: 1rem; background-color: #ffffff; }
</style>
```

## CSS Export Example

```go
package main

import (
    "fmt"
    "github.com/conneroisu/twerge"
    "os"
)

func main() {
    // Generate some classes
    twerge.Generate("flex items-center p-4")
    twerge.Generate("text-xl font-bold text-gray-800")
    twerge.Generate("bg-blue-500 text-white rounded")

    // Export to CSS file
    err := twerge.ExportCSS("styles.css")
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    fmt.Println("CSS file exported successfully to styles.css")

    // Generate Tailwind input
    err = twerge.GenerateInputCSSForTailwind("tailwind-input.css", "tailwind-output.css")
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    fmt.Println("Tailwind input file generated successfully")
}
```

## Code Generation Example

```go
package main

import (
    "fmt"
    "github.com/conneroisu/twerge"
    "os"
)

func main() {
    // Generate some classes
    twerge.Generate("flex items-center p-4")
    twerge.Generate("text-xl font-bold text-gray-800")
    twerge.Generate("bg-blue-500 text-white rounded")

    // Generate Go code
    code := twerge.GenerateClassMapCode()
    fmt.Println("Generated Code Preview:")
    fmt.Println("--------------------")
    fmt.Println(code)
    fmt.Println("--------------------")

    // Write to file
    err := twerge.WriteClassMapFile("tailwind_classes_gen.go")
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    fmt.Println("Go code written to tailwind_classes_gen.go")
}
```

## Configuration Example

```go
package main

import (
    "fmt"
    "github.com/conneroisu/twerge"
)

func main() {
    // Configure cache
    twerge.ConfigureCache(1000)
    fmt.Println("Cache configured with size 1000")

    // Set class prefix
    twerge.SetClassPrefix("app-")
    fmt.Println("Class prefix set to 'app-'")

    // Generate a class with the new prefix
    className := twerge.Generate("flex items-center p-4")
    fmt.Println("Generated class name:", className) // Should start with "app-"

    // Disable cache
    twerge.DisableCache()
    fmt.Println("Cache disabled")

    // Check if cache is enabled
    isEnabled := twerge.IsCacheEnabled()
    fmt.Println("Is cache enabled?", isEnabled) // Should be false
}
```
