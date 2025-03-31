# Installation

## Using Go (recommended)

```bash
go get github.com/conneroisu/twerge
```

## Using Git

You can also clone the repository directly:

```bash
git clone https://github.com/conneroisu/twerge.git
cd twerge
```

## Nix Integration

Twerge comes with Nix support for a reproducible development environment:

```bash
# Start a nix shell with all dependencies
nix-shell

# Or with nix flakes
nix develop
```

## Verifying Installation

To verify the installation, create a simple Go program:

```go
package main

import (
    "fmt"
    "github.com/conneroisu/twerge"
)

func main() {
    // Merge conflicting Tailwind classes
    merged := twerge.Merge("text-red-500 bg-blue-300 text-xl")
    fmt.Println(merged) // Should output "bg-blue-300 text-xl text-red-500" or similar order
}
```

Run the program:

```bash
go run main.go
```
