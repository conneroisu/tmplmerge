# Advanced Tailwind CSS Integration with Twerge

This example demonstrates advanced integration between Twerge and Tailwind CSS, showcasing the full range of CSS export utilities.

## Features Demonstrated

1. **Multiple Export Formats**

   - Standard CSS with `@apply` directives
   - SCSS format for Sass workflows
   - LESS format for Less workflows
   - Minified CSS for production

2. **Customization Options**

   - Class name prefixing
   - Custom markers for CSS injection
   - Comment toggling
   - Minification settings

3. **Build System Integration**

   - PostCSS configuration generation
   - Tailwind CLI integration
   - Watch mode for development

4. **Export Utilities Showcased**
   - `ExportCSS`: Basic CSS export
   - `ExportCSSWithOptions`: Advanced options
   - `ExportOptimizedCSS`: Production-ready CSS
   - `GeneratePostCSSConfig`: PostCSS integration
   - `WatchAndExportCSS`: Development mode with live updates

## Running the Example

1. Install dependencies:

   ```
   npm install
   ```

2. Run the example:

   ```
   go run main.go
   ```

3. Open `dist/index.html` in your browser to see the results.

## CSS Output Files

The example generates multiple CSS output files to demonstrate different export options:

- `src/input.css`: Standard CSS with Twerge markers
- `src/prefix-input.css`: CSS with prefixed class names
- `src/minified-input.css`: Minified CSS without comments
- `src/scss/input.scss`: SCSS format export
- `src/less/input.less`: LESS format export
- `src/optimized-input.css`: Production-optimized CSS
- `src/custom-markers.css`: CSS with custom marker syntax
- `src/watched-input.css`: Automatically updated CSS file

## HTML Demo

The generated HTML file demonstrates the use of Twerge-generated class names in a realistic web page, including:

- Responsive layout components
- UI elements (buttons, cards, forms)
- Typography styles
- Navigation components

## Key Concepts

1. **CSS Export Options**: Control the format, minification, comments, and prefix of the exported CSS.

2. **Markers**: Designate where Twerge should insert generated CSS in your files.

3. **Watch Mode**: Automatically update CSS files when classes change, ideal for development.

4. **Build Integration**: Integrate with popular build tools like PostCSS and Tailwind CLI.

5. **Class Reuse**: Register commonly used Tailwind utility combinations under semantic names.
