package model

import (
	"sort"

	"github.com/samber/lo"
)

// PrimaryKeyColumns returns all columns that contain identity or are numeric, sorted by identity.
func (s *Table) PrimaryKeyColumns() []Column {
	// get all columns contains identity
	cols := lo.Filter(s.Columns, func(c Column, _ int) bool {
		return c.Identity == "Y" || isNumeric(c.Identity)
	})
	// sort the the identity columns
	sort.Slice(cols, func(i, j int) bool {
		return toInt(cols[i].Identity) < toInt(cols[j].Identity)
	})
	return cols
}
