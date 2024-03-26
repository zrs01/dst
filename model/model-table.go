package model

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/samber/lo"
)

func (s *Table) PrimaryColumns() []Column {
	cols := lo.Filter(s.Columns, func(c Column, _ int) bool {
		return c.Identity == "Y" || is_numeric(c.Identity)
	})
	sort.Slice(cols, func(i, j int) bool {
		return identity_num(cols[i].Identity) < identity_num(cols[j].Identity)
	})
	return cols
}

func is_numeric(word string) bool {
	return regexp.MustCompile(`\d`).MatchString(word)
}

func identity_num(value string) int {
	if value == "Y" {
		return 0
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return i
}
