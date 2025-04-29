package utils

import (
	"regexp"
	"strings"
)

func WildCardToRegexp(pattern string) string {
	components := regexp.MustCompile("[*%]+").Split(pattern, -1)
	if len(components) == 1 {
		// if len is 1, there are no *'s, return exact match pattern
		return "^" + "(?i)" + pattern + "$"
	}
	var sb strings.Builder
	for i, literal := range components {

		// Replace char with .*
		if i > 0 {
			sb.WriteString(".*")
		}

		// Quote any regular expression meta characters in the literal text.
		sb.WriteString(regexp.QuoteMeta(literal))
	}
	return "^" + "(?i)" + sb.String() + "$"
}

// WildCardMatch checks if a given value matches a wildcard pattern.
//
// pattern: the wildcard pattern to match against.
// value: the value to check.
// bool: true if the value matches the pattern, false otherwise.
func WildCardMatch(pattern string, value string) bool {
	result, _ := regexp.MatchString(WildCardToRegexp(strings.TrimSpace(pattern)), value)
	return result
}

func WildCardMatchs(pattern []string, value string) bool {
	for i := 0; i < len(pattern); i++ {
		if WildCardMatch(pattern[i], value) {
			return true
		}
	}
	return false
}
