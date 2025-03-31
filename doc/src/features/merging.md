# Merging Tailwind CSS Classes

One of the core features of Twerge is the ability to intelligently merge Tailwind CSS classes, resolving conflicts according to Tailwind's specificity rules.

## The Problem

When working with Tailwind CSS, you often need to combine multiple sets of classes, which can lead to conflicts. For example:

```go
// These classes conflict (two text colors)
"text-red-500 text-blue-500"
```

In this case, you want the last class to win (`text-blue-500`), following Tailwind's specificity rules.

## How Twerge Solves It

The `Merge` function in Twerge intelligently combines Tailwind classes, resolving conflicts correctly:

```go
import "github.com/conneroisu/twerge"

// Example usage
mergedClasses := twerge.Merge("text-red-500 bg-blue-300 text-xl")

// Result: "bg-blue-300 text-xl text-red-500"
// Classes are resolved and ordered by type
```

## Class Resolution Rules

Twerge follows Tailwind's conflict resolution rules:

1. **Last Declaration Wins** - For conflicting classes of the same type, the last one in the string takes precedence
2. **Type Preservation** - Non-conflicting classes are preserved
3. **Order Optimization** - The resulting class string is optimized for readability and consistency

## Supported Class Categories

Twerge understands and correctly handles conflicts between these Tailwind categories:

- Layout (display, position, etc.)
- Flexbox & Grid
- Spacing (margin, padding)
- Sizing (width, height)
- Typography (font, text)
- Backgrounds
- Borders
- Effects
- Filters
- Tables
- Transitions & Animation
- Transforms
- Interactivity
- SVG
- Accessibility

## Advanced Merging

Twerge can handle complex combinations, including:

```go
// Merging multiple complex class sets
classes1 := "flex items-center space-x-4 text-sm"
classes2 := "grid text-lg font-bold"
merged := twerge.Merge(classes1 + " " + classes2)

// Result includes non-conflicting classes and resolves conflicts
// "items-center space-x-4 grid text-lg font-bold"
```

## Performance Optimization

Twerge uses an LRU cache for frequently used class combinations:

```go
// Subsequent calls with the same input use the cache
result1 := twerge.Merge("p-4 m-2 p-8")  // Computed
result2 := twerge.Merge("p-4 m-2 p-8")  // Retrieved from cache
```

## Integration Examples

In Go-templ templates, you can use it like this:

```go
// In a templ file
<div class={ twerge.Merge("bg-blue-500 p-4 bg-red-500") }>
  This will have a red background with padding
</div>
```

## Related Functions

- `Merge(classes string) string` - Merges Tailwind classes
- `ConfigureCache(size int)` - Configures the cache size for merging operations
- `DisableCache()` - Disables caching for merging operations
