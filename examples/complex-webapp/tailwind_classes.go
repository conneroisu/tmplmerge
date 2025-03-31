package main

import "github.com/conneroisu/twerge"

// InitTailwindClasses sets up the tailwind class mapping
func InitTailwindClasses() {
	// Register common layout classes
	twerge.RegisterClasses(map[string]string{
		// Layout & Container
		"container mx-auto px-4 sm:px-6 lg:px-8": "tw-container",
		"max-w-7xl mx-auto":                      "tw-container-lg",
		"max-w-5xl mx-auto":                      "tw-container-md",
		"max-w-3xl mx-auto":                      "tw-container-sm",

		// Flexbox layouts
		"flex items-center justify-between": "tw-flex-between",
		"flex items-center justify-center":  "tw-flex-center",
		"flex flex-col space-y-4":           "tw-flex-col",
		"flex items-center space-x-4":       "tw-flex-row",

		// Grid layouts
		"grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6": "tw-grid-responsive",
		"grid grid-cols-1 sm:grid-cols-2 gap-4":                "tw-grid-cards",

		// Spacing
		"py-2 px-4":         "tw-padding-btn",
		"p-4 sm:p-6":        "tw-padding-card",
		"py-8 px-4 sm:px-6": "tw-padding-section",
		"mt-4 mb-8":         "tw-margin-section",
		"mx-auto my-6":      "tw-margin-center",

		// Typography
		"text-xl font-bold text-gray-900":     "tw-heading-xl",
		"text-lg font-semibold text-gray-800": "tw-heading-lg",
		"text-md font-medium text-gray-800":   "tw-heading-md",
		"text-sm text-gray-600":               "tw-text-small",
		"text-xs text-gray-500":               "tw-text-xs",

		// Buttons
		"py-2 px-4 bg-indigo-600 hover:bg-indigo-700 text-white rounded-md shadow-sm": "tw-btn-primary",
		"py-2 px-4 bg-gray-200 hover:bg-gray-300 text-gray-900 rounded-md shadow-sm":  "tw-btn-secondary",
		"py-2 px-4 bg-red-600 hover:bg-red-700 text-white rounded-md shadow-sm":       "tw-btn-danger",
		"py-2 px-4 border border-gray-300 hover:bg-gray-50 rounded-md shadow-sm":      "tw-btn-outline",

		// Cards
		"bg-white rounded-lg shadow-md overflow-hidden":                       "tw-card",
		"bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow": "tw-card-hover",
		"bg-gray-50 rounded-lg p-4":                                           "tw-card-alt",

		// Forms
		"block w-full rounded-md border-gray-300 shadow-sm":                                                                             "tw-input",
		"block w-full text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md": "tw-input-focus",

		// Navigation
		"bg-white shadow": "tw-navbar",
		"flex items-center px-4 py-6 sm:px-6 md:justify-start md:space-x-10": "tw-navbar-container",

		// Page sections
		"py-12 bg-white":   "tw-section",
		"py-16 bg-gray-50": "tw-section-alt",

		// States
		"opacity-50 cursor-not-allowed": "tw-disabled",
		"animate-pulse":                 "tw-loading",

		// Utilities
		"sr-only":         "tw-sr-only",
		"hidden sm:block": "tw-hide-mobile",
		"block sm:hidden": "tw-hide-desktop",
	})
}
