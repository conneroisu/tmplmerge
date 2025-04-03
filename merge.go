package twerge

import (
	"fmt"
	"slices"
	"strings"
	"sync"
)

const (
	// HashLimit is the length of the generated class name without the prefix "tw-".
	HashLimit = 7
)

// ClassMapStr is a map of class strings to their generated class names
// This variable can be populated by code generation or manually
// It is protected by mapMutex for concurrent access
var (
	// Merge is the default template merger
	// It takes a space-delimited string of TailwindCSS classes and returns a merged string
	// It also adds the merged class to the ClassMapStr when used
	// It will quickly return the generated class name from ClassMapStr if available
	Merge = createTwMerge(nil, nil)

	ClassMapStr = make(map[string]string)

	// GenClassMergeStr is a map of merged class strings to their generated class names
	// This variable can be populated by code generation or manually
	// It is protected by mapMutex for concurrent access
	GenClassMergeStr = make(map[string]string)

	// mapMutex protects ClassMapStr for concurrent access
	mapMutex sync.RWMutex

	classID int
)

// twMergeFn is the type of the template merger.
type twMergeFn func(classes string) string

// splitModifiersFn is the type of the function used to split modifiers
type splitModifiersFn = func(string) (
	baseClass string,
	modifiers []string,
	hasImportant bool,
	maybePostfixModPosition int,
)

// createTwMerge creates a new template merger
func createTwMerge(
	config *config,
	cache icache,
) twMergeFn {
	var (
		fnToCall        twMergeFn
		splitModifiers  splitModifiersFn
		getClassGroupID getClassGroupIDFn
		mergeClassList  func(classList string) string
	)

	merger := func(classes string) string {
		classList := strings.TrimSpace(classes)
		if classList == "" {
			return ""
		}

		// Check if we've seen this class list before in the cache
		cached := cache.Get(classList)
		if cached != "" {
			return cached
		}

		// Merge the classes
		merged := mergeClassList(classList)
		cache.Set(classList, merged)

		// Add to ClassMapStr for lookup by other functions
		if classList != merged {
			mapMutex.Lock()
			className := fmt.Sprintf("tw-%d", classID)
			ClassMapStr[classList] = className
			GenClassMergeStr[className] = merged
			classID++
			mapMutex.Unlock()
		}

		return merged
	}

	init := func(classes string) string {
		if config == nil {
			config = defaultConfig
		}
		if cache == nil {
			cache = newCache(config.MaxCacheSize)
		}

		splitModifiers = makeSplitModifiers(config)

		getClassGroupID = makeGetClassGroupID(config)

		mergeClassList = makeMergeClassList(config, splitModifiers, getClassGroupID)

		fnToCall = merger
		return fnToCall(classes)
	}

	fnToCall = init
	return func(classes string) string {
		return fnToCall(classes)
	}
}

// makeMergeClassList creates a function that merges a class list
func makeMergeClassList(
	conf *config,
	splitModifiers splitModifiersFn,
	getClassGroupID getClassGroupIDFn,
) func(classList string) string {
	return func(classList string) string {
		classes := strings.Split(strings.TrimSpace(classList), " ")
		unqClasses := make(map[string]string, len(classes))
		resultClassList := ""

		for _, class := range classes {
			baseClass, modifiers, hasImportant, postFixMod := splitModifiers(class)

			// there is a postfix modifier -> text-lg/8
			if postFixMod != -1 {
				baseClass = baseClass[:postFixMod]
			}
			isTwClass, groupID := getClassGroupID(baseClass)
			if !isTwClass {
				resultClassList += class + " "
				continue
			}
			// we have to sort the modifiers bc hover:focus:bg-red-500 == focus:hover:bg-red-500
			modifiers = sortModifiers(modifiers)
			if hasImportant {
				modifiers = append(modifiers, "!")
			}
			unqClasses[groupID+strings.Join(modifiers, string(conf.ModifierSeparator))] = class

			conflicts := conf.ConflictingClassGroups[groupID]
			if conflicts == nil {
				continue
			}
			for _, conflict := range conflicts {
				// erase the conflicts with the same modifiers
				unqClasses[conflict+strings.Join(modifiers, string(conf.ModifierSeparator))] = ""
			}
		}

		for _, class := range unqClasses {
			if class == "" {
				continue
			}
			resultClassList += class + " "
		}
		return strings.TrimSpace(resultClassList)
	}

}

// sortModifiers Sorts modifiers according to following schema:
// - Predefined modifiers are sorted alphabetically
// - When an arbitrary variant appears, it must be preserved which modifiers are before and after it
func sortModifiers(modifiers []string) []string {
	if len(modifiers) < 2 {
		return modifiers
	}

	unsortedModifiers := []string{}
	sorted := make([]string, len(modifiers))

	for _, modifier := range modifiers {
		isArbitraryVariant := modifier[0] == '['
		if isArbitraryVariant {
			slices.Sort(unsortedModifiers)
			sorted = append(sorted, unsortedModifiers...)
			sorted = append(sorted, modifier)
			unsortedModifiers = []string{}
			continue
		}
		unsortedModifiers = append(unsortedModifiers, modifier)
	}

	slices.Sort(unsortedModifiers)
	sorted = append(sorted, unsortedModifiers...)

	return sorted
}

