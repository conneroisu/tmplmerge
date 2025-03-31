package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ClassOccurrence stores information about where a class string occurs
type ClassOccurrence struct {
	FilePath     string
	LineNum      int
	FullElement  string
	OriginalText string
}

// Result stores the final results of the analysis
type Result struct {
	NormalizedClasses string
	ClassesCount      int
	Occurrences       []ClassOccurrence
	HasDifferentOrder bool
}

var (
	classPattern = regexp.MustCompile(`class\s*=\s*["']([^"']+)["']`)
)

// normalizeClassString normalizes a class string by sorting its components
func normalizeClassString(classStr string) string {
	classes := strings.Fields(classStr)
	sort.Strings(classes)
	return strings.Join(classes, " ")
}

// extractClassesFromFile extracts all class strings from a file with their locations
func extractClassesFromFile(filePath string, classOccurrences map[string][]ClassOccurrence) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	contentStr := string(content)

	// For line number tracking
	scanner := bufio.NewScanner(strings.NewReader(contentStr))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Find all class attributes
	matches := classPattern.FindAllStringSubmatchIndex(contentStr, -1)
	for _, match := range matches {
		// Extract the class attribute value (between quotes)
		originalClassStr := strings.TrimSpace(contentStr[match[2]:match[3]])
		if originalClassStr == "" {
			continue
		}

		// Normalize the class string
		normalizedClassStr := normalizeClassString(originalClassStr)

		// Find line number for this match
		position := match[0]
		lineNum := 1
		linePos := 0
		for i, line := range lines {
			if linePos+len(line)+1 > position {
				lineNum = i + 1
				break
			}
			linePos += len(line) + 1
		}

		// Get context for the element
		lineStartPos := strings.LastIndex(contentStr[:match[0]], "\n")
		if lineStartPos == -1 {
			lineStartPos = 0
		} else {
			lineStartPos++
		}

		lineEndPos := strings.Index(contentStr[match[0]:], "\n")
		if lineEndPos == -1 {
			lineEndPos = len(contentStr) - match[0]
		} else {
			lineEndPos += match[0]
		}

		elementContext := strings.TrimSpace(contentStr[lineStartPos:lineEndPos])
		if len(elementContext) > 100 {
			elementContext = elementContext[:97] + "..."
		}

		// Add the occurrence
		classOccurrences[normalizedClassStr] = append(
			classOccurrences[normalizedClassStr],
			ClassOccurrence{
				FilePath:     filePath,
				LineNum:      lineNum,
				FullElement:  elementContext,
				OriginalText: originalClassStr,
			},
		)
	}

	return nil
}

// processFiles processes all files in the patterns and finds duplicate class strings
func processFiles(patterns []string, minLength, minOccurrences, minClassCount int) ([]Result, error) {
	allClassOccurrences := make(map[string][]ClassOccurrence)

	// Process each pattern
	for _, pattern := range patterns {
		// Expand glob pattern
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("error with pattern %s: %w", pattern, err)
		}

		if len(matches) == 0 {
			fmt.Printf("%s Warning: No files match pattern '%s'\n", ("⚠"), pattern)
			continue
		}

		// Process each matching file
		for _, filePath := range matches {
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				fmt.Printf("%s Error accessing %s: %v\n", ("✗"), filePath, err)
				continue
			}

			if fileInfo.IsDir() {
				continue
			}

			if err := extractClassesFromFile(filePath, allClassOccurrences); err != nil {
				fmt.Printf("%s %v\n", ("✗"), err)
			}
		}
	}

	// Filter and prepare results
	var results []Result
	for normalizedClassStr, occurrences := range allClassOccurrences {
		if len(occurrences) < minOccurrences {
			continue
		}

		classes := strings.Fields(normalizedClassStr)
		if len(classes) < minClassCount || len(normalizedClassStr) < minLength {
			continue
		}

		// Check if there are different orderings
		orderingsMap := make(map[string]bool)
		for _, occ := range occurrences {
			orderingsMap[occ.OriginalText] = true
		}
		hasDifferentOrder := len(orderingsMap) > 1

		results = append(results, Result{
			NormalizedClasses: normalizedClassStr,
			ClassesCount:      len(classes),
			Occurrences:       occurrences,
			HasDifferentOrder: hasDifferentOrder,
		})
	}

	// Sort by potential space savings (length * occurrences)
	sort.Slice(results, func(i, j int) bool {
		savingsI := len(results[i].NormalizedClasses) * len(results[i].Occurrences)
		savingsJ := len(results[j].NormalizedClasses) * len(results[j].Occurrences)
		return savingsI > savingsJ
	})

	return results, nil
}

