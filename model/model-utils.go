package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/zrs01/dst/utils"
	"github.com/ztrue/tracerr"
)

// FilterData filters the data based on the provided schema, table, and column patterns.
func FilterData(data *DataDef, schemaPattern string, tablePattern string, columnPattern string) (*DataDef, error) {
	d := &DataDef{
		Fixed:   data.Fixed,
		Schemas: make([]Schema, 0),
	}

	for i := 0; i < len(data.Schemas); i++ {
		schema := data.Schemas[i]
		if schemaPattern == "" || utils.WildCardMatchs(strings.Split(schemaPattern, ","), schema.Name) {
			var tables []Table
			tbs := lo.Filter(schema.Tables, func(t Table, _ int) bool {
				return tablePattern == "" || utils.WildCardMatchs(strings.Split(tablePattern, ","), t.Name)
			})
			for j := 0; j < len(tbs); j++ {
				columns := lo.Filter(tbs[j].Columns, func(c Column, _ int) bool {
					return columnPattern == "" || utils.WildCardMatchs(strings.Split(columnPattern, ","), c.Name)
				})
				if len(columns) > 0 {
					tbs[j].Columns = columns
					tables = append(tables, tbs[j])
				}
			}

			if len(tables) > 0 {
				schema.Tables = tables
				d.Schemas = append(d.Schemas, schema)
			}
		}
	}
	tables := lo.FlatMap(d.Schemas, func(s Schema, _ int) []Table {
		return s.Tables
	})
	if len(tables) == 0 {
		return nil, tracerr.New("no schema/table/column matched")
	}
	return d, nil
}

// Verify checks the integrity of the data in the DataDef struct.
func Verify(data *DataDef) []string {
	tables := make(map[string][]Column)

	// convert to map for easy searching
	for _, schema := range data.Schemas {
		for _, table := range schema.Tables {
			tables[table.Name] = table.Columns
		}
	}
	tables["fixed"] = data.Fixed

	isFKExist := func(fk string) bool {
		// fk format: table.field
		tf := strings.Split(fk, ".")
		if columns, found := tables[tf[0]]; found {
			for _, column := range columns {
				if column.Name == tf[1] {
					return true
				}
			}
		}
		return false
	}

	result := make([]string, 0)

	for k, columns := range tables {
		for i, column := range columns {
			if column.Name == "" {
				result = append(result, fmt.Sprintf("[Table: %s] missing column name at line %d", k, i))
			}
			if column.DataType == "" {
				result = append(result, fmt.Sprintf("[Table: %s] missing data type of the column '%s' at line %d", k, column.Name, i))
			}
			if column.ForeignKey != "" {
				// check the foreign key whether exists
				if !isFKExist(column.ForeignKey) {
					result = append(result, fmt.Sprintf("[Table: %s] [FK: %s] cannot be found at line %d", k, column.ForeignKey, i))
				}
			}
		}
	}
	return result
}

func isNumeric(word string) bool {
	return regexp.MustCompile(`\d`).MatchString(word)
}

func toInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return i
}