// makeSplitModifiers creates a function that splits modifiers
func makeSplitModifiers(conf *config) splitModifiersFn {
	separator := conf.ModifierSeparator

	return func(className string) (string, []string, bool, int) {
		modifiers := []string{}
		modifierStart := 0
		bracketDepth := 0
		// used for bg-red-500/50 (50% opacity)
		maybePostfixModPosition := -1

		for i := range len(className) {
			char := rune(className[i])

			if char == '[' {
				bracketDepth++
				continue
			}
			if char == ']' {
				bracketDepth--
				continue
			}

			if bracketDepth == 0 {
				if char == separator {
					modifiers = append(modifiers, className[modifierStart:i])
					modifierStart = i + 1
					continue
				}

				if char == conf.PostfixModifier {
					maybePostfixModPosition = i
				}
			}
		}

		baseClassWithImportant := className[modifierStart:]
		hasImportant := baseClassWithImportant[0] == byte(conf.ImportantModifier)

		var baseClass string
		if hasImportant {
			baseClass = baseClassWithImportant[1:]
		} else {
			baseClass = baseClassWithImportant
		}

		// fix case where there is modifier & maybePostfix which causes maybePostfix to be beyond size of baseClass!
		if maybePostfixModPosition != -1 && maybePostfixModPosition > modifierStart {
			maybePostfixModPosition -= modifierStart
		} else {
			maybePostfixModPosition = -1
		}

		return baseClass, modifiers, hasImportant, maybePostfixModPosition

	}
}

// // sanitizeCSSClassName replaces disallowed characters in CSS class names with a replacement string.
// // CSS class names can contain letters (a-z, A-Z), digits (0-9), hyphens (-), and underscores (_),
// // but cannot start with a digit, hyphen, or underscore.
// func sanitizeCSSClassName(className string, replacement string) string {
// 	if className == "" {
// 		return ""
// 	}
//
// 	// Handle the first character separately to enforce it must be a letter
// 	var result strings.Builder
// 	firstChar := rune(className[0])
// 	if unicode.IsLetter(firstChar) {
// 		result.WriteRune(firstChar)
// 	} else {
// 		result.WriteString(replacement)
// 	}
//
// 	// Process the remaining characters
// 	for _, char := range className[1:] {
// 		if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '-' || char == '_' {
// 			result.WriteRune(char)
// 		} else {
// 			result.WriteString(replacement)
// 		}
// 	}
//
// 	return result.String()
// }

// LintReport represents a report of duplicate merged class values
type LintReport struct {
	// MergedValue is the merged class string that has duplicates
	MergedValue string
	// OriginalClasses is a list of original class combinations that merge to the same value
	OriginalClasses []string
}

// Lint checks for multiple different class combinations that merge to the same final value
// Returns a slice of LintReport structures identifying duplicates
func Lint() []LintReport {
	// Create a map of merged class strings to original class strings
	mergedToOriginal := make(map[string][]string)
	var reports []LintReport

	// First, add all the current class mappings to the internal maps
	// This ensures we're not relying solely on cached maps from previous operations
	mapMutex.RLock()
	for original, generated := range ClassMapStr {
		merged, exists := GenClassMergeStr[generated]
		if !exists {
			continue // Skip if no merged value exists
		}

		mergedToOriginal[merged] = append(mergedToOriginal[merged], "'"+original+"'")
	}
	mapMutex.RUnlock()

	// Look for cases where different original classes result in the same merged output
	for merged, originals := range mergedToOriginal {
		if len(originals) > 1 {
			reports = append(reports, LintReport{
				MergedValue:     merged,
				OriginalClasses: originals,
			})
		}
	}

	return reports
}

// LintString returns a formatted string representation of the lint results
// The output format is designed to be easily readable, showing each merged value
// and the original class combinations that produce it
func LintString() string {
	reports := Lint()
	if len(reports) == 0 {
		return "No duplicate class combinations found."
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d cases where multiple class combinations merge to the same value:\n\n", len(reports)))

	for i, report := range reports {
		result.WriteString(fmt.Sprintf("%d. Merged value: \"%s\"\n", i+1, report.MergedValue))
		result.WriteString("   Original combinations:\n")

		for j, original := range report.OriginalClasses {
			result.WriteString(fmt.Sprintf("   %d.%d. \"%s\"\n", i+1, j+1, original))
		}

		result.WriteString("\n")
	}

	return result.String()
}