// printResults prints the results in a readable format with colors
func printResults(results []Result) {
	if len(results) == 0 {
		fmt.Println(("✓ No duplicate class strings found!"))
		return
	}

	totalDupes := len(results)
	totalOccurrences := 0
	totalChars := 0

	for _, result := range results {
		totalOccurrences += len(result.Occurrences)

		// Calculate character savings
		originalTotal := 0
		for _, occ := range result.Occurrences {
			originalTotal += len(occ.OriginalText)
		}
		optimizedTotal := len(result.Occurrences[0].OriginalText)
		totalChars += originalTotal - optimizedTotal
	}

	fmt.Printf("\n%s Found %d duplicate class sets across %d occurrences\n",
		("🔍"), totalDupes, totalOccurrences)
	fmt.Printf("%s Potential savings: approximately %d characters\n\n",
		("💾"), totalChars)

	for _, result := range results {
		// Calculate savings for this group
		originalTotal := 0
		for _, occ := range result.Occurrences {
			originalTotal += len(occ.OriginalText)
		}
		optimizedTotal := len(result.Occurrences[0].OriginalText)
		charSavings := originalTotal - optimizedTotal

		fmt.Printf("%s Duplicate Class Set (%d classes, %d occurrences, ~%d chars savings):\n",
			("🔄"), result.ClassesCount, len(result.Occurrences), charSavings)
		fmt.Printf("%s\n", (result.NormalizedClasses))

		if result.HasDifferentOrder {
			fmt.Printf("\n%s Different orderings found:\n", ("⚠"))

			// Get unique orderings
			orderings := make(map[string]bool)
			for _, occ := range result.Occurrences {
				orderings[occ.OriginalText] = true
			}

			i := 1
			for ordering := range orderings {
				fmt.Printf("  %d. %s\n", i, ordering)
				i++
			}
		}

		fmt.Printf("\n%s Occurrences:\n", "📍")

		for _, occ := range result.Occurrences {
			orderIndicator := ""
			if result.HasDifferentOrder {
				orderIndicator = fmt.Sprintf(" [order: %s]", occ.OriginalText)
			}

			fmt.Printf("  %s%s\n", (fmt.Sprintf("%s:%d", occ.FilePath, occ.LineNum)), orderIndicator)
			fmt.Printf("  %s\n\n", occ.FullElement)
		}

		fmt.Println(strings.Repeat("-", 80))
	}
}

func main() {
	// Simple command-line argument parsing
	if len(os.Args) < 2 {
		fmt.Println("Usage: class-duplication-finder [--min-length N] [--min-occurrences N] [--min-class-count N] [--no-color] <glob_patterns...>")
		os.Exit(1)
	}

	// Parse arguments
	var patterns []string
	minLength := 20
	minOccurrences := 2
	minClassCount := 1

	i := 1
	for i < len(os.Args) {
		arg := os.Args[i]

		switch arg {
		case "--min-length":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &minLength)
				i += 2
			} else {
				i++
			}
		case "--min-occurrences":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &minOccurrences)
				i += 2
			} else {
				i++
			}
		case "--min-class-count":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &minClassCount)
				i += 2
			} else {
				i++
			}
		case "--no-color":
			i++
		default:
			patterns = append(patterns, arg)
			i++
		}
	}

	if len(patterns) == 0 {
		fmt.Println(("Error: No glob patterns provided"))
		os.Exit(1)
	}

	// Process files and print results
	results, err := processFiles(patterns, minLength, minOccurrences, minClassCount)
	if err != nil {
		fmt.Printf("%s Error: %v\n", ("✗"), err)
		os.Exit(1)
	}

	printResults(results)

	if len(results) > 0 {
		os.Exit(1) // Return error code if duplicates found
	}
}
