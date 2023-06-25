package xfmr

import (
	"regexp"
	"strings"
)

func wildCardToRegexp(pattern string) string {
	// components := strings.Split(pattern, "*")
	components := regexp.MustCompile("[*,%]+").Split(pattern, -1)
	if len(components) == 1 {
		// if len is 1, there are no *'s, return exact match pattern
		return "^" + pattern + "$"
	}
	var result strings.Builder
	for i, literal := range components {

		// Replace * with .*
		if i > 0 {
			result.WriteString(".*")
		}

		// Quote any regular expression meta characters in the
		// literal text.
		result.WriteString(regexp.QuoteMeta(literal))
	}
	return "^" + result.String() + "$"
}

func wildCardMatch(pattern string, value string) bool {
	result, _ := regexp.MatchString(wildCardToRegexp(pattern), value)
	return result
}
