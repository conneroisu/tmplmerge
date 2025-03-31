# Web Server Integration

This example demonstrates how to integrate Twerge with a web server in Go, showing how to use it with HTML templates and dynamic content.

## Basic Web Server Example

This example uses the standard `net/http` package and shows how to integrate Twerge for class management:

```go
package main

import (
    "fmt"
    "github.com/conneroisu/twerge"
    "html/template"
    "net/http"
)

// Template data structure
type PageData struct {
    Title string
    Content string
    IsDarkMode bool
}

// HTML template with Twerge classes
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        {{.CSS}}
    </style>
</head>
<body class="{{.BodyClass}}">
    <header class="{{.HeaderClass}}">
        <h1 class="{{.TitleClass}}">{{.Title}}</h1>
    </header>
    <main class="{{.ContentClass}}">
        <p>{{.Content}}</p>
    </main>
    <footer class="{{.FooterClass}}">
        &copy; 2023 Twerge Example
    </footer>
</body>
</html>
`

func main() {
    // Initialize Twerge with common classes
    twerge.InitWithCommonClasses()

    // Register custom component classes
    twerge.RegisterClasses(map[string]string{
        "bg-white dark:bg-gray-800 min-h-screen": "tw-body",
        "flex justify-between items-center p-4 border-b": "tw-header",
        "text-2xl font-bold text-gray-800 dark:text-white": "tw-title",
        "p-6 max-w-4xl mx-auto": "tw-content",
        "mt-8 p-4 text-center text-gray-600 dark:text-gray-400 text-sm": "tw-footer",
    })

    // Define HTTP handler
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Determine if dark mode from query param
        isDarkMode := r.URL.Query().Get("darkMode") == "true"

        // Prepare template data
        data := PageData{
            Title: "Twerge Web Server Example",
            Content: "This example demonstrates how to use Twerge in a web server.",
            IsDarkMode: isDarkMode,
        }

        // Prepare template with Twerge-generated classes
        tmplData := struct {
            PageData
            CSS        template.HTML
            BodyClass  string
            HeaderClass string
            TitleClass  string
            ContentClass string
            FooterClass string
        }{
            PageData: data,
            CSS:        template.HTML(twerge.GetRuntimeClassHTML()),
            BodyClass:  twerge.RuntimeGenerate("bg-white dark:bg-gray-800 min-h-screen"),
            HeaderClass: twerge.RuntimeGenerate("flex justify-between items-center p-4 border-b"),
            TitleClass:  twerge.RuntimeGenerate("text-2xl font-bold text-gray-800 dark:text-white"),
            ContentClass: twerge.RuntimeGenerate("p-6 max-w-4xl mx-auto"),
            FooterClass: twerge.RuntimeGenerate("mt-8 p-4 text-center text-gray-600 dark:text-gray-400 text-sm"),
        }

        // Parse and execute template
        tmpl, err := template.New("page").Parse(htmlTemplate)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = tmpl.Execute(w, tmplData)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    })

    // Start server
    fmt.Println("Server starting at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

## Integration with Templ

This example shows how to integrate Twerge with the [templ](https://github.com/a-h/templ) templating language:

```go
// File: view.templ
package main

import "github.com/conneroisu/twerge"

// Layout component
templ Layout(title string) {
    <!DOCTYPE html>
    <html>
        <head>
            <title>{ title }</title>
            <style>
                { twerge.GetRuntimeClassHTML() }
            </style>
        </head>
        <body class={ twerge.RuntimeGenerate("bg-white dark:bg-gray-800 min-h-screen") }>
            { children... }
        </body>
    </html>
}

// Header component
templ Header(title string) {
    <header class={ twerge.RuntimeGenerate("flex justify-between items-center p-4 border-b") }>
        <h1 class={ twerge.RuntimeGenerate("text-2xl font-bold text-gray-800 dark:text-white") }>{ title }</h1>
        <nav class={ twerge.RuntimeGenerate("flex gap-4") }>
            <a href="/" class={ twerge.RuntimeGenerate("text-blue-600 hover:text-blue-800 dark:text-blue-400") }>Home</a>
            <a href="/about" class={ twerge.RuntimeGenerate("text-blue-600 hover:text-blue-800 dark:text-blue-400") }>About</a>
        </nav>
    </header>
}

// Content component
templ Content() {
    <main class={ twerge.RuntimeGenerate("p-6 max-w-4xl mx-auto") }>
        { children... }
    </main>
}

// Footer component
templ Footer() {
    <footer class={ twerge.RuntimeGenerate("mt-8 p-4 text-center text-gray-600 dark:text-gray-400 text-sm") }>
        &copy; 2023 Twerge Example
    </footer>
}

// Home page
templ HomePage() {
    @Layout("Twerge + Templ Example") {
        @Header("Twerge + Templ Example")
        @Content() {
            <div class={ twerge.RuntimeGenerate("prose dark:prose-invert") }>
                <p>This example demonstrates how to use Twerge with the Templ templating language.</p>
                <p>The classes are dynamically generated and the CSS is included in the page.</p>
            </div>
        }
        @Footer()
    }
}
```

```go
// File: main.go
package main

import (
    "fmt"
    "github.com/conneroisu/twerge"
    "net/http"
)

func main() {
    // Initialize Twerge
    twerge.InitWithCommonClasses()

    // Register routes
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Render the home page
        component := HomePage()
        component.Render(r.Context(), w)
    })

    // Start server
    fmt.Println("Server starting at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

## Advanced Example: API with Dynamic Classes

This example shows how to use Twerge in an API that dynamically generates classes based on data:

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/conneroisu/twerge"
    "net/http"
)

// Component data
type Component struct {
    ID          string `json:"id"`
    Type        string `json:"type"`
    Label       string `json:"label"`
    Description string `json:"description"`
    Color       string `json:"color"`
    Size        string `json:"size"`
}

// Generate component classes based on data
func generateComponentClasses(c Component) map[string]string {
    // Base classes for different component types
    baseClasses := map[string]string{
        "button": "inline-flex items-center justify-center rounded font-medium focus:outline-none focus:ring-2 focus:ring-offset-2",
        "card":   "rounded overflow-hidden shadow",
        "badge":  "inline-flex items-center rounded-full text-xs font-medium",
        "alert":  "p-4 rounded",
    }

    // Color classes
    colorClasses := map[string]string{
        "blue":  "bg-blue-500 text-white hover:bg-blue-600 focus:ring-blue-500",
        "green": "bg-green-500 text-white hover:bg-green-600 focus:ring-green-500",
        "red":   "bg-red-500 text-white hover:bg-red-600 focus:ring-red-500",
        "gray":  "bg-gray-200 text-gray-800 hover:bg-gray-300 focus:ring-gray-500",
    }

    // Size classes
    sizeClasses := map[string]string{
        "sm": "px-2.5 py-1.5 text-xs",
        "md": "px-4 py-2 text-sm",
        "lg": "px-6 py-3 text-base",
    }

    // Get base classes for the component type
    classes := baseClasses[c.Type]
    if classes == "" {
        classes = baseClasses["button"] // Default to button
    }

    // Add color classes
    if color := colorClasses[c.Color]; color != "" {
        classes += " " + color
    }

    // Add size classes
    if size := sizeClasses[c.Size]; size != "" {
        classes += " " + size
    }

    // Generate a short class name
    shortClassName := twerge.RuntimeGenerate(classes)

    return map[string]string{
        "classes": classes,
        "className": shortClassName,
    }
}

func main() {
    // Initialize Twerge
    twerge.InitWithCommonClasses()

    // API routes
    http.HandleFunc("/api/components", func(w http.ResponseWriter, r *http.Request) {
        // Example components
        components := []Component{
            {ID: "btn1", Type: "button", Label: "Save", Color: "blue", Size: "md"},
            {ID: "btn2", Type: "button", Label: "Cancel", Color: "gray", Size: "md"},
            {ID: "card1", Type: "card", Label: "User Profile", Color: "blue", Size: "lg"},
            {ID: "badge1", Type: "badge", Label: "New", Color: "red", Size: "sm"},
        }

        // Process each component
        result := make([]map[string]interface{}, 0, len(components))
        for _, comp := range components {
            classes := generateComponentClasses(comp)

            // Create response object
            compResponse := map[string]interface{}{
                "id":          comp.ID,
                "type":        comp.Type,
                "label":       comp.Label,
                "description": comp.Description,
                "classes":     classes["classes"],
                "className":   classes["className"],
            }

            result = append(result, compResponse)
        }

        // Add CSS to the response
        css := twerge.GetRuntimeClassHTML()

        // Create final response
        response := map[string]interface{}{
            "components": result,
            "css":        css,
        }

        // Return JSON response
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })

    // Start server
    fmt.Println("API server starting at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

See the [web-server example](https://github.com/conneroisu/twerge/tree/main/examples/web-server) in the GitHub repository for a complete working example.
