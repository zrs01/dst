package utils

import (
	"strings"

	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

// Filter filters the data based on the provided schema, table, and column patterns.
//
// Parameters:
// - data: the data to be filtered.
// - schemaPattern: a comma-separated string of schema patterns. Matches if any pattern matches.
// - tablePattern: a comma-separated string of table patterns. Matches if any pattern matches.
// - columnPattern: a comma-separated string of column patterns. Matches if any pattern matches.
//
// Returns:
// - a new DataDef object containing the filtered data.
// - an error if no schema/table/column matched.
func Filter(data *model.DataDef, schemaPattern string, tablePattern string, columnPattern string) (*model.DataDef, error) {
	d := &model.DataDef{
		Fixed:   data.Fixed,
		Schemas: make([]model.Schema, 0),
	}

	for i := 0; i < len(data.Schemas); i++ {
		schema := data.Schemas[i]
		if schemaPattern == "" || WildCardMatchs(strings.Split(schemaPattern, ","), schema.Name) {
			var tables []model.Table
			filteredTables := lo.Filter(schema.Tables, func(t model.Table, _ int) bool {
				return tablePattern == "" || WildCardMatchs(strings.Split(tablePattern, ","), t.Name)
			})
			for j := 0; j < len(filteredTables); j++ {
				columns := lo.Filter(filteredTables[j].Columns, func(c model.Column, _ int) bool {
					return columnPattern == "" || WildCardMatchs(strings.Split(columnPattern, ","), c.Name)
				})
				if len(columns) > 0 {
					filteredTables[j].Columns = columns
					tables = append(tables, filteredTables[j])
				}
			}

			if len(tables) > 0 {
				schema.Tables = tables
				d.Schemas = append(d.Schemas, schema)
			}
		}
	}
	tables := lo.FlatMap(d.Schemas, func(s model.Schema, _ int) []model.Table {
		return s.Tables
	})
	if len(tables) == 0 {
		return nil, tracerr.New("no schema/table/column matched")
	}
	return d, nil
}
