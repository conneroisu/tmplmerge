# Introduction to Twerge

Twerge is a Go library designed to enhance your experience working with Tailwind CSS in Go applications. The name "Twerge" comes from "Tailwind + Merge".

## What is Twerge?

Twerge is a comprehensive Go library that performs four key functions for Tailwind CSS integration:

1. **Intelligent Class Merging** - Resolves conflicts between Tailwind CSS classes according to their specificity rules
2. **Class Name Generation** - Creates short, unique CSS class names based on hashes of the merged classes
3. **Class Mapping Management** - Maintains mappings between original class strings and generated class names, with code generation capabilities
4. **Runtime Static Hashmap** - Provides a fast runtime lookup for direct class name resolution without a generation step

## Why Use Twerge?

If you're developing Go-based web applications with Tailwind CSS, Twerge offers significant advantages:

- **Smaller HTML output** - By merging conflicting classes and generating short class names
- **Better performance** - Through intelligent caching and efficient lookups
- **Build-time optimization** - Via code generation capabilities
- **Runtime flexibility** - Through the runtime static hashmap for dynamic class handling
- **Simplified workflow** - By integrating seamlessly with Go templates, particularly [templ](https://github.com/a-h/templ)

## Key Features

- **Intelligent class merging** - Resolves conflicts according to Tailwind CSS specificity rules
- **Short class name generation** - Creates compact, unique class names for reduced HTML size
- **Runtime class management** - Provides a fast lookup system for dynamic applications
- **Code generation** - Produces optimized Go code for class mappings
- **CSS integration** - Works with Tailwind CLI and CSS build pipelines
- **Flexible configuration** - Customizable caching, hash algorithms, and more
- **Nix integration** - Reproducible development environment

## Target Use Cases

Twerge is particularly well-suited for:

- Go web applications using Tailwind CSS
- Projects using the [templ](https://github.com/a-h/templ) templating language
- Applications requiring build-time CSS optimization
- Static site generators with Tailwind CSS integration
- Dynamic web applications needing runtime class management

## Next Steps

To get started with Twerge, check out:

- [Installation](./install.md) - How to install Twerge in your Go project
- [Configuration](./configure.md) - How to configure Twerge for your needs
- [Merging Classes](./features/merging.md) - Learn about the core class merging functionality
