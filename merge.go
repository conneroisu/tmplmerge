package twerge

import (
	"regexp"
	"slices"
	"strings"
)

var (
	// Merge is the default template merger
	Merge        = CreateTwMerge(nil, nil)
	splitPattern = regexp.MustCompile(SplitClassesRegex)
)

// SplitClassesRegex is the regex used to split classes
const SplitClassesRegex = `\s+`

// TwMergeFn is the type of the template merger.
type TwMergeFn func(args ...string) string

// SplitModifiersFn is the type of the function used to split modifiers
type SplitModifiersFn = func(string) (
	baseClass string,
	modifiers []string,
	hasImportant bool,
	maybePostfixModPosition int,
)

// CreateTwMerge creates a new template merger
func CreateTwMerge(
	config *Config,
	cache ICache,
) TwMergeFn {
	var (
		fnToCall        TwMergeFn
		splitModifiers  SplitModifiersFn
		getClassGroupID GetClassGroupIDFn
		mergeClassList  func(classList string) string
	)

	merger := func(args ...string) string {
		classList := strings.TrimSpace(strings.Join(args, " "))
		if classList == "" {
			return ""
		}
		cached := cache.Get(classList)
		if cached != "" {
			return cached
		}
		// check if in cache
		merged := mergeClassList(classList)
		cache.Set(classList, merged)
		return merged
	}

	init := func(args ...string) string {
		if config == nil {
			config = DefaultConfig
		}
		if cache == nil {
			cache = Make(config.MaxCacheSize)
		}

		splitModifiers = MakeSplitModifiers(config)

		getClassGroupID = MakeGetClassGroupID(config)

		mergeClassList = MakeMergeClassList(config, splitModifiers, getClassGroupID)

		fnToCall = merger
		return fnToCall(args...)
	}

	fnToCall = init
	return func(args ...string) string {
		return fnToCall(args...)
	}
}

// MakeMergeClassList creates a function that merges a class list
func MakeMergeClassList(
	conf *Config,
	splitModifiers SplitModifiersFn,
	getClassGroupID GetClassGroupIDFn,
) func(classList string) string {
	return func(classList string) string {
		classes := splitPattern.Split(strings.TrimSpace(classList), -1)
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
			modifiers = SortModifiers(modifiers)
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

// SortModifiers Sorts modifiers according to following schema:
// - Predefined modifiers are sorted alphabetically
// - When an arbitrary variant appears, it must be preserved which modifiers are before and after it
func SortModifiers(modifiers []string) []string {
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

// MakeSplitModifiers creates a function that splits modifiers
func MakeSplitModifiers(conf *Config) SplitModifiersFn {
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
