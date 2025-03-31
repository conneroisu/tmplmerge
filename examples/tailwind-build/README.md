# Twerge with Tailwind CLI Build Process

This example demonstrates how to use Twerge in a Tailwind CSS build process.

## Overview

This example shows how to:

1. Register Tailwind utility classes with Twerge
2. Generate CSS content with `@apply` directives
3. Insert the CSS into an input file between specific markers
4. Run the Tailwind CLI to process the input file 
5. Generate optimized CSS output for production

## How It Works

The example creates a complete workflow:

1. **Class Registration**: Define reusable components using Tailwind utility classes
2. **CSS Generation**: Insert `@apply` rules between `<!-- twerge:begin -->` and `<!-- twerge:end -->` markers
3. **Tailwind Processing**: Run the Tailwind CLI to process the `@apply` directives
4. **HTML Integration**: Use the generated class names in HTML

## Utility Functions

The example uses the following Twerge utility functions:

- `RegisterClasses`: Register Tailwind utility classes with custom names
- `GenerateInputCSSForTailwind`: Create an input CSS file for the Tailwind CLI
- `GetRuntimeClassHTML`: Generate CSS with `@apply` directives
- `ExportCSS`: Export CSS to a file between markers

## Running the Example

To run this example:

1. Make sure you have Node.js and npm installed (for Tailwind CLI)
2. Navigate to the example directory:
   ```
   cd examples/tailwind-build
   ```
3. Run the example:
   ```
   go run main.go
   ```
4. Open the generated HTML file in your browser:
   ```
   open dist/index.html
   ```

## Production Use Cases

In a real-world scenario, you would:

1. Define your utility class compositions during development
2. Generate the optimized CSS as part of your build process
3. Deploy the generated CSS and HTML with the short class names

This approach gives you the benefits of Tailwind's utility-first approach during development, while producing optimized CSS for production.

## Files Generated

- `src/input.css`: Input CSS file with Twerge-generated `@apply` directives
- `dist/css/styles.css`: Processed CSS file from Tailwind CLI
- `dist/index.html`: Sample HTML using the generated class names

## Additional Resources

For more information, see:

- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Twerge Documentation](https://github.com/conneroisu/twerge)