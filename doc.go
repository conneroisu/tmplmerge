// Package twerge provides a tailwind merger for go-templ with class generation and runtime static hashmap.
//
// It performs four key functions:
// 1. Merges TailwindCSS classes intelligently (resolving conflicts)
// 2. Generates short unique CSS class names from the merged classes
// 3. Creates a mapping from original class strings to generated class names
// 4. Provides a runtime static hashmap for direct class name lookup
//
// Basic Usage:
//
//	import "github.com/conneroisu/twerge"
//
//	// Merge TailwindCSS classes from a space-delimited string
//	merged := twerge.Merge("text-red-500 bg-blue-500 text-blue-700")
//	// Returns: "bg-blue-500 text-blue-700"
//
//	// Generate a short unique class name
//	className := twerge.Generate("text-red-500 bg-blue-500")
//	// Returns something like: "tw-Ab3F5g7"
//
// Runtime Static Map Usage:
//
//	// Pre-register common classes
//	twerge.RegisterClasses(map[string]string{
//	  "bg-blue-500 text-white": "tw-btn-blue",
//	  "bg-red-500 text-white": "tw-btn-red",
//	})
//
//	// Generate a class name at runtime
//	className := twerge.RuntimeGenerate("p-4 m-2")
//	// Returns a deterministic class name, stored in the runtime map
//
//	// Generate CSS for all registered classes
//	css := twerge.GetRuntimeClassHTML()
//	// Returns CSS like: ".tw-btn-blue { @apply bg-blue-500 text-white; }"
//
// For templ users:
//
//	<div class={ twerge.Merge("bg-blue-500 p-4 bg-red-500") }>
//	  Using merged classes directly
//	</div>
//
//	<div class={ twerge.RuntimeGenerate("bg-blue-500 p-4") }>
//	  Using runtime generated class name
//	</div>
//
//	<style>
//	  @unsafe {
//	    twerge.GetRuntimeClassHTML()
//	  }
//	</style>
package twerge
