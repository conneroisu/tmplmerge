package twerge

import (
	"regexp"
	"strconv"
)

var (
	stringLengths = map[string]bool{
		"px":     true,
		"full":   true,
		"screen": true,
	}
	lengthUnitRegex = regexp.MustCompile(`\d+(%|px|r?em|[sdl]?v([hwib]|min|max)|pt|pc|in|cm|mm|cap|ch|ex|r?lh|cq(w|h|i|b|min|max))|\b(calc|min|max|clamp)\(.+\)|^0$`)
	colorFnRegex    = regexp.MustCompile(`^(rgba?|hsla?|hwb|(ok)?(lab|lch))\(.+\)$`)
	arbitraryRegex  = regexp.MustCompile(`(?i)^\[(?:([a-z-]+):)?(.+)\]$`)
	shirtPattern    = regexp.MustCompile(`^(\d+(\.\d+)?)?(xs|sm|md|lg|xl)$`)
	shardowPattern  = regexp.MustCompile(`^(inset_)?-?((\d+)?\.?(\d+)[a-z]+|0)_-?((\d+)?\.?(\d+)[a-z]+|0)`)

	sizeLabels  = map[string]bool{"length": true, "size": true, "percentage": true}
	imageLabels = map[string]bool{"image": true, "url": true}
)

// config is the configuration for the template merger
type config struct {
	// defaults should be good enough
	// hover:bg-red-500 -> :
	ModifierSeparator rune
	// bg-red-500 -> -
	ClassSeparator rune
	// !bg-red-500 -> !
	ImportantModifier rune
	// used for bg-red-500/50 (50% opacity) -> /
	PostfixModifier rune
	// optional
	Prefix string
	// CACHE
	MaxCacheSize int
	// This is a large map of all the classes and their validators -> see default-config.go
	ClassGroups classPart
	// class group with conflict + conflicting groups -> if "p" is set all others are removed
	// p: ['px', 'py', 'ps', 'pe', 'pt', 'pr', 'pb', 'pl']
	ConflictingClassGroups conflictingClassGroups
}

// classGroupValidator is a validator for a class group
type classGroupValidator struct {
	Fn           func(string) bool
	ClassGroupID string
}

// classPart is a part of a class group
type classPart struct {
	NextPart     map[string]classPart
	Validators   []classGroupValidator
	ClassGroupID string
}

// conflictingClassGroups is a map of class groups that conflict with each other
type conflictingClassGroups map[string][]string

func getBreaks(groupID string) map[string]classPart {
	return map[string]classPart{
		"auto": {
			NextPart:     map[string]classPart{},
			Validators:   []classGroupValidator{},
			ClassGroupID: groupID,
		},
		"avoid": {
			NextPart:     make(map[string]classPart),
			Validators:   []classGroupValidator{},
			ClassGroupID: groupID,
		},
		"all": {
			NextPart:     map[string]classPart{},
			Validators:   []classGroupValidator{},
			ClassGroupID: groupID,
		},
		"page": {
			NextPart:     map[string]classPart{},
			Validators:   []classGroupValidator{},
			ClassGroupID: groupID,
		},
		"left": {
			NextPart:     map[string]classPart{},
			Validators:   []classGroupValidator{},
			ClassGroupID: groupID,
		},
		"right": {
			NextPart:     map[string]classPart{},
			Validators:   []classGroupValidator{},
			ClassGroupID: groupID,
		},
		"column": {
			NextPart:     map[string]classPart{},
			Validators:   []classGroupValidator{},
			ClassGroupID: groupID,
		},
	}
}

func isAny(_ string) bool {
	return true
}

func isNever(_ string) bool {
	return false
}

func isLength(val string) bool {
	if isNumber(val) || stringLengths[val] || isFraction(val) {
		return true
	}
	return false
}

func isArbitraryLength(val string) bool {
	return labelIsArbitraryValue(val, "length", isLengthOnly)
}

// isArbitraryNumber returns true if the given value is an arbitrary number
func isArbitraryNumber(val string) bool {
	return labelIsArbitraryValue(val, "number", isNumber)
}

// isArbitraryPosition returns true if the given value is an arbitrary position
func isArbitraryPosition(val string) bool {
	return labelIsArbitraryValue(val, "position", isNever)
}

// isArbitrarySize returns true if the given value is an arbitrary size
func isArbitrarySize(val string) bool {
	return labelIsArbitraryValue(val, sizeLabels, isNever)
}

func isArbitraryImage(val string) bool {
	return labelIsArbitraryValue(val, imageLabels, isImage)
}
func isArbitraryShadow(val string) bool {
	return labelIsArbitraryValue(val, "", isShadow)
}

func isArbitraryValue(val string) bool {
	return arbitraryRegex.MatchString(val)
}

func isPercent(val string) bool {
	return val[len(val)-1] == '%' && isNumber(val[:len(val)-1])
}

func isTshirtSize(val string) bool {
	return shirtPattern.MatchString(val)
}

func isShadow(val string) bool {
	return shardowPattern.MatchString(val)
}

func isImage(val string) bool {
	pattern := regexp.MustCompile(`^(url|image|image-set|cross-fade|element|(repeating-)?(linear|radial|conic)-gradient)\(.+\)$`)
	return pattern.MatchString(val)
}

func isFraction(val string) bool {
	pattern := regexp.MustCompile(`^\d+\/\d+$`)
	return pattern.MatchString(val)
}

func isNumber(val string) bool {
	return isInteger(val) || isFloat(val)
}

func isInteger(val string) bool {
	_, err := strconv.Atoi(val)
	return err == nil
}

func isFloat(val string) bool {
	_, err := strconv.ParseFloat(val, 64)
	return err == nil
}

func isLengthOnly(val string) bool {
	return lengthUnitRegex.MatchString(val) && !colorFnRegex.MatchString(val)
}

// labelIsArbitraryValue returns true if the given value is an arbitrary value
// with the given label. The label can be a string, a map[string]bool or a
// function that takes a string and returns a bool.
func labelIsArbitraryValue(
	val string,
	label any,
	testValue func(string) bool,
) bool {
	res := arbitraryRegex.FindStringSubmatch(val)
	if len(res) > 1 {
		if res[1] != "" {
			if t, ok := label.(string); ok {
				return res[1] == t
			}

			if t, ok := label.(map[string]bool); ok {
				return t[res[1]]
			}
		}
		if len(res) > 2 {
			return testValue(res[2])
		}
	}
	return false
}

