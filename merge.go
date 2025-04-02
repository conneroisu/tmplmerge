package twerge

import (
	"crypto/md5"
	"encoding/base64"
	"slices"
	"strings"
	"sync"
)

// ClassMapStr is a map of class strings to their generated class names
// This variable can be populated by code generation or manually
// It is protected by mapMutex for concurrent access
var ClassMapStr = make(map[string]string)

// GenClassMergeStr is a map of class strings to their generated class names
// This variable can be populated by code generation or manually
// It is protected by mapMutex for concurrent access
var GenClassMergeStr = make(map[string]string)

// mapMutex protects ClassMapStr for concurrent access
var mapMutex sync.RWMutex

var (
	// Merge is the default template merger
	// It takes a space-delimited string of TailwindCSS classes and returns a merged string
	// It also adds the merged class to the ClassMapStr when used
	// It will quickly return the generated class name from ClassMapStr if available
	Merge = createTwMerge(nil, nil)
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
			// Add both the original and merged versions to ClassMapStr
			hash := md5.Sum([]byte(merged))
			className := "tw-" + base64.URLEncoding.EncodeToString(hash[:])[:7]

			mapMutex.Lock()
			ClassMapStr[classList] = className
			GenClassMergeStr[className] = merged
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
