package tmplmerge

import (
	"regexp"
	"strings"
)

// GetClassGroupIdfn returns the class group id for a given class
type GetClassGroupIdfn func(string) (isTwClass bool, groupId string)

// MakeGetClassGroupID returns a GetClassGroupIdfn
func MakeGetClassGroupID(conf *Config) GetClassGroupIdfn {
	var getClassGroupIDRecursive func(classParts []string, i int, classMap *ClassPart) (isTwClass bool, groupId string)
	getClassGroupIDRecursive = func(classParts []string, i int, classMap *ClassPart) (isTwClass bool, groupId string) {
		if i >= len(classParts) {
			if classMap.ClassGroupID != "" {
				return true, classMap.ClassGroupID
			}

			return false, ""
		}

		if classMap.NextPart != nil {
			nextClassMap := classMap.NextPart[classParts[i]]
			isTw, id := getClassGroupIDRecursive(classParts, i+1, &nextClassMap)
			if isTw {
				return isTw, id
			}
		}

		if len(classMap.Validators) > 0 {
			remainingClass := strings.Join(classParts[i:], string(conf.ClassSeparator))

			for _, validator := range classMap.Validators {
				if validator.Fn(remainingClass) {
					return true, validator.ClassGroupID
				}
			}

		}
		return false, ""
	}

	var arbitraryPropertyRegex = regexp.MustCompile(`^\[(.+)\]$`)

	getGroupIDForArbitraryProperty := func(class string) (bool, string) {
		if arbitraryPropertyRegex.MatchString(class) {
			arbitraryPropertyClassName := arbitraryPropertyRegex.FindStringSubmatch(class)[1]
			property := arbitraryPropertyClassName[:strings.Index(arbitraryPropertyClassName, ":")]

			if property != "" {
				// I use two dots here because one dot is used as prefix for class groups in plugins
				return true, "arbitrary.." + property
			}
		}

		return false, ""
	}

	return func(baseClass string) (isTwClass bool, groupdId string) {

		classParts := strings.Split(baseClass, string(conf.ClassSeparator))
		// remove first element if empty for things like -px-4
		if len(classParts) > 0 && classParts[0] == "" {
			classParts = classParts[1:]
		}
		isTwClass, groupID := getClassGroupIDRecursive(classParts, 0, &conf.ClassGroups)
		if isTwClass {
			return isTwClass, groupID
		}

		return getGroupIDForArbitraryProperty(baseClass)
	}

}
