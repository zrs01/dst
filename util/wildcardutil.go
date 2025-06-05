package util

import (
	"regexp"
	"strings"

	"github.com/samber/lo"
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

func WildCardMatchWithCommaPattern(pattern string, value string) bool {
	return WildCardMatchWithArrayPattern(SplitWithTrim(pattern, ","), value)
}

func WildCardMatchWithArrayPattern(pattern []string, value string) bool {
	return lo.SomeBy(pattern, func(x string) bool {
		return WildCardMatch(x, value)
	})
}

// WildCardMatch checks if a given value matches a wildcard pattern.
func WildCardMatch(pattern string, value string) bool {
	result, _ := regexp.MatchString(WildCardToRegexp(strings.TrimSpace(pattern)), value)
	return result
}