// defaultConfig is the default TwMergeConfig
var defaultConfig = &config{
	ModifierSeparator: ':',
	ClassSeparator:    '-',
	ImportantModifier: '!',
	PostfixModifier:   '/',
	MaxCacheSize:      1000,
	ConflictingClassGroups: conflictingClassGroups{
		"overflow":         {"overflow-x", "overflow-y"},
		"overscroll":       {"overscroll-x", "overscroll-y"},
		"inset":            {"inset-x", "inset-y", "start", "end", "top", "right", "bottom", "left"},
		"inset-x":          {"right", "left"},
		"inset-y":          {"top", "bottom"},
		"flex":             {"basis", "grow", "shrink"},
		"gap":              {"gap-x", "gap-y"},
		"p":                {"px", "py", "ps", "pe", "pt", "pr", "pb", "pl"},
		"px":               {"pr", "pl"},
		"py":               {"pt", "pb"},
		"m":                {"mx", "my", "ms", "me", "mt", "mr", "mb", "ml"},
		"mx":               {"mr", "ml"},
		"my":               {"mt", "mb"},
		"size":             {"w", "h"},
		"font-size":        {"leading"},
		"fvn-normal":       {"fvn-ordinal", "fvn-slashed-zero", "fvn-figure", "fvn-spacing", "fvn-fraction"},
		"fvn-ordinal":      {"fvn-normal"},
		"fvn-slashed-zero": {"fvn-normal"},
		"fvn-figure":       {"fvn-normal"},
		"fvn-spacing":      {"fvn-normal"},
		"fvn-fraction":     {"fvn-normal"},
		"line-clamp":       {"display", "overflow"},
		"rounded":          {"rounded-s", "rounded-e", "rounded-t", "rounded-r", "rounded-b", "rounded-l", "rounded-ss", "rounded-se", "rounded-ee", "rounded-es", "rounded-tl", "rounded-tr", "rounded-br", "rounded-bl"},
		"rounded-s":        {"rounded-ss", "rounded-es"},
		"rounded-e":        {"rounded-se", "rounded-ee"},
		"rounded-t":        {"rounded-tl", "rounded-tr"},
		"rounded-r":        {"rounded-tr", "rounded-br"},
		"rounded-b":        {"rounded-br", "rounded-bl"},
		"rounded-l":        {"rounded-tl", "rounded-bl"},
		"border-spacing":   {"border-spacing-x", "border-spacing-y"},
		"border-w":         {"border-w-s", "border-w-e", "border-w-t", "border-w-r", "border-w-b", "border-w-l"},
		"border-w-x":       {"border-w-r", "border-w-l"},
		"border-w-y":       {"border-w-t", "border-w-b"},
		"border-color":     {"border-color-t", "border-color-r", "border-color-b", "border-color-l"},
		"border-color-x":   {"border-color-r", "border-color-l"},
		"border-color-y":   {"border-color-t", "border-color-b"},
		"scroll-m":         {"scroll-mx", "scroll-my", "scroll-ms", "scroll-me", "scroll-mt", "scroll-mr", "scroll-mb", "scroll-ml"},
		"scroll-mx":        {"scroll-mr", "scroll-ml"},
		"scroll-my":        {"scroll-mt", "scroll-mb"},
		"scroll-p":         {"scroll-px", "scroll-py", "scroll-ps", "scroll-pe", "scroll-pt", "scroll-pr", "scroll-pb", "scroll-pl"},
		"scroll-px":        {"scroll-pr", "scroll-pl"},
		"scroll-py":        {"scroll-pt", "scroll-pb"},
		"touch":            {"touch-x", "touch-y", "touch-pz"},
		"touch-x":          {"touch"},
		"touch-y":          {"touch"},
		"touch-pz":         {"touch"},
	},
	ClassGroups: classPart{
		NextPart: map[string]classPart{
			// Aspect Ratio
			// @see https://tailwindcss.com/docs/aspect-ratio
			"aspect": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "aspect",
					},
					"square": {
						ClassGroupID: "aspect",
					},
					"video": {
						ClassGroupID: "aspect",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "aspect",
					},
				},
			},
			// Container
			// @see https://tailwindcss.com/docs/container
			"container": {
				NextPart:     map[string]classPart{},
				ClassGroupID: "container",
			},
			// Columns
			// @see https://tailwindcss.com/docs/columns
			"columns": {
				NextPart: map[string]classPart{},
				Validators: []classGroupValidator{
					{
						Fn:           isTshirtSize,
						ClassGroupID: "columns",
					},
				},
			},
			"break": {
				NextPart: map[string]classPart{
					// Break After
					// @see https://tailwindcss.com/docs/break-after
					"after": {
						NextPart: getBreaks("break-after"),
					},

					// Break Before @see https://tailwindcss.com/docs/break-before
					"before": {
						NextPart: getBreaks("break-before"),
					},

					// Break Inside
					// @see https://tailwindcss.com/docs/break-inside
					"inside": {
						NextPart: map[string]classPart{
							"auto": {
								ClassGroupID: "break-inside",
							},
							"avoid": {
								NextPart: map[string]classPart{
									"page": {
										ClassGroupID: "break-inside",
									},
									"column": {
										ClassGroupID: "break-inside",
									},
								},
								ClassGroupID: "break-inside",
							},
						},
					},

					// Word Break
					// @see https://tailwindcss.com/docs/word-break
					"normal": {
						ClassGroupID: "break",
					},
					"words": {
						ClassGroupID: "break",
					},
					"all": {
						ClassGroupID: "break",
					},
					"keep": {
						ClassGroupID: "break",
					},
				},
				Validators: []classGroupValidator{},
			},

			"box": {
				NextPart: map[string]classPart{
					// Box Sizing
					// @see https://tailwindcss.com/docs/box-sizing
					"border": {
						ClassGroupID: "box",
					},
					"content": {
						ClassGroupID: "box",
					},

					// Box Decoration Break
					// @see https://tailwindcss.com/docs/box-decoration-break
					"decoration": {
						NextPart: map[string]classPart{
							"slice": {
								ClassGroupID: "box-decoration"},
							"clone": {
								ClassGroupID: "box-decoration",
							},
						},
					},
				},
			},

			// Display
			// @see https://tailwindcss.com/docs/display
			"block": {
				ClassGroupID: "display",
			},
			"inline": {
				NextPart: map[string]classPart{
					"block": {ClassGroupID: "display"},
					"flex":  {ClassGroupID: "display"},
					"grid":  {ClassGroupID: "display"},
					"table": {ClassGroupID: "display"},
				},
				ClassGroupID: "display",
			},
			"flex": {
				NextPart: map[string]classPart{
					"row": {
						NextPart: map[string]classPart{
							"reverse": {
								ClassGroupID: "flex-direction",
							},
						},
						ClassGroupID: "flex-direction",
					},
					"col": {
						NextPart: map[string]classPart{
							"reverse": {
								ClassGroupID: "flex-direction",
							},
						},
						ClassGroupID: "flex-direction",
					},
					"wrap": {
						NextPart: map[string]classPart{
							"reverse": {
								ClassGroupID: "flex-wrap",
							},
						},
						ClassGroupID: "flex-wrap",
					},
					"nowrap": {
						ClassGroupID: "flex-wrap",
					},
					"1": {
						ClassGroupID: "flex",
					},
					"auto": {
						ClassGroupID: "flex",
					},
					"initial": {
						ClassGroupID: "flex",
					},
					"none": {
						ClassGroupID: "flex",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "flex",
					},
				},
				ClassGroupID: "display",
			},
			"table": {
				NextPart: map[string]classPart{
					"caption": {
						ClassGroupID: "display",
					},
					"cell": {
						ClassGroupID: "display",
					},
					"column": {
						NextPart: map[string]classPart{
							"group": {
								ClassGroupID: "display",
							},
						},
						ClassGroupID: "display",
					},
					"footer": {
						NextPart: map[string]classPart{
							"group": {
								ClassGroupID: "display",
							},
						},
					},
					"header": {
						NextPart: map[string]classPart{
							"group": {
								ClassGroupID: "display",
							},
						},
					},
					"row": {
						NextPart: map[string]classPart{
							"group": {
								ClassGroupID: "display",
							},
						},
						ClassGroupID: "display",
					},
					"auto": {
						ClassGroupID: "table-layout",
					},
					"fixed": {
						ClassGroupID: "table-layout",
					},
				},
				ClassGroupID: "display",
			},
			"flow": {
				NextPart: map[string]classPart{"root": {ClassGroupID: "display"}},
			},
			"grid": {
				NextPart: map[string]classPart{
					"cols": {
						Validators: []classGroupValidator{
							{
								Fn:           isAny,
								ClassGroupID: "grid-cols",
							},
						},
					},
					"rows": {
						Validators: []classGroupValidator{
							{
								Fn:           isAny,
								ClassGroupID: "grid-rows",
							},
						},
					},
					"flow": {
						NextPart: map[string]classPart{
							"row": {
								NextPart: map[string]classPart{
									"dense": {
										ClassGroupID: "grid-flow",
									},
								},
								ClassGroupID: "grid-flow",
							},
							"col": {
								NextPart: map[string]classPart{
									"dense": {
										ClassGroupID: "grid-flow",
									},
								},
								ClassGroupID: "grid-flow",
							},
							"dense": {
								ClassGroupID: "grid-flow",
							},
						},
					},
				},
				Validators:   []classGroupValidator{},
				ClassGroupID: "display",
			},
			"contents": {ClassGroupID: "display"},
			"list": {
				NextPart: map[string]classPart{
					"item": {
						ClassGroupID: "display",
					},
					"image": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "list-image",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "list-image",
							},
						},
					},
					"none": {
						ClassGroupID: "list-style-type",
					},
					"disc": {
						ClassGroupID: "list-style-type",
					},
					"decimal": {
						ClassGroupID: "list-style-type",
					},
					"inside": {
						ClassGroupID: "list-style-position",
					},
					"outside": {
						ClassGroupID: "list-style-position",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn: func(s string) bool {
							return s == "none" || s == "disc" || s == "decimal" || s == "square" || s == "circle"
						},
						ClassGroupID: "list-style-type",
					},
				},
			},
			"hidden": {ClassGroupID: "display"},
			"float": {
				NextPart: map[string]classPart{
					"right": {
						ClassGroupID: "float",
					},
					"left": {
						ClassGroupID: "float",
					},
					"none": {
						ClassGroupID: "float",
					},
					"start": {
						ClassGroupID: "float",
					},
					"end": {
						ClassGroupID: "float",
					},
				},
			},
			"clear": {
				NextPart: map[string]classPart{
					"left": {
						ClassGroupID: "clear",
					},
					"right": {
						ClassGroupID: "clear",
					},
					"both": {
						ClassGroupID: "clear",
					},
					"none": {
						ClassGroupID: "clear",
					},
					"start": {
						ClassGroupID: "clear",
					},
					"end": {
						ClassGroupID: "clear",
					},
				},
			},
			"isolate": {ClassGroupID: "isolation"},
			"isolation": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "isolation",
					},
				},
			},
			"object": {
				NextPart: map[string]classPart{
					"contain": {
						ClassGroupID: "object-fit",
					},
					"cover": {
						ClassGroupID: "object-fit",
					},
					"fill": {
						ClassGroupID: "object-fit",
					},
					"none": {
						ClassGroupID: "object-fit",
					},
					"scale": {
						NextPart: map[string]classPart{
							"down": {
								ClassGroupID: "object-fit",
							},
						},
					},
					"bottom": {
						ClassGroupID: "object-position",
					},
					"center": {
						ClassGroupID: "object-position",
					},
					"left": {
						NextPart: map[string]classPart{
							"bottom": {
								ClassGroupID: "object-position",
							},
							"top": {
								ClassGroupID: "object-position",
							},
						},
					},
					"right": {
						NextPart: map[string]classPart{
							"bottom": {
								ClassGroupID: "object-position",
							},
							"top": {
								ClassGroupID: "object-position",
							},
						},
					},
					"top": {
						ClassGroupID: "object-position",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "object-position",
					},
				},
			},

			"overflow": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "overflow",
					},
					"hidden": {
						ClassGroupID: "overflow",
					},
					"clip": {
						ClassGroupID: "overflow",
					},
					"visible": {
						ClassGroupID: "overflow",
					},
					"scroll": {
						ClassGroupID: "overflow",
					},
					"x": {
						NextPart: map[string]classPart{
							"auto": {
								ClassGroupID: "overflow-x",
							},
							"hidden": {
								ClassGroupID: "overflow-x",
							},
							"clip": {
								ClassGroupID: "overflow-x",
							},
							"visible": {
								ClassGroupID: "overflow-x",
							},
							"scroll": {
								ClassGroupID: "overflow-x",
							},
						},
					},
					"y": {
						NextPart: map[string]classPart{
							"auto": {
								ClassGroupID: "overflow-y",
							},
							"hidden": {
								ClassGroupID: "overflow-y",
							},
							"clip": {
								ClassGroupID: "overflow-y",
							},
							"visible": {
								ClassGroupID: "overflow-y",
							},
							"scroll": {
								ClassGroupID: "overflow-y",
							},
						},
					},
				},
			},
			"overscroll": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "overscroll",
					},
					"contain": {
						ClassGroupID: "overscroll",
					},
					"none": {
						ClassGroupID: "overscroll",
					},
					"x": {
						NextPart: map[string]classPart{
							"auto": {
								ClassGroupID: "overscroll-x",
							},
							"contain": {
								ClassGroupID: "overscroll-x",
							},
							"none": {
								ClassGroupID: "overscroll-x",
							},
						},
					},
					"y": {
						NextPart: map[string]classPart{
							"auto": {
								ClassGroupID: "overscroll-y",
							},
							"contain": {
								ClassGroupID: "overscroll-y",
							},
							"none": {
								ClassGroupID: "overscroll-y",
							},
						},
					},
				},
				Validators: []classGroupValidator{},
			},

			"static": {
				ClassGroupID: "position",
			},
			"fixed": {
				ClassGroupID: "position",
			},
			"absolute": {
				ClassGroupID: "position",
			},
			"relative": {
				ClassGroupID: "position",
			},
			"sticky": {
				ClassGroupID: "position",
			},

			"inset": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "inset",
					},
					"x": {
						NextPart: map[string]classPart{
							"auto": {
								ClassGroupID: "inset-x",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "inset-x",
							},
							{
								Fn:           isLength,
								ClassGroupID: "inset-x",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "inset-x",
							},
						},
					},
					"y": {
						NextPart: map[string]classPart{
							"auto": {
								ClassGroupID: "inset-y",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "inset-y",
							},
							{
								Fn:           isLength,
								ClassGroupID: "inset-y",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "inset-y",
							},
						},
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "inset",
					},
					{
						Fn:           isLength,
						ClassGroupID: "inset",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "inset",
					},
				},
			},
			"start": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "start",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "start",
					},
					{
						Fn:           isLength,
						ClassGroupID: "start",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "start",
					},
				},
			},
			"end": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "end",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "end",
					},
					{
						Fn:           isLength,
						ClassGroupID: "end",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "end",
					},
				},
			},
			"top": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "top",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "top",
					},
					{
						Fn:           isLength,
						ClassGroupID: "top",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "top",
					},
				},
			},
			"right": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "right",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "right",
					},
					{
						Fn:           isLength,
						ClassGroupID: "right",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "right",
					},
				},
			},
			"bottom": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "bottom",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "bottom",
					},
					{
						Fn:           isLength,
						ClassGroupID: "bottom",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "bottom",
					},
				},
			},
			"left": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "left",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "left",
					},
					{
						Fn:           isLength,
						ClassGroupID: "left",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "left",
					},
				},
			},
			"visible": {
				ClassGroupID: "visibility",
			},
			"invisible": {
				ClassGroupID: "visibility",
			},
			"collapse": {
				ClassGroupID: "visibility",
			},
			"z": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "z",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isInteger,
						ClassGroupID: "z",
					},
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "z",
					},
				},
			},
			"basis": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "basis",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "basis",
					},
					{
						Fn:           isLength,
						ClassGroupID: "basis",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "basis",
					},
				},
			},
			"grow": {
				NextPart: map[string]classPart{
					"0": {
						ClassGroupID: "grow",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "grow",
					},
				},
				ClassGroupID: "grow",
			},
			"shrink": {
				NextPart: map[string]classPart{
					"0": {
						ClassGroupID: "shrink",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "shrink",
					},
				},
				ClassGroupID: "shrink",
			},
			"order": {
				NextPart: map[string]classPart{
					"first": {
						ClassGroupID: "order",
					},
					"last": {
						ClassGroupID: "order",
					},
					"none": {
						ClassGroupID: "order",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isInteger,
						ClassGroupID: "order",
					},
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "order",
					},
				},
			},
			"col": {
				NextPart: map[string]classPart{
					"auto": {
						NextPart:     map[string]classPart{},
						Validators:   []classGroupValidator{},
						ClassGroupID: "col-start-end",
					},
					"span": {
						NextPart: map[string]classPart{
							"full": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "col-start-end",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isInteger,
								ClassGroupID: "col-start-end",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "col-start-end",
							},
						},
					},
					"start": {
						NextPart: map[string]classPart{
							"auto": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "col-start",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "col-start",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "col-start",
							},
						},
					},
					"end": {
						NextPart: map[string]classPart{
							"auto": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "col-end",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "col-end",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "col-end",
							},
						},
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "col-start-end",
					},
				},
			},
			"row": {
				NextPart: map[string]classPart{
					"auto": {
						NextPart:     map[string]classPart{},
						Validators:   []classGroupValidator{},
						ClassGroupID: "row-start-end",
					},
					"span": {
						NextPart: map[string]classPart{},
						Validators: []classGroupValidator{
							{
								Fn:           isInteger,
								ClassGroupID: "row-start-end",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "row-start-end",
							},
						},
					},
					"start": {
						NextPart: map[string]classPart{
							"auto": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "row-start",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "row-start",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "row-start",
							},
						},
					},
					"end": {
						NextPart: map[string]classPart{
							"auto": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "row-end",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "row-end",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "row-end",
							},
						},
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "row-start-end",
					},
				},
			},
			"auto": {
				NextPart: map[string]classPart{
					"cols": {
						NextPart: map[string]classPart{
							"auto": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "auto-cols",
							},
							"min": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "auto-cols",
							},
							"max": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "auto-cols",
							},
							"fr": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "auto-cols",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "auto-cols",
							},
						},
					},
					"rows": {
						NextPart: map[string]classPart{
							"auto": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "auto-rows",
							},
							"min": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "auto-rows",
							},
							"max": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "auto-rows",
							},
							"fr": {
								NextPart:     map[string]classPart{},
								Validators:   []classGroupValidator{},
								ClassGroupID: "auto-rows",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "auto-rows",
							},
						},
					},
				},
				Validators: []classGroupValidator{},
			},
			"gap": {
				NextPart: map[string]classPart{
					"x": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "gap-x",
							},
							{
								Fn:           isLength,
								ClassGroupID: "gap-x",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "gap-x",
							},
						},
					},
					"y": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "gap-y",
							},
							{
								Fn:           isLength,
								ClassGroupID: "gap-y",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "gap-y",
							},
						},
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "gap",
					},
					{
						Fn:           isLength,
						ClassGroupID: "gap",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "gap",
					},
				},
			},
			"justify": {
				NextPart: map[string]classPart{
					"normal": {
						ClassGroupID: "justify-content",
					},
					"start": {
						ClassGroupID: "justify-content",
					},
					"end": {
						ClassGroupID: "justify-content",
					},
					"center": {
						ClassGroupID: "justify-content",
					},
					"between": {
						ClassGroupID: "justify-content",
					},
					"around": {
						ClassGroupID: "justify-content",
					},
					"evenly": {
						ClassGroupID: "justify-content",
					},
					"stretch": {
						ClassGroupID: "justify-content",
					},
					"items": {
						NextPart: map[string]classPart{
							"start": {
								ClassGroupID: "justify-items",
							},
							"end": {
								ClassGroupID: "justify-items",
							},
							"center": {
								ClassGroupID: "justify-items",
							},
							"stretch": {
								ClassGroupID: "justify-items",
							},
						},
					},
					"self": {
						NextPart: map[string]classPart{
							"auto": {
								ClassGroupID: "justify-self",
							},
							"start": {
								ClassGroupID: "justify-self",
							},
							"end": {
								ClassGroupID: "justify-self",
							},
							"center": {
								ClassGroupID: "justify-self",
							},
							"stretch": {
								ClassGroupID: "justify-self",
							},
						},
					},
				},
			},
			"content": {
				NextPart: map[string]classPart{
					"normal": {
						ClassGroupID: "align-content",
					},
					"start": {
						ClassGroupID: "align-content",
					},
					"end": {
						ClassGroupID: "align-content",
					},
					"center": {
						ClassGroupID: "align-content",
					},
					"between": {
						ClassGroupID: "align-content",
					},
					"around": {
						ClassGroupID: "align-content",
					},
					"evenly": {
						ClassGroupID: "align-content",
					},
					"stretch": {
						ClassGroupID: "align-content",
					},
					"baseline": {
						ClassGroupID: "align-content",
					},
					"none": {
						ClassGroupID: "content",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "content",
					},
				},
			},
			"items": {
				NextPart: map[string]classPart{
					"start": {
						ClassGroupID: "align-items",
					},
					"end": {
						ClassGroupID: "align-items",
					},
					"center": {
						ClassGroupID: "align-items",
					},
					"baseline": {
						ClassGroupID: "align-items",
					},
					"stretch": {
						ClassGroupID: "align-items",
					},
				},
			},
			"self": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "align-self",
					},
					"start": {
						ClassGroupID: "align-self",
					},
					"end": {
						ClassGroupID: "align-self",
					},
					"center": {
						ClassGroupID: "align-self",
					},
					"stretch": {
						ClassGroupID: "align-self",
					},
					"baseline": {
						ClassGroupID: "align-self",
					},
				},
			},
			"place": {
				NextPart: map[string]classPart{
					"content": {
						NextPart: map[string]classPart{
							"start": {
								ClassGroupID: "place-content",
							},
							"end": {
								ClassGroupID: "place-content",
							},
							"center": {
								ClassGroupID: "place-content",
							},
							"between": {
								ClassGroupID: "place-content",
							},
							"around": {
								ClassGroupID: "place-content",
							},
							"evenly": {
								ClassGroupID: "place-content",
							},
							"stretch": {
								ClassGroupID: "place-content",
							},
							"baseline": {
								ClassGroupID: "place-content",
							},
						},
					},
					"items": {
						NextPart: map[string]classPart{
							"start": {
								ClassGroupID: "place-items",
							},
							"end": {
								ClassGroupID: "place-items",
							},
							"center": {
								ClassGroupID: "place-items",
							},
							"baseline": {
								ClassGroupID: "place-items",
							},
							"stretch": {
								ClassGroupID: "place-items",
							},
						},
					},
					"self": {
						NextPart: map[string]classPart{
							"auto": {
								ClassGroupID: "place-self",
							},
							"start": {
								ClassGroupID: "place-self",
							},
							"end": {
								ClassGroupID: "place-self",
							},
							"center": {
								ClassGroupID: "place-self",
							},
							"stretch": {
								ClassGroupID: "place-self",
							},
						},
					},
				},
			},
			"p": {
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "p",
					},
					{
						Fn:           isLength,
						ClassGroupID: "p",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "p",
					},
				},
			},
			"px": {
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "px",
					},
					{
						Fn:           isLength,
						ClassGroupID: "px",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "px",
					},
				},
			},
			"py": {
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "py",
					},
					{
						Fn:           isLength,
						ClassGroupID: "py",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "py",
					},
				},
			},
			"ps": {
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "ps",
					},
					{
						Fn:           isLength,
						ClassGroupID: "ps",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "ps",
					},
				},
			},
			"pe": {
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "pe",
					},
					{
						Fn:           isLength,
						ClassGroupID: "pe",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "pe",
					},
				},
			},
			"pt": {
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "pt",
					},
					{
						Fn:           isLength,
						ClassGroupID: "pt",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "pt",
					},
				},
			},
			"pr": {
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "pr",
					},
					{
						Fn:           isLength,
						ClassGroupID: "pr",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "pr",
					},
				},
			},
			"pb": {
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "pb",
					},
					{
						Fn:           isLength,
						ClassGroupID: "pb",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "pb",
					},
				},
			},
			"pl": {
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "pl",
					},
					{
						Fn:           isLength,
						ClassGroupID: "pl",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "pl",
					},
				},
			},
			"m": {
				NextPart: map[string]classPart{
					"auto": {
						Validators:   []classGroupValidator{},
						ClassGroupID: "m",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "m",
					},
					{
						Fn:           isLength,
						ClassGroupID: "m",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "m",
					},
				},
			},
			"mx": {
				NextPart: map[string]classPart{
					"auto": {
						Validators:   []classGroupValidator{},
						ClassGroupID: "mx",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "mx",
					},
					{
						Fn:           isLength,
						ClassGroupID: "mx",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "mx",
					},
				},
			},
			"my": {
				NextPart: map[string]classPart{
					"auto": {
						Validators:   []classGroupValidator{},
						ClassGroupID: "my",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "my",
					},
					{
						Fn:           isLength,
						ClassGroupID: "my",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "my",
					},
				},
			},
			"ms": {
				NextPart: map[string]classPart{
					"auto": {
						Validators:   []classGroupValidator{},
						ClassGroupID: "ms",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "ms",
					},
					{
						Fn:           isLength,
						ClassGroupID: "ms",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "ms",
					},
				},
			},
			"me": {
				NextPart: map[string]classPart{
					"auto": {
						Validators:   []classGroupValidator{},
						ClassGroupID: "me",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "me",
					},
					{
						Fn:           isLength,
						ClassGroupID: "me",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "me",
					},
				},
			},
			"mt": {
				NextPart: map[string]classPart{
					"auto": {
						Validators:   []classGroupValidator{},
						ClassGroupID: "mt",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "mt",
					},
					{
						Fn:           isLength,
						ClassGroupID: "mt",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "mt",
					},
				},
			},
			"mr": {
				NextPart: map[string]classPart{
					"auto": {
						Validators:   []classGroupValidator{},
						ClassGroupID: "mr",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "mr",
					},
					{
						Fn:           isLength,
						ClassGroupID: "mr",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "mr",
					},
				},
			},
			"mb": {
				NextPart: map[string]classPart{
					"auto": {
						Validators:   []classGroupValidator{},
						ClassGroupID: "mb",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "mb",
					},
					{
						Fn:           isLength,
						ClassGroupID: "mb",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "mb",
					},
				},
			},
			"ml": {
				NextPart: map[string]classPart{
					"auto": {
						Validators:   []classGroupValidator{},
						ClassGroupID: "ml",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "ml",
					},
					{
						Fn:           isLength,
						ClassGroupID: "ml",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "ml",
					},
				},
			},
			"space": {
				NextPart: map[string]classPart{
					"x": {
						NextPart: map[string]classPart{
							"reverse": {
								Validators:   []classGroupValidator{},
								ClassGroupID: "space-x-reverse",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "space-x",
							},
							{
								Fn:           isLength,
								ClassGroupID: "space-x",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "space-x",
							},
						},
					},
					"y": {
						NextPart: map[string]classPart{
							"reverse": {
								ClassGroupID: "space-y-reverse",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "space-y",
							},
							{
								Fn:           isLength,
								ClassGroupID: "space-y",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "space-y",
							},
						},
					},
				},
				Validators: []classGroupValidator{},
			},
			"w": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "w",
					},
					"min": {
						ClassGroupID: "w",
					},
					"max": {
						ClassGroupID: "w",
					},
					"fit": {
						ClassGroupID: "w",
					},
					"svw": {
						ClassGroupID: "w",
					},
					"lvw": {
						ClassGroupID: "w",
					},
					"dvw": {
						ClassGroupID: "w",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "w",
					},
					{
						Fn:           isLength,
						ClassGroupID: "w",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "w",
					},
				},
			},
			"min": {
				NextPart: map[string]classPart{
					"w": {
						NextPart: map[string]classPart{
							"min": {
								ClassGroupID: "min-w",
							},
							"max": {
								ClassGroupID: "min-w",
							},
							"fit": {
								ClassGroupID: "min-w",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "min-w",
							},
							{
								Fn:           isLength,
								ClassGroupID: "min-w",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "min-w",
							},
						},
					},
					"h": {
						NextPart: map[string]classPart{
							"min": {
								ClassGroupID: "min-h",
							},
							"max": {
								ClassGroupID: "min-h",
							},
							"fit": {
								ClassGroupID: "min-h",
							},
							"svh": {
								ClassGroupID: "min-h",
							},
							"lvh": {
								ClassGroupID: "min-h",
							},
							"dvh": {
								ClassGroupID: "min-h",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "min-h",
							},
							{
								Fn:           isLength,
								ClassGroupID: "min-h",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "min-h",
							},
						},
					},
				},
				Validators: []classGroupValidator{},
			},
			"max": {
				NextPart: map[string]classPart{
					"w": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "max-w",
							},
							"full": {
								ClassGroupID: "max-w",
							},
							"min": {
								ClassGroupID: "max-w",
							},
							"max": {
								ClassGroupID: "max-w",
							},
							"fit": {
								ClassGroupID: "max-w",
							},
							"prose": {
								ClassGroupID: "max-w",
							},
							"screen": {
								Validators: []classGroupValidator{
									{
										Fn:           isTshirtSize,
										ClassGroupID: "max-w",
									},
								},
								ClassGroupID: "max-w",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "max-w",
							},
							{
								Fn:           isLength,
								ClassGroupID: "max-w",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "max-w",
							},
							{
								Fn:           isTshirtSize,
								ClassGroupID: "max-w",
							},
						},
						ClassGroupID: "max-w",
					},
					"h": {
						NextPart: map[string]classPart{
							"min": {
								ClassGroupID: "max-h",
							},
							"max": {
								ClassGroupID: "max-h",
							},
							"fit": {
								ClassGroupID: "max-h",
							},
							"svh": {
								ClassGroupID: "max-h",
							},
							"lvh": {
								ClassGroupID: "max-h",
							},
							"dvh": {
								ClassGroupID: "max-h",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "max-h",
							},
							{
								Fn:           isLength,
								ClassGroupID: "max-h",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "max-h",
							},
						},
						ClassGroupID: "max-h",
					},
				},
			},
			"h": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "h",
					},
					"min": {
						ClassGroupID: "h",
					},
					"max": {
						ClassGroupID: "h",
					},
					"fit": {
						ClassGroupID: "h",
					},
					"svh": {
						ClassGroupID: "h",
					},
					"lvh": {
						ClassGroupID: "h",
					},
					"dvh": {
						ClassGroupID: "h",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "h",
					},
					{
						Fn:           isLength,
						ClassGroupID: "h",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "h",
					},
				},
				ClassGroupID: "h",
			},
			"size": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "size",
					},
					"min": {
						ClassGroupID: "size",
					},
					"max": {
						ClassGroupID: "size",
					},
					"fit": {
						ClassGroupID: "size",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "size",
					},
					{
						Fn:           isLength,
						ClassGroupID: "size",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "size",
					},
				},
				ClassGroupID: "size",
			},
			"text": {
				NextPart: map[string]classPart{
					"base": {
						ClassGroupID: "font-size",
					},
					"left": {
						ClassGroupID: "text-alignment",
					},
					"center": {
						ClassGroupID: "text-alignment",
					},
					"right": {
						ClassGroupID: "text-alignment",
					},
					"justify": {
						ClassGroupID: "text-alignment",
					},
					"start": {
						ClassGroupID: "text-alignment",
					},
					"end": {
						ClassGroupID: "text-alignment",
					},
					"opacity": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "text-opacity",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "text-opacity",
							},
						},
						ClassGroupID: "text-opacity",
					},
					"ellipsis": {
						ClassGroupID: "text-overflow",
					},
					"clip": {
						ClassGroupID: "text-overflow",
					},
					"wrap": {
						ClassGroupID: "text-wrap",
					},
					"nowrap": {
						ClassGroupID: "text-wrap",
					},
					"balance": {
						ClassGroupID: "text-wrap",
					},
					"pretty": {
						ClassGroupID: "text-wrap",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isTshirtSize,
						ClassGroupID: "font-size",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "font-size",
					},
					{
						Fn:           isAny,
						ClassGroupID: "text-color",
					},
				},
			},
			"antialiased": {
				ClassGroupID: "font-smoothing",
			},
			"subpixel": {
				NextPart: map[string]classPart{
					"antialiased": {
						ClassGroupID: "font-smoothing",
					},
				},
			},
			"italic": {
				ClassGroupID: "font-style",
			},
			"not": {
				NextPart: map[string]classPart{
					"italic": {
						ClassGroupID: "font-style",
					},
					"sr": {
						NextPart: map[string]classPart{
							"only": {
								ClassGroupID: "sr",
							},
						},
					},
				},
			},
			"font": {
				NextPart: map[string]classPart{
					"thin": {
						ClassGroupID: "font-weight",
					},
					"extralight": {
						ClassGroupID: "font-weight",
					},
					"light": {
						ClassGroupID: "font-weight",
					},
					"normal": {
						ClassGroupID: "font-weight",
					},
					"medium": {
						ClassGroupID: "font-weight",
					},
					"semibold": {
						ClassGroupID: "font-weight",
					},
					"bold": {
						ClassGroupID: "font-weight",
					},
					"extrabold": {
						ClassGroupID: "font-weight",
					},
					"black": {
						ClassGroupID: "font-weight",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryNumber,
						ClassGroupID: "font-weight",
					},
					{
						Fn:           isAny,
						ClassGroupID: "font-family",
					},
				},
			},
			"normal": {
				NextPart: map[string]classPart{
					"nums": {
						ClassGroupID: "fvn-normal",
					},
					"case": {
						ClassGroupID: "text-transform",
					},
				},
			},
			"ordinal": {
				ClassGroupID: "fvn-ordinal",
			},
			"slashed": {
				NextPart: map[string]classPart{
					"zero": {
						ClassGroupID: "fvn-slashed-zero",
					},
				},
			},
			"lining": {
				NextPart: map[string]classPart{
					"nums": {
						ClassGroupID: "fvn-figure",
					},
				},
			},
			"oldstyle": {
				NextPart: map[string]classPart{
					"nums": {
						ClassGroupID: "fvn-figure",
					},
				},
			},
			"proportional": {
				NextPart: map[string]classPart{
					"nums": {
						ClassGroupID: "fvn-spacing",
					},
				},
			},
			"tabular": {
				NextPart: map[string]classPart{
					"nums": {
						ClassGroupID: "fvn-spacing",
					},
				},
			},
			"diagonal": {
				NextPart: map[string]classPart{
					"fractions": {
						ClassGroupID: "fvn-fraction",
					},
				},
			},
			"stacked": {
				NextPart: map[string]classPart{
					"fractons": {
						ClassGroupID: "fvn-fraction",
					},
				},
			},
			"tracking": {
				NextPart: map[string]classPart{
					"tighter": {
						ClassGroupID: "tracking",
					},
					"tight": {
						ClassGroupID: "tracking",
					},
					"normal": {
						ClassGroupID: "tracking",
					},
					"wide": {
						ClassGroupID: "tracking",
					},
					"wider": {
						ClassGroupID: "tracking",
					},
					"widest": {
						ClassGroupID: "tracking",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "tracking",
					},
				},
			},
			"line": {
				NextPart: map[string]classPart{
					"clamp": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "line-clamp",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "line-clamp",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "line-clamp",
							},
						},
					},
					"through": {
						ClassGroupID: "text-decoration",
					},
				},
			},
			"leading": {
				NextPart: map[string]classPart{
					"none": {
						ClassGroupID: "leading",
					},
					"tight": {
						ClassGroupID: "leading",
					},
					"snug": {
						ClassGroupID: "leading",
					},
					"normal": {
						ClassGroupID: "leading",
					},
					"relaxed": {
						ClassGroupID: "leading",
					},
					"loose": {
						ClassGroupID: "leading",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isLength,
						ClassGroupID: "leading",
					},
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "leading",
					},
				},
			},
			"placeholder": {
				NextPart: map[string]classPart{
					"opacity": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "placeholder-opacity",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "placeholder-opacity",
							},
						},
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isAny,
						ClassGroupID: "placeholder-color",
					},
				},
			},
			"underline": {
				NextPart: map[string]classPart{
					"offset": {
						NextPart: map[string]classPart{
							"auto": {
								ClassGroupID: "underline-offset",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "underline-offset",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "underline-offset",
							},
						},
					},
				},
				ClassGroupID: "text-decoration",
			},
			"overline": {
				ClassGroupID: "text-decoration",
			},
			"no": {
				NextPart: map[string]classPart{
					"underline": {
						ClassGroupID: "text-decoration",
					},
				},
			},
			"decoration": {
				NextPart: map[string]classPart{
					"solid": {
						ClassGroupID: "text-decoration-style",
					},
					"dashed": {
						ClassGroupID: "text-decoration-style",
					},
					"dotted": {
						ClassGroupID: "text-decoration-style",
					},
					"double": {
						ClassGroupID: "text-decoration-style",
					},
					"none": {
						ClassGroupID: "text-decoration-style",
					},
					"wavy": {
						ClassGroupID: "text-decoration-style",
					},
					"auto": {
						ClassGroupID: "text-decoration-thickness",
					},
					"from": {
						NextPart: map[string]classPart{
							"font": {
								ClassGroupID: "text-decoration-thickness",
							},
						},
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isLength,
						ClassGroupID: "text-decoration-thickness",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "text-decoration-thickness",
					},
					{
						Fn:           isAny,
						ClassGroupID: "text-decoration-color",
					},
				},
				ClassGroupID: "",
			},
			"uppercase": {
				ClassGroupID: "text-transform",
			},
			"lowercase": {
				ClassGroupID: "text-transform",
			},
			"capitalize": {
				ClassGroupID: "text-transform",
			},
			"truncate": {
				ClassGroupID: "text-overflow",
			},
			"indent": {
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "indent",
					},
					{
						Fn:           isLength,
						ClassGroupID: "indent",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "indent",
					},
				},
			},
			"align": {
				NextPart: map[string]classPart{
					"baseline": {
						ClassGroupID: "vertical-align",
					},
					"top": {
						ClassGroupID: "vertical-align",
					},
					"middle": {
						ClassGroupID: "vertical-align",
					},
					"bottom": {
						ClassGroupID: "vertical-align",
					},
					"text": {
						NextPart: map[string]classPart{
							"top": {
								ClassGroupID: "vertical-align",
							},
							"bottom": {
								ClassGroupID: "vertical-align",
							},
						},
					},
					"sub": {
						ClassGroupID: "vertical-align",
					},
					"super": {
						ClassGroupID: "vertical-align",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "vertical-align",
					},
				},
			},
			"whitespace": {
				NextPart: map[string]classPart{
					"normal": {
						ClassGroupID: "whitespace",
					},
					"nowrap": {
						ClassGroupID: "whitespace",
					},
					"pre": {
						NextPart: map[string]classPart{
							"line": {
								ClassGroupID: "whitespace",
							},
							"wrap": {
								ClassGroupID: "whitespace",
							},
						},
						ClassGroupID: "whitespace",
					},
					"break": {
						NextPart: map[string]classPart{
							"spaces": {
								ClassGroupID: "whitespace",
							},
						},
						ClassGroupID: "",
					},
				},
			},
			"hyphens": {
				NextPart: map[string]classPart{
					"none": {
						ClassGroupID: "hyphens",
					},
					"manual": {
						ClassGroupID: "hyphens",
					},
					"auto": {
						ClassGroupID: "hyphens",
					},
				},
			},
			"bg": {
				NextPart: map[string]classPart{
					"fixed": {
						ClassGroupID: "bg-attachment",
					},
					"local": {
						ClassGroupID: "bg-attachment",
					},
					"scroll": {
						ClassGroupID: "bg-attachment",
					},
					"clip": {
						NextPart: map[string]classPart{
							"border": {
								ClassGroupID: "bg-clip",
							},
							"padding": {
								ClassGroupID: "bg-clip",
							},
							"content": {
								ClassGroupID: "bg-clip",
							},
							"text": {
								ClassGroupID: "bg-clip",
							},
						},
					},
					"opacity": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "bg-opacity",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "bg-opacity",
							},
						},
					},
					"origin": {
						NextPart: map[string]classPart{
							"border": {
								ClassGroupID: "bg-origin",
							},
							"padding": {
								ClassGroupID: "bg-origin",
							},
							"content": {
								ClassGroupID: "bg-origin",
							},
						},
					},
					"bottom": {
						ClassGroupID: "bg-position",
					},
					"center": {
						ClassGroupID: "bg-position",
					},
					"left": {
						NextPart: map[string]classPart{
							"bottom": {
								ClassGroupID: "bg-position",
							},
							"top": {
								ClassGroupID: "bg-position",
							},
						},
						ClassGroupID: "bg-position",
					},
					"right": {
						NextPart: map[string]classPart{
							"bottom": {
								ClassGroupID: "bg-position",
							},
							"top": {
								ClassGroupID: "bg-position",
							},
						},
						ClassGroupID: "bg-position",
					},
					"top": {
						ClassGroupID: "bg-position",
					},
					"no": {
						NextPart: map[string]classPart{
							"repeat": {
								ClassGroupID: "bg-repeat",
							},
						},
					},
					"repeat": {
						NextPart: map[string]classPart{
							"x": {
								ClassGroupID: "bg-repeat",
							},
							"y": {
								ClassGroupID: "bg-repeat",
							},
							"round": {
								ClassGroupID: "bg-repeat",
							},
							"space": {
								ClassGroupID: "bg-repeat",
							},
						},
						ClassGroupID: "bg-repeat",
					},
					"auto": {
						ClassGroupID: "bg-size",
					},
					"cover": {
						ClassGroupID: "bg-size",
					},
					"contain": {
						ClassGroupID: "bg-size",
					},
					"none": {
						ClassGroupID: "bg-image",
					},
					"gradient": {
						NextPart: map[string]classPart{
							"to": {
								NextPart: map[string]classPart{
									"t": {
										ClassGroupID: "bg-image",
									},
									"tr": {
										ClassGroupID: "bg-image",
									},
									"r": {
										ClassGroupID: "bg-image",
									},
									"br": {
										ClassGroupID: "bg-image",
									},
									"b": {
										ClassGroupID: "bg-image",
									},
									"bl": {
										ClassGroupID: "bg-image",
									},
									"l": {
										ClassGroupID: "bg-image",
									},
									"tl": {
										ClassGroupID: "bg-image",
									},
								},
							},
						},
					},
					"blend": {
						NextPart: map[string]classPart{
							"normal": {
								ClassGroupID: "bg-blend",
							},
							"multiply": {
								ClassGroupID: "bg-blend",
							},
							"screen": {
								ClassGroupID: "bg-blend",
							},
							"overlay": {
								ClassGroupID: "bg-blend",
							},
							"darken": {
								ClassGroupID: "bg-blend",
							},
							"lighten": {
								ClassGroupID: "bg-blend",
							},
							"color": {
								NextPart: map[string]classPart{
									"dodge": {
										ClassGroupID: "bg-blend",
									},
									"burn": {
										ClassGroupID: "bg-blend",
									},
								},
							},
							"hard": {
								NextPart: map[string]classPart{
									"light": {
										ClassGroupID: "bg-blend",
									},
								},
							},
							"soft": {
								NextPart: map[string]classPart{
									"light": {
										ClassGroupID: "bg-blend",
									},
								},
							},
							"difference": {
								ClassGroupID: "bg-blend",
							},
							"exclusion": {
								ClassGroupID: "bg-blend",
							},
							"hue": {
								ClassGroupID: "bg-blend",
							},
							"saturation": {
								ClassGroupID: "bg-blend",
							},
							"luminosity": {
								ClassGroupID: "bg-blend",
							},
							"plus": {
								NextPart: map[string]classPart{
									"lighter": {
										ClassGroupID: "bg-blend",
									},
								},
							},
						},
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryPosition,
						ClassGroupID: "bg-position",
					},
					{
						Fn:           isArbitrarySize,
						ClassGroupID: "bg-size",
					},
					{
						Fn:           isArbitraryImage,
						ClassGroupID: "bg-image",
					},
					{
						Fn:           isAny,
						ClassGroupID: "bg-color",
					},
				},
			},
			"from": {
				Validators: []classGroupValidator{
					{
						Fn:           isPercent,
						ClassGroupID: "gradient-from-pos",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "gradient-from-pos",
					},
					{
						Fn:           isAny,
						ClassGroupID: "gradient-from",
					},
				},
			},
			"via": {
				Validators: []classGroupValidator{
					{
						Fn:           isPercent,
						ClassGroupID: "gradient-via-pos",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "gradient-via-pos",
					},
					{
						Fn:           isAny,
						ClassGroupID: "gradient-via",
					},
				},
			},
			"to": {
				Validators: []classGroupValidator{
					{
						Fn:           isPercent,
						ClassGroupID: "gradient-to-pos",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "gradient-to-pos",
					},
					{
						Fn:           isAny,
						ClassGroupID: "gradient-to",
					},
				},
			},
			"rounded": {
				NextPart: map[string]classPart{
					"none": {
						ClassGroupID: "rounded",
					},
					"full": {
						ClassGroupID: "rounded",
					},
					"s": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-s",
							},
							"full": {
								ClassGroupID: "rounded-s",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-s",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-s",
							},
						},
						ClassGroupID: "rounded-s",
					},
					"e": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-e",
							},
							"full": {
								ClassGroupID: "rounded-e",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-e",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-e",
							},
						},
						ClassGroupID: "rounded-e",
					},
					"t": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-t",
							},
							"full": {
								ClassGroupID: "rounded-t",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-t",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-t",
							},
						},
						ClassGroupID: "rounded-t",
					},
					"r": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-r",
							},
							"full": {
								ClassGroupID: "rounded-r",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-r",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-r",
							},
						},
						ClassGroupID: "rounded-r",
					},
					"b": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-b",
							},
							"full": {
								ClassGroupID: "rounded-b",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-b",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-b",
							},
						},
						ClassGroupID: "rounded-b",
					},
					"l": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-l",
							},
							"full": {
								ClassGroupID: "rounded-l",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-l",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-l",
							},
						},
						ClassGroupID: "rounded-l",
					},
					"ss": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-ss",
							},
							"full": {
								ClassGroupID: "rounded-ss",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-ss",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-ss",
							},
						},
						ClassGroupID: "rounded-ss",
					},
					"se": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-se",
							},
							"full": {
								ClassGroupID: "rounded-se",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-se",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-se",
							},
						},
						ClassGroupID: "rounded-se",
					},
					"ee": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-ee",
							},
							"full": {
								ClassGroupID: "rounded-ee",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-ee",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-ee",
							},
						},
						ClassGroupID: "rounded-ee",
					},
					"es": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-es",
							},
							"full": {
								ClassGroupID: "rounded-es",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-es",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-es",
							},
						},
						ClassGroupID: "rounded-es",
					},
					"tl": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-tl",
							},
							"full": {
								ClassGroupID: "rounded-tl",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-tl",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-tl",
							},
						},
						ClassGroupID: "rounded-tl",
					},
					"tr": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-tr",
							},
							"full": {
								ClassGroupID: "rounded-tr",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-tr",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-tr",
							},
						},
						ClassGroupID: "rounded-tr",
					},
					"br": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-br",
							},
							"full": {
								ClassGroupID: "rounded-br",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-br",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-br",
							},
						},
						ClassGroupID: "rounded-br",
					},
					"bl": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "rounded-bl",
							},
							"full": {
								ClassGroupID: "rounded-bl",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "rounded-bl",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "rounded-bl",
							},
						},
						ClassGroupID: "rounded-bl",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isTshirtSize,
						ClassGroupID: "rounded",
					},
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "rounded",
					},
				},
				ClassGroupID: "rounded",
			},
			"border": {
				NextPart: map[string]classPart{
					"x": {
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "border-w-x",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "border-w-x",
							},
							{
								Fn:           isAny,
								ClassGroupID: "border-color-x",
							},
						},
						ClassGroupID: "border-w-x",
					},
					"y": {
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "border-w-y",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "border-w-y",
							},
							{
								Fn:           isAny,
								ClassGroupID: "border-color-y",
							},
						},
						ClassGroupID: "border-w-y",
					},
					"s": {
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "border-w-s",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "border-w-s",
							},
						},
						ClassGroupID: "border-w-s",
					},
					"e": {
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "border-w-e",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "border-w-e",
							},
						},
						ClassGroupID: "border-w-e",
					},
					"t": {
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "border-w-t",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "border-w-t",
							},
							{
								Fn:           isAny,
								ClassGroupID: "border-color-t",
							},
						},
						ClassGroupID: "border-w-t",
					},
					"r": {
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "border-w-r",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "border-w-r",
							},
							{
								Fn:           isAny,
								ClassGroupID: "border-color-r",
							},
						},
						ClassGroupID: "border-w-r",
					},
					"b": {
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "border-w-b",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "border-w-b",
							},
							{
								Fn:           isAny,
								ClassGroupID: "border-color-b",
							},
						},
						ClassGroupID: "border-w-b",
					},
					"l": {
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "border-w-l",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "border-w-l",
							},
							{
								Fn:           isAny,
								ClassGroupID: "border-color-l",
							},
						},
						ClassGroupID: "border-w-l",
					},
					"opacity": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "border-opacity",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "border-opacity",
							},
						},
						ClassGroupID: "border-opacity",
					},
					"solid": {
						ClassGroupID: "border-style",
					},
					"dashed": {
						ClassGroupID: "border-style",
					},
					"dotted": {
						ClassGroupID: "border-style",
					},
					"double": {
						ClassGroupID: "border-style",
					},
					"none": {
						ClassGroupID: "border-style",
					},
					"hidden": {
						ClassGroupID: "border-style",
					},
					"collapse": {
						ClassGroupID: "border-collapse",
					},
					"separate": {
						ClassGroupID: "border-collapse",
					},
					"spacing": {
						NextPart: map[string]classPart{
							"x": {
								Validators: []classGroupValidator{
									{
										Fn:           isArbitraryValue,
										ClassGroupID: "border-spacing-x",
									},
									{
										Fn:           isLength,
										ClassGroupID: "border-spacing-x",
									},
									{
										Fn:           isArbitraryLength,
										ClassGroupID: "border-spacing-x",
									},
								},
							},
							"y": {
								Validators: []classGroupValidator{
									{
										Fn:           isArbitraryValue,
										ClassGroupID: "border-spacing-y",
									},
									{
										Fn:           isLength,
										ClassGroupID: "border-spacing-y",
									},
									{
										Fn:           isArbitraryLength,
										ClassGroupID: "border-spacing-y",
									},
								},
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "border-spacing",
							},
							{
								Fn:           isLength,
								ClassGroupID: "border-spacing",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "border-spacing",
							},
						},
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isLength,
						ClassGroupID: "border-w",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "border-w",
					},
					{
						Fn:           isAny,
						ClassGroupID: "border-color",
					},
				},
				ClassGroupID: "border-w",
			},
			"divide": {
				NextPart: map[string]classPart{
					"x": {
						NextPart: map[string]classPart{
							"reverse": {
								ClassGroupID: "divide-x-reverse",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "divide-x",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "divide-x",
							},
						},
						ClassGroupID: "divide-x",
					},
					"y": {
						NextPart: map[string]classPart{
							"reverse": {
								ClassGroupID: "divide-y-reverse",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "divide-y",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "divide-y",
							},
						},
						ClassGroupID: "divide-y",
					},
					"opacity": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "divide-opacity",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "divide-opacity",
							},
						},
					},
					"solid": {
						ClassGroupID: "divide-style",
					},
					"dashed": {
						ClassGroupID: "divide-style",
					},
					"dotted": {
						ClassGroupID: "divide-style",
					},
					"double": {
						ClassGroupID: "divide-style",
					},
					"none": {
						ClassGroupID: "divide-style",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isAny,
						ClassGroupID: "divide-color",
					},
				},
			},
			"outline": {
				NextPart: map[string]classPart{
					"solid": {
						ClassGroupID: "outline-style",
					},
					"dashed": {
						ClassGroupID: "outline-style",
					},
					"dotted": {
						ClassGroupID: "outline-style",
					},
					"double": {
						ClassGroupID: "outline-style",
					},
					"none": {
						ClassGroupID: "outline-style",
					},
					"offset": {
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "outline-offset",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "outline-offset",
							},
						},
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isLength,
						ClassGroupID: "outline-w",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "outline-w",
					},
					{
						Fn:           isAny,
						ClassGroupID: "outline-color",
					},
				},
				ClassGroupID: "outline-style",
			},
			"ring": {
				NextPart: map[string]classPart{
					"inset": {
						ClassGroupID: "ring-w-inset",
					},
					"opacity": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "ring-opacity",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "ring-opacity",
							},
						},
					},
					"offset": {
						Validators: []classGroupValidator{
							{
								Fn:           isLength,
								ClassGroupID: "ring-offset-w",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "ring-offset-w",
							},
							{
								Fn:           isAny,
								ClassGroupID: "ring-offset-color",
							},
						},
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isLength,
						ClassGroupID: "ring-w",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "ring-w",
					},
					{
						Fn:           isAny,
						ClassGroupID: "ring-color",
					},
				},
				ClassGroupID: "ring-w",
			},

			"shadow": {
				NextPart: map[string]classPart{
					"inner": {
						ClassGroupID: "shadow",
					},
					"none": {
						ClassGroupID: "shadow",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isTshirtSize,
						ClassGroupID: "shadow",
					},
					{
						Fn:           isArbitraryShadow,
						ClassGroupID: "shadow",
					},
					{
						Fn:           isAny,
						ClassGroupID: "shadow-color",
					},
				},
				ClassGroupID: "shadow",
			},
			"opacity": {
				NextPart: map[string]classPart{},
				Validators: []classGroupValidator{
					{
						Fn:           isNumber,
						ClassGroupID: "opacity",
					},
					{
						Fn:           isArbitraryNumber,
						ClassGroupID: "opacity",
					},
				},
			},

			"mix": {
				NextPart: map[string]classPart{
					"blend": {
						NextPart: map[string]classPart{
							"normal": {
								ClassGroupID: "mix-blend",
							},
							"multiply": {
								ClassGroupID: "mix-blend",
							},
							"screen": {
								ClassGroupID: "mix-blend",
							},
							"overlay": {
								ClassGroupID: "mix-blend",
							},
							"darken": {
								ClassGroupID: "mix-blend",
							},
							"lighten": {
								ClassGroupID: "mix-blend",
							},
							"color": {
								NextPart: map[string]classPart{
									"dodge": {
										ClassGroupID: "mix-blend",
									},
									"burn": {
										ClassGroupID: "mix-blend",
									},
								},
								ClassGroupID: "mix-blend",
							},
							"hard": {
								NextPart: map[string]classPart{
									"light": {
										ClassGroupID: "mix-blend",
									},
								},
							},
							"soft": {
								NextPart: map[string]classPart{
									"light": {
										ClassGroupID: "mix-blend",
									},
								},
							},
							"difference": {
								ClassGroupID: "mix-blend",
							},
							"exclusion": {
								ClassGroupID: "mix-blend",
							},
							"hue": {
								ClassGroupID: "mix-blend",
							},
							"saturation": {
								ClassGroupID: "mix-blend",
							},
							"luminosity": {
								ClassGroupID: "mix-blend",
							},
							"plus": {
								NextPart: map[string]classPart{
									"lighter": {
										ClassGroupID: "mix-blend",
									},
								},
							},
						},
					},
				},
			},
			"filter": {
				NextPart: map[string]classPart{
					"none": {
						ClassGroupID: "filter",
					},
				},
				ClassGroupID: "filter",
			},
			"blur": {
				NextPart: map[string]classPart{
					"none": {
						ClassGroupID: "blur",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isTshirtSize,
						ClassGroupID: "blur",
					},
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "blur",
					},
				},
				ClassGroupID: "blur",
			},

			"brightness": {
				NextPart: map[string]classPart{},
				Validators: []classGroupValidator{
					{
						Fn:           isNumber,
						ClassGroupID: "brightness",
					},
					{
						Fn:           isArbitraryNumber,
						ClassGroupID: "brightness",
					},
				},
				ClassGroupID: "brightness",
			},

			"contrast": {
				NextPart: map[string]classPart{},
				Validators: []classGroupValidator{
					{
						Fn:           isNumber,
						ClassGroupID: "contrast",
					},
					{
						Fn:           isArbitraryNumber,
						ClassGroupID: "contrast",
					},
				},
				ClassGroupID: "contrast",
			},

			"drop": {
				NextPart: map[string]classPart{
					"shadow": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "drop-shadow",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "drop-shadow",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "drop-shadow",
							},
						},
						ClassGroupID: "drop-shadow",
					},
				},
			},

			"grayscale": {
				NextPart: map[string]classPart{
					"0": {
						ClassGroupID: "grayscale",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "grayscale",
					},
				},
				ClassGroupID: "grayscale",
			},

			"hue": {
				NextPart: map[string]classPart{
					"rotate": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "hue-rotate",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "hue-rotate",
							},
						},
					},
				},
			},
			"invert": {
				NextPart: map[string]classPart{
					"0": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "invert",
							},
						},
						ClassGroupID: "invert",
					},
				},
				ClassGroupID: "invert",
			},

			"saturate": {
				NextPart: map[string]classPart{},
				Validators: []classGroupValidator{
					{
						Fn:           isNumber,
						ClassGroupID: "saturate",
					},
					{
						Fn:           isArbitraryNumber,
						ClassGroupID: "saturate",
					},
				},
				ClassGroupID: "saturate",
			},

			"sepia": {
				NextPart: map[string]classPart{
					"0": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "sepia",
							},
						},
						ClassGroupID: "sepia",
					},
				},
				ClassGroupID: "sepia",
			},

			"backdrop": {
				NextPart: map[string]classPart{
					"filter": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "backdrop-filter",
							},
						},
						ClassGroupID: "backdrop-filter",
					},
					"blur": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "backdrop-blur",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isTshirtSize,
								ClassGroupID: "backdrop-blur",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "backdrop-blur",
							},
						},
						ClassGroupID: "backdrop-blur",
					},
					"brightness": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "backdrop-brightness",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "backdrop-brightness",
							},
						},
						ClassGroupID: "backdrop-brightness",
					},
					"contrast": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "backdrop-contrast",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "backdrop-contrast",
							},
						},
						ClassGroupID: "backdrop-contrast",
					},
					"grayscale": {
						NextPart: map[string]classPart{
							"0": {
								ClassGroupID: "backdrop-grayscale",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "backdrop-grayscale",
							},
						},
						ClassGroupID: "backdrop-grayscale",
					},
					"hue": {
						NextPart: map[string]classPart{
							"rotate": {
								Validators: []classGroupValidator{
									{
										Fn:           isNumber,
										ClassGroupID: "backdrop-hue-rotate",
									},
									{
										Fn:           isArbitraryValue,
										ClassGroupID: "backdrop-hue-rotate",
									},
								},
							},
						},
						Validators: []classGroupValidator{},
					},
					"invert": {
						NextPart: map[string]classPart{
							"0": {
								Validators: []classGroupValidator{
									{
										Fn:           isArbitraryValue,
										ClassGroupID: "backdrop-invert",
									},
								},
								ClassGroupID: "backdrop-invert",
							},
						},
						Validators:   []classGroupValidator{},
						ClassGroupID: "backdrop-invert",
					},
					"opacity": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "backdrop-opacity",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "backdrop-opacity",
							},
						},
						ClassGroupID: "backdrop-opacity",
					},
					"saturate": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "backdrop-saturate",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "backdrop-saturate",
							},
						},
						ClassGroupID: "backdrop-saturate",
					},
					"sepia": {
						NextPart: map[string]classPart{
							"0": {
								Validators: []classGroupValidator{
									{
										Fn:           isArbitraryValue,
										ClassGroupID: "backdrop-sepia",
									},
								},
								ClassGroupID: "backdrop-sepia",
							},
						},
						ClassGroupID: "backdrop-sepia",
					},
				},
			},
			"caption": {
				NextPart: map[string]classPart{
					"top": {
						ClassGroupID: "caption",
					},
					"bottom": {
						ClassGroupID: "caption",
					},
				},
				Validators: []classGroupValidator{},
			},
			"transition": {
				NextPart: map[string]classPart{
					"none": {
						ClassGroupID: "transition",
					},
					"all": {
						ClassGroupID: "transition",
					},
					"colors": {
						ClassGroupID: "transition",
					},
					"opacity": {
						ClassGroupID: "transition",
					},
					"shadow": {
						ClassGroupID: "transition",
					},
					"transform": {
						ClassGroupID: "transition",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "transition",
					},
				},
				ClassGroupID: "transition",
			},
			"duration": {
				Validators: []classGroupValidator{
					{
						Fn:           isNumber,
						ClassGroupID: "duration",
					},
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "duration",
					},
				},
			},
			"ease": {
				NextPart: map[string]classPart{
					"linear": {
						ClassGroupID: "ease",
					},
					"in": {
						NextPart: map[string]classPart{
							"out": {
								ClassGroupID: "ease",
							},
						},
						ClassGroupID: "ease",
					},
					"out": {
						ClassGroupID: "ease",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "ease",
					},
				},
			},
			"delay": {
				Validators: []classGroupValidator{
					{
						Fn:           isNumber,
						ClassGroupID: "delay",
					},
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "delay",
					},
				},
			},
			"animate": {
				NextPart: map[string]classPart{
					"none": {
						ClassGroupID: "animate",
					},
					"spin": {
						ClassGroupID: "animate",
					},
					"ping": {
						ClassGroupID: "animate",
					},
					"pulse": {
						ClassGroupID: "animate",
					},
					"bounce": {
						ClassGroupID: "animate",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "animate",
					},
				},
			},
			"transform": {
				NextPart: map[string]classPart{
					"gpu": {
						ClassGroupID: "transform",
					},
					"none": {
						ClassGroupID: "transform",
					},
				},
				ClassGroupID: "transform",
			},
			"scale": {
				NextPart: map[string]classPart{
					"x": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "scale-x",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "scale-x",
							},
						},
					},
					"y": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "scale-y",
							},
							{
								Fn:           isArbitraryNumber,
								ClassGroupID: "scale-y",
							},
						},
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isNumber,
						ClassGroupID: "scale",
					},
					{
						Fn:           isArbitraryNumber,
						ClassGroupID: "scale",
					},
				},
			},
			"rotate": {
				Validators: []classGroupValidator{
					{
						Fn:           isInteger,
						ClassGroupID: "rotate",
					},
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "rotate",
					},
				},
			},
			"translate": {
				NextPart: map[string]classPart{
					"x": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "translate-x",
							},
							{
								Fn:           isLength,
								ClassGroupID: "translate-x",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "translate-x",
							},
						},
					},
					"y": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "translate-y",
							},
							{
								Fn:           isLength,
								ClassGroupID: "translate-y",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "translate-y",
							},
						},
					},
				},
			},
			"skew": {
				NextPart: map[string]classPart{
					"x": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "skew-x",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "skew-x",
							},
						},
					},
					"y": {
						Validators: []classGroupValidator{
							{
								Fn:           isNumber,
								ClassGroupID: "skew-y",
							},
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "skew-y",
							},
						},
					},
				},
			},
			"origin": {
				NextPart: map[string]classPart{
					"center": {
						ClassGroupID: "transform-origin",
					},
					"top": {
						NextPart: map[string]classPart{
							"right": {
								ClassGroupID: "transform-origin",
							},
							"left": {
								ClassGroupID: "transform-origin",
							},
						},
						ClassGroupID: "transform-origin",
					},
					"right": {
						ClassGroupID: "transform-origin",
					},
					"bottom": {
						NextPart: map[string]classPart{
							"right": {
								ClassGroupID: "transform-origin",
							},
							"left": {
								ClassGroupID: "transform-origin",
							},
						},
						ClassGroupID: "transform-origin",
					},
					"left": {
						ClassGroupID: "transform-origin",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "transform-origin",
					},
				},
			},
			"accent": {
				NextPart: map[string]classPart{
					"auto": {
						NextPart:     map[string]classPart{},
						Validators:   []classGroupValidator{},
						ClassGroupID: "accent",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isAny,
						ClassGroupID: "accent",
					},
				},
			},
			"appearance": {
				NextPart: map[string]classPart{
					"none": {
						ClassGroupID: "appearance",
					},
					"auto": {
						ClassGroupID: "appearance",
					},
				},
			},
			"cursor": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "cursor",
					},
					"default": {
						ClassGroupID: "cursor",
					},
					"pointer": {
						ClassGroupID: "cursor",
					},
					"wait": {
						ClassGroupID: "cursor",
					},
					"text": {
						ClassGroupID: "cursor",
					},
					"move": {
						ClassGroupID: "cursor",
					},
					"help": {
						ClassGroupID: "cursor",
					},
					"not": {
						NextPart: map[string]classPart{
							"allowed": {
								ClassGroupID: "cursor",
							},
						},
					},
					"none": {
						ClassGroupID: "cursor",
					},
					"context": {
						NextPart: map[string]classPart{
							"menu": {
								ClassGroupID: "cursor",
							},
						},
					},
					"progress": {
						ClassGroupID: "cursor",
					},
					"cell": {
						ClassGroupID: "cursor",
					},
					"crosshair": {
						ClassGroupID: "cursor",
					},
					"vertical": {
						NextPart: map[string]classPart{
							"text": {
								ClassGroupID: "cursor",
							},
						},
					},
					"alias": {
						ClassGroupID: "cursor",
					},
					"copy": {
						ClassGroupID: "cursor",
					},
					"no": {
						NextPart: map[string]classPart{
							"drop": {
								ClassGroupID: "cursor",
							},
						},
					},
					"grab": {
						ClassGroupID: "cursor",
					},
					"grabbing": {
						ClassGroupID: "cursor",
					},
					"all": {
						NextPart: map[string]classPart{
							"scroll": {
								ClassGroupID: "cursor",
							},
						},
					},
					"col": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"row": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"n": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"e": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"s": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"w": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"ne": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"nw": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"se": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"sw": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"ew": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"ns": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"nesw": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"nwse": {
						NextPart: map[string]classPart{
							"resize": {
								ClassGroupID: "cursor",
							},
						},
					},
					"zoom": {
						NextPart: map[string]classPart{
							"in": {
								ClassGroupID: "cursor",
							},
							"out": {
								ClassGroupID: "cursor",
							},
						},
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isArbitraryValue,
						ClassGroupID: "cursor",
					},
				},
			},
			"caret": {
				NextPart: map[string]classPart{},
				Validators: []classGroupValidator{
					{
						Fn:           isAny,
						ClassGroupID: "caret-color",
					},
				},
			},
			"pointer": {
				NextPart: map[string]classPart{
					"events": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "pointer-events",
							},
							"auto": {
								ClassGroupID: "pointer-events",
							},
						},
					},
				},
			},
			"resize": {
				NextPart: map[string]classPart{
					"none": {
						ClassGroupID: "resize",
					},
					"y": {
						ClassGroupID: "resize",
					},
					"x": {
						ClassGroupID: "resize",
					},
				},
				ClassGroupID: "resize",
			},
			"scroll": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "scroll-behavior",
					},
					"smooth": {
						ClassGroupID: "scroll-behavior",
					},
					"m": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-m",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-m",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-m",
							},
						},
					},
					"mx": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-mx",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-mx",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-mx",
							},
						},
					},
					"my": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-my",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-my",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-my",
							},
						},
					},
					"ms": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-ms",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-ms",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-ms",
							},
						},
					},
					"me": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-me",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-me",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-me",
							},
						},
					},
					"mt": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-mt",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-mt",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-mt",
							},
						},
					},
					"mr": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-mr",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-mr",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-mr",
							},
						},
					},
					"mb": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-mb",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-mb",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-mb",
							},
						},
					},
					"ml": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-ml",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-ml",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-ml",
							},
						},
					},
					"p": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-p",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-p",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-p",
							},
						},
					},
					"px": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-px",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-px",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-px",
							},
						},
					},
					"py": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-py",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-py",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-py",
							},
						},
					},
					"ps": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-ps",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-ps",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-ps",
							},
						},
					},
					"pe": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-pe",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-pe",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-pe",
							},
						},
					},
					"pt": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-pt",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-pt",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-pt",
							},
						},
					},
					"pr": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-pr",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-pr",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-pr",
							},
						},
					},
					"pb": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-pb",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-pb",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-pb",
							},
						},
					},
					"pl": {
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "scroll-pl",
							},
							{
								Fn:           isLength,
								ClassGroupID: "scroll-pl",
							},
							{
								Fn:           isArbitraryLength,
								ClassGroupID: "scroll-pl",
							},
						},
					},
				},
			},
			"snap": {
				NextPart: map[string]classPart{
					"start": {
						ClassGroupID: "snap-align",
					},
					"end": {
						ClassGroupID: "snap-align",
					},
					"center": {
						ClassGroupID: "snap-align",
					},
					"align": {
						NextPart: map[string]classPart{
							"none": {
								ClassGroupID: "snap-align",
							},
						},
					},
					"normal": {
						ClassGroupID: "snap-stop",
					},
					"always": {
						ClassGroupID: "snap-stop",
					},
					"none": {
						ClassGroupID: "snap-type",
					},
					"x": {
						ClassGroupID: "snap-type",
					},
					"y": {
						ClassGroupID: "snap-type",
					},
					"both": {
						ClassGroupID: "snap-type",
					},
					"mandatory": {
						ClassGroupID: "snap-strictness",
					},
					"proximity": {
						ClassGroupID: "snap-strictness",
					},
				},
			},
			"touch": {
				NextPart: map[string]classPart{
					"auto": {
						ClassGroupID: "touch",
					},
					"none": {
						ClassGroupID: "touch",
					},
					"manipulation": {
						ClassGroupID: "touch",
					},
					"pan": {
						NextPart: map[string]classPart{
							"x": {
								ClassGroupID: "touch-x",
							},
							"left": {
								ClassGroupID: "touch-x",
							},
							"right": {
								ClassGroupID: "touch-x",
							},
							"y": {
								ClassGroupID: "touch-y",
							},
							"up": {
								ClassGroupID: "touch-y",
							},
							"down": {
								ClassGroupID: "touch-y",
							},
						},
					},
					"pinch": {
						NextPart: map[string]classPart{
							"zoom": {
								ClassGroupID: "touch-pz",
							},
						},
					},
				},
			},
			"select": {
				NextPart: map[string]classPart{
					"none": {
						ClassGroupID: "select",
					},
					"text": {
						ClassGroupID: "select",
					},
					"all": {
						ClassGroupID: "select",
					},
					"auto": {
						ClassGroupID: "select",
					},
				},
				Validators: []classGroupValidator{},
			},
			"will": {
				NextPart: map[string]classPart{
					"change": {
						NextPart: map[string]classPart{
							"auto": {
								ClassGroupID: "will-change",
							},
							"scroll": {
								ClassGroupID: "will-change",
							},
							"contents": {
								ClassGroupID: "will-change",
							},
							"transform": {
								ClassGroupID: "will-change",
							},
						},
						Validators: []classGroupValidator{
							{
								Fn:           isArbitraryValue,
								ClassGroupID: "will-change",
							},
						},
					},
				},
			},
			"fill": {
				NextPart: map[string]classPart{
					"none": {
						ClassGroupID: "fill",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isAny,
						ClassGroupID: "fill",
					},
				},
			},
			"stroke": {
				NextPart: map[string]classPart{
					"none": {
						ClassGroupID: "stroke",
					},
				},
				Validators: []classGroupValidator{
					{
						Fn:           isLength,
						ClassGroupID: "stroke-w",
					},
					{
						Fn:           isArbitraryLength,
						ClassGroupID: "stroke-w",
					},
					{
						Fn:           isArbitraryNumber,
						ClassGroupID: "stroke-w",
					},
					{
						Fn:           isAny,
						ClassGroupID: "stroke",
					},
				},
			},
			"sr": {
				NextPart: map[string]classPart{
					"only": {
						ClassGroupID: "sr",
					},
				},
			},
			"forced": {
				NextPart: map[string]classPart{"color": {
					NextPart: map[string]classPart{"adjust": {
						NextPart: map[string]classPart{
							"auto": {
								ClassGroupID: "forced-color-adjust",
							},
							"none": {
								ClassGroupID: "forced-color-adjust",
							},
						},
					},
					},
				},
				},
			},
		},
	},
}
