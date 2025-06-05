package util

import (
	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

// FilterData filters the data based on the provided schema, table, and column patterns.
func FilterData(data *model.SchemaDef, schemaPattern string, tablePattern string, columnPattern string) (*model.SchemaDef, error) {
	newData := &model.SchemaDef{
		Fixed:   data.Fixed,
		Schemas: make([]*model.Schema, 0),
	}

	for i := 0; i < len(data.Schemas); i++ {
		schema := data.Schemas[i]
		if schemaPattern == "" || WildCardMatchWithArrayPattern(SplitWithTrim(schemaPattern, ","), schema.Name) {
			var tables []*model.Table

			// filter tables based on tablePattern
			filteredTables := lo.Filter(schema.Tables, func(t *model.Table, _ int) bool {
				return tablePattern == "" || WildCardMatchWithArrayPattern(SplitWithTrim(tablePattern, ","), t.Name)
			})

			// if no tables were filtered, skip this schema
			if len(filteredTables) == 0 {
				continue
			}

			// process each table
			for j := 0; j < len(filteredTables); j++ {
				columns := lo.Filter(filteredTables[j].Columns, func(c *model.Column, _ int) bool {
					return columnPattern == "" || WildCardMatchWithArrayPattern(SplitWithTrim(columnPattern, ","), c.Name)
				})
				if len(columns) > 0 {
					filteredTables[j].Columns = columns
					tables = append(tables, filteredTables[j])
				}
			}

			if len(tables) > 0 {
				schema.Tables = tables
				newData.Schemas = append(newData.Schemas, schema)
			}
		}
	}
	tables := lo.FlatMap(newData.Schemas, func(s *model.Schema, _ int) []*model.Table {
		return s.Tables
	})
	if len(tables) == 0 {
		return nil, tracerr.New("no schema/table/column matched")
	}
	return newData, nil
}
