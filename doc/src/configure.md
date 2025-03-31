# Configuration

Twerge can be configured in several ways to customize its behavior. The library doesn't require a configuration file by default, but you can configure its behavior programmatically.

## Cache Configuration

Twerge uses an LRU cache for performance. You can configure its size:

```go
import "github.com/conneroisu/twerge"

func main() {
    // Configure cache size
    twerge.ConfigureCache(1000) // Set cache size to 1000 entries

    // Or disable caching entirely
    twerge.DisableCache()

    // Check if caching is enabled
    isEnabled := twerge.IsCacheEnabled()
}
```

## Class Generation Configuration

You can customize how class names are generated:

```go
import "github.com/conneroisu/twerge"

func main() {
    // Configure prefix for generated class names
    twerge.SetClassPrefix("tw-") // Generated classes will start with "tw-"

    // Configure custom hash algorithm (default is FNV-1a)
    twerge.SetHashFunction(twerge.HashFunctionSHA1)

    // Configure hash length
    twerge.SetHashLength(6) // Default is 8 characters
}
```

## Runtime Configuration

For the runtime static hashmap, you can configure pre-registered classes:

```go
import "github.com/conneroisu/twerge"

func main() {
    // Initialize with common Tailwind class combinations
    twerge.InitWithCommonClasses()

    // Or register your own set of classes with custom names
    customClasses := map[string]string{
        "flex items-center justify-between": "tw-header",
        "text-sm font-medium text-gray-500": "tw-label",
    }
    twerge.RegisterClasses(customClasses)
}
```

## CSS Integration Configuration

For CSS integration, you can configure the markers used in CSS files:

```go
import "github.com/conneroisu/twerge"

func main() {
    // Configure custom CSS markers
    twerge.SetCSSMarkers("/* TWERGE-START */", "/* TWERGE-END */")
}
```
