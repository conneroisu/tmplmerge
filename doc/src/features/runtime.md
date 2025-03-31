# Runtime Class Management

Twerge provides a runtime static hashmap for direct class name lookups, allowing for dynamic class handling without a build or generation step.

## The Problem

In dynamic applications, you may not know all the class combinations in advance, making it challenging to use build-time generation. You need a way to:

1. Handle classes that are determined at runtime
2. Maintain performance with dynamic class generation
3. Still benefit from class shortening and merging

## How Twerge Solves It

The runtime functions in Twerge provide a static hashmap that manages class mappings at runtime:

```go
import "github.com/conneroisu/twerge"

// Generate a class name at runtime
shortClassName := twerge.RuntimeGenerate("flex items-center justify-between p-4")

// Result: "tw-a1b2c3d4" (example - actual output will vary)
// The mapping is automatically stored in the runtime map
```

## Key Runtime Functions

### Generate and Store

```go
// Generate a class name and store it in the runtime map
className := twerge.RuntimeGenerate("flex p-4 text-lg")

// This stores the mapping internally for later use
```

### Pre-register Classes

```go
// Pre-register classes with custom names
customClasses := map[string]string{
    "flex items-center justify-between": "tw-header",
    "text-sm font-medium text-gray-500": "tw-label",
}
twerge.RegisterClasses(customClasses)

// Now you can use these custom names
className := twerge.RuntimeGenerate("flex items-center justify-between")
// Result: "tw-header"
```

### Common Classes

```go
// Initialize with common Tailwind class combinations
twerge.InitWithCommonClasses()

// This pre-populates the runtime map with frequently used combinations
```

### Access the Mapping

```go
// Get the current runtime mapping
mapping := twerge.GetRuntimeMapping()

// This returns a copy of the current runtime class map
// map[string]string{"flex p-4 text-lg": "tw-a1b2c3d4", ...}
```

### Generate HTML/CSS

```go
// Generate CSS/HTML for all registered class mappings
cssHTML := twerge.GetRuntimeClassHTML()

// This produces HTML like:
// <style>
// .tw-a1b2c3d4 { display: flex; padding: 1rem; font-size: 1.125rem; }
// </style>
```

## Integration Example

In a Go-templ application:

```go
// Initialize common classes at startup
func init() {
    twerge.InitWithCommonClasses()
}

// In your component template
templ Button(classes string) {
    <button class={ twerge.RuntimeGenerate(classes) }>
        { children... }
    </button>
}

// In your layout template
templ Layout() {
    <html>
        <head>
            <style>{ twerge.GetRuntimeClassHTML() }</style>
        </head>
        <body>
            { children... }
        </body>
    </html>
}
```

## Performance Considerations

- The runtime map is thread-safe
- Lookups are very fast (O(1) complexity)
- The map grows with unique class combinations
- For very large applications, consider using the build-time generation instead

## Combining with Other Approaches

You can combine the runtime approach with build-time generation:

1. Use build-time generation for known, common class combinations
2. Use runtime generation for dynamic, user-defined classes

```go
// Pre-register known classes
twerge.RegisterClasses(generatedClassMap)

// Use runtime generation for dynamic classes
className := twerge.RuntimeGenerate(dynamicClasses)
```

## Exporting Runtime Mappings

You can export the runtime mappings for persistence:

```go
// Export the runtime mapping to a CSS file
twerge.ExportCSS("runtime-styles.css")

// Save the mapping for later use
mapping := twerge.GetRuntimeMapping()
// Save to file, database, etc.
```
