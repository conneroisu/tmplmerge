# Generating Short Class Names

Twerge can generate short, unique class names based on hashes of the merged Tailwind classes, allowing for smaller HTML output and improved performance.

## The Problem

Tailwind CSS is incredibly powerful, but it can lead to long class strings:

```html
<div
  class="flex flex-col items-center justify-between p-4 rounded-lg shadow-lg bg-white hover:bg-gray-50 dark:bg-gray-800 dark:hover:bg-gray-700 transition-all duration-300"
>
  <!-- Content -->
</div>
```

These long strings increase HTML size, reduce readability, and can impact performance.

## How Twerge Solves It

The `Generate` function creates short, unique class names based on the hash of the merged Tailwind classes:

```go
import "github.com/conneroisu/twerge"

// Long class string
classes := "flex flex-col items-center justify-between p-4 rounded-lg shadow-lg bg-white hover:bg-gray-50"

// Generate a short unique class name
shortClassName := twerge.It(classes)

// Result: "tw-a1b2c3d4" (example - actual output will vary)
```

## Benefits of Generated Class Names

- **Smaller HTML** - Dramatically reduces HTML file size
- **Better Caching** - Improves browser caching of HTML
- **Consistent Naming** - Same class string always generates the same short name
- **Collision-Free** - Hash-based generation ensures uniqueness
- **Automatic Conflict Resolution** - Classes are merged before hashing

## How It Works

1. First, Twerge merges the provided classes to resolve conflicts
2. The merged class string is hashed using a fast algorithm (default is FNV-1a)
3. The hash is converted to a short alphanumeric string
4. A prefix is added (default is "tw-")
5. The mapping from original classes to the generated name is stored

## Customizing Generation

You can customize how class names are generated:

```go
// Set a custom prefix for generated class names
twerge.SetClassPrefix("app-")

// Set a custom hash length (default is 8)
twerge.SetHashLength(6)

// Set a different hash algorithm
twerge.SetHashFunction(twerge.HashFunctionSHA1)
```

## Integration Example

In a Go-templ template:

```go
// Instead of this with long class strings
<div class="flex items-center justify-between p-4 bg-white dark:bg-gray-800">...</div>

// You can do this with a short generated class
<div class={ twerge.It("flex items-center justify-between p-4 bg-white dark:bg-gray-800") }>...</div>

// Which results in HTML like:
// <div class="tw-a1b2c3d4">...</div>
```

## Mapping and CSS Generation

When using generated class names, you'll need to generate the corresponding CSS. Twerge provides several ways to do this:

```go
// Get the current mapping of original classes to generated names
mapping := twerge.GetMapping()

// Export CSS with the mappings to a file
twerge.ExportCSS("styles.css")
```

These aspects are covered in more detail in the [CSS Integration](./css-integration.md) and [Class Mapping](./mapping.md) sections.

## Performance Considerations

- Generated class names are cached for performance
- The same input always produces the same output
- Short names reduce HTML size and parsing time
