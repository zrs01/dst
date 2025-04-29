package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Verify checks the integrity of the data in the DataDef struct.
func Verify(data *DataDef) []string {
	tables := make(map[string][]*Column)

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
