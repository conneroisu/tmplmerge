package split_benchmark

import (
	"regexp"
	"strings"
	"testing"
)

// Regular expression patterns
var splitPatternSpace = regexp.MustCompile(`\s+`)
var splitPatternComma = regexp.MustCompile(`,\s*`)

// Test strings with various characteristics
var (
	// Simple space-separated string
	simpleString = "a b c d"

	// String with multiple spaces
	multiSpaceString = "a  b   c    d"

	// Long string with consistent spacing
	longString = "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12 word13 word14 word15"

	// String with mixed whitespace (spaces, tabs, newlines)
	mixedWhitespaceString = "word1 word2\tword3\nword4\r\nword5   word6"

	// Comma-separated values with varied spacing
	commaString = "item1, item2,item3,   item4, item5"
)

// Benchmark functions for simple string with single spaces
func BenchmarkRegexSplitSimple(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		splitPatternSpace.Split(simpleString, -1)
	}
}

func BenchmarkStringsSplitSimple(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		strings.Split(simpleString, " ")
	}
}

func BenchmarkStringsFieldsSimple(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		strings.Fields(simpleString)
	}
}

// Benchmark functions for string with multiple spaces
func BenchmarkRegexSplitMultiSpace(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		splitPatternSpace.Split(multiSpaceString, -1)
	}
}

func BenchmarkStringsSplitMultiSpace(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		strings.Split(multiSpaceString, " ")
	}
}

func BenchmarkStringsFieldsMultiSpace(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		strings.Fields(multiSpaceString)
	}
}

// Benchmark functions for long string
func BenchmarkRegexSplitLong(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		splitPatternSpace.Split(longString, -1)
	}
}

func BenchmarkStringsSplitLong(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		strings.Split(longString, " ")
	}
}

func BenchmarkStringsFieldsLong(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		strings.Fields(longString)
	}
}

// Benchmark functions for mixed whitespace string
func BenchmarkRegexSplitMixedWhitespace(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		splitPatternSpace.Split(mixedWhitespaceString, -1)
	}
}

func BenchmarkStringsSplitMixedWhitespace(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		strings.Split(mixedWhitespaceString, " ")
	}
}

func BenchmarkStringsFieldsMixedWhitespace(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		strings.Fields(mixedWhitespaceString)
	}
}

// Benchmark functions for comma-separated values
func BenchmarkRegexSplitComma(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		splitPatternComma.Split(commaString, -1)
	}
}

func BenchmarkStringsSplitComma(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		strings.Split(commaString, ",")
	}
}

// Additional benchmarks to test compiled vs non-compiled regex
func BenchmarkRegexCompileAndSplit(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		re := regexp.MustCompile(`\s+`)
		re.Split(simpleString, -1)
	}
}

// Test with pre-allocated result slice
func BenchmarkStringsSplitPreAllocated(b *testing.B) {
	result := make([]string, 0, 4) // Pre-allocate capacity for expected result size
	b.ResetTimer()
	for b.Loop() {
		result = strings.SplitN(simpleString, " ", 4)
	}
	b.Log(result)
}

// Test with specific limit for comparison
func BenchmarkRegexSplitWithLimit(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		splitPatternSpace.Split(simpleString, 4)
	}
}

func BenchmarkStringsSplitWithLimit(b *testing.B) {

	for b.Loop() {
		strings.SplitN(simpleString, " ", 4)
	}
}
