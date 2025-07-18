package model

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

type Schema struct {
	Name    string   `yaml:"name,omitempty"`
	Desc    string   `yaml:"desc,omitempty"`
	Tables  []*Table `yaml:"tables,omitempty"`
	Indexes []*Index `yaml:"indexes,omitempty"`
}

func (s *Schema) Filter(tablePattern, columnPattern string) error {
	var isPatternMatch = func(pattern, value string) bool {
		return pattern == "" || util.WildCardMatchWithCommaPattern(pattern, value)
	}

	var dstTables []*Table

	// Filter tables that match the table pattern specified in settings
	// Uses wildcard matching to include/exclude tables based on their names
	matchedTables := lo.Filter(s.Tables, func(t *Table, _ int) bool {
		return isPatternMatch(tablePattern, t.Name)
	})

	// return empty schema if no tables match the filter pattern
	if len(matchedTables) == 0 {
		return tracerr.New(fmt.Sprintf("no tables matched the filter pattern '%s'", tablePattern))
	}

	for j := 0; j < len(matchedTables); j++ {
		columns := lo.Filter(matchedTables[j].Columns, func(c *Column, _ int) bool {
			return isPatternMatch(columnPattern, c.Name)
		})
		if len(columns) > 0 {
			matchedTables[j].Columns = columns
			dstTables = append(dstTables, matchedTables[j])
		}
	}

	if len(dstTables) > 0 {
		s.Tables = dstTables
	} else {
		return tracerr.New(fmt.Sprintf("no tables with matching columns found for filter pattern '%s'", columnPattern))
	}
	return nil
}
