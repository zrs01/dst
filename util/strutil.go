package util

import (
	"strings"

	"github.com/samber/lo"
)

func IsYes(value string) bool {
	return strings.Contains(strings.ToUpper(strings.TrimSpace(value)), "Y")
}

func SplitWithTrim(str string, sep string) []string {
	return lo.FilterMap(strings.Split(str, sep), func(x string, _ int) (string, bool) {
		return strings.TrimSpace(x), x != ""
	})
}
