# Twerge Complex Web Application Example

This example demonstrates how to use Twerge in a more complex web application setting, showcasing how to efficiently manage Tailwind CSS classes in a real-world project.

## Overview

The example application is an e-commerce store called "TechShop" with the following features:

- Responsive design using Tailwind CSS
- Multiple page templates with shared components
- Dynamic content rendering
- Interactive elements using HTMX
- Intelligent class name management with Twerge

## Components

The application includes:

1. **Layouts**:

   - Base layout with header and footer

2. **Pages**:

   - Home page with featured products and promotional sections
   - Products page with filtering capabilities
   - Shopping cart page with interactive elements

3. **Components**:
   - Header with navigation and user menu
   - Footer with multiple sections
   - Product card for displaying product details
   - Cart item for managing cart contents

## Twerge Implementation

This example demonstrates three different approaches to using Twerge:

1. **Direct Class Names**: Using predefined class names like `tw-container`, `tw-card`, etc.

   ```html
   <div class="tw-container">...</div>
   ```

2. **RuntimeGenerate**: Generating optimized class names at runtime

   ```html
   <div class={ twerge.RuntimeGenerate("flex items-center justify-between") }>...</div>
   ```

3. **Merged Classes**: Using the raw merged class strings
   ```html
   <div class={ twerge.Merge("flex items-center justify-between") }>...</div>
   ```

## Running the Example

To run this example:

1. Navigate to the example directory:

   ```
   cd examples/complex-webapp
   ```

2. Run the application:

   ```
   go run main.go tailwind_classes.go
   ```

3. Open your browser and visit:
   ```
   http://localhost:8080
   ```

## Code Generation

In a production environment, you would typically:

1. Generate a static CSS file from the runtime class map:

   ```go
   css := twerge.GetRuntimeClassHTML()
   os.WriteFile("public/css/tailwind-classes.css", []byte(css), 0644)
   ```

2. Or generate a Go file with predefined class mappings:
   ```go
   twerge.WriteClassMapFile("classmap_generated.go")
   ```

This approach allows for both development flexibility and production optimization.
