package transform

import (
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"strings"

	yamlOut "github.com/goccy/go-yaml"
	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
	yamlIn "gopkg.in/yaml.v3"
)

func ReadYml(file string) (*model.DataDef, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var d model.DataDef
	if err := yamlIn.Unmarshal(yamlFile, &d); err != nil {
		return nil, tracerr.Wrap(err)
	}

	/* ------------------------- update reference tables ------------------------ */
	// the map table to speed up the lookup process
	var tableMap = make(map[string]*model.Table)
	for i := 0; i < len(d.Schemas); i++ {
		schema := &d.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			table := &schema.Tables[j]
			tableMap[table.Name] = table
		}
	}

	for i := 0; i < len(d.Schemas); i++ {
		schema := &d.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			table := &schema.Tables[j]
			for k := 0; k < len(table.Columns); k++ {
				column := &table.Columns[k]
				if column.ForeignKey != "" {
					fkTableName, fkColumnName, found := strings.Cut(column.ForeignKey, ".")
					if found {
						fkTable, ok := tableMap[fkTableName]
						if ok {
							fkTable.References = append(fkTable.References, model.Reference{
								ColumnName: fkColumnName,
								Foreign:    []model.ForeignTable{{Table: table.Name, Column: column.Name}},
							})
						} else {
							fmt.Printf("failed to find table '%s'", fkTableName)
						}
					}
				}
			}
		}
	}
	return &d, err
}

func FilterYml(ifile, schema, table string, column string) (*model.DataDef, error) {
	rawData, err := ReadYml(ifile)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	// copy fixed columns to each tables
	for i := 0; i < len(rawData.Schemas); i++ {
		schema := rawData.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			schema.Tables[j].Columns = append(schema.Tables[j].Columns, rawData.Fixed...)
		}
	}
	// clean the fixed column
	rawData.Fixed = []model.Column{}

	validateResult := model.Verify(rawData)
	if len(validateResult) > 0 {
		lo.ForEach(validateResult, func(v string, _ int) {
			fmt.Println(v)
		})
		return nil, tracerr.Errorf("invalid data")
	}
	data, err := model.FilterData(rawData, schema, table, column)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return data, nil
}

func WriteYml(data *model.DataDef, outfile string) error {
	// restoreFixColumns(data)

	// modify the columns in fixed to flow style
	pFixed := &data.Fixed
	for i := 0; i < len(*pFixed); i++ {
		(*data).OutFixed = make([]model.OutColumn, len(*pFixed))
		for j, column := range *pFixed {
			(*data).OutFixed[j].Value = column
		}
	}
	data.Fixed = nil

	// modify the columns to flow style
	pSchemas := &data.Schemas
	for i := 0; i < len(*pSchemas); i++ {
		var pTables = &(*pSchemas)[i].Tables
		for j := 0; j < len(*pTables); j++ {
			// references is a runtime content, should not show in output
			(*pTables)[j].References = nil

			// lowercase the column type
			for k := 0; k < len((*pTables)[j].Columns); k++ {
				(*pTables)[j].Columns[k].DataType = strings.ToLower((*pTables)[j].Columns[k].DataType)
			}

			// move columns to out_columns for flow style output
			(*pTables)[j].OutColumns = make([]model.OutColumn, len((*pTables)[j].Columns))
			for k, column := range (*pTables)[j].Columns {
				(*pTables)[j].OutColumns[k].Value = column
				(*pTables)[j].Columns = nil
			}
		}
	}

	bytes, err := yamlOut.Marshal(data)
	if err != nil {
		return tracerr.Wrap(err)
	}

	output := string(bytes)
	// correct the names
	output = strings.ReplaceAll(output, "out_fixed", "fixed")
	output = strings.ReplaceAll(output, "out_columns", "columns")
	output = strings.ReplaceAll(output, "_column_values: ", "")
	// remove the quote for boolean
	output = strings.ReplaceAll(output, "\"N\"", "N")
	output = strings.ReplaceAll(output, "\"n\"", "n")
	output = strings.ReplaceAll(output, "\"Y\"", "Y")
	output = strings.ReplaceAll(output, "\"y\"", "y")

	if outfile == "" || outfile == "stdout" {
		fmt.Println(output)
	} else {
		os.WriteFile(outfile, []byte(output), fs.FileMode(0744))
	}
	return nil
}

func restoreFixColumns(data *model.DataDef) {
	// Map to store column attributes as keys and a list of tables as values
	columnMap := make(map[string][]string)

	// Iterate over schemas
	for _, schema := range data.Schemas {
		// Iterate over tables
		for _, table := range schema.Tables {
			// Iterate over columns
			for _, column := range table.Columns {
				// Generate a unique key based on column attributes
				key := generateColumnKey(column)

				// Append the current table to the list of tables with the same attributes
				columnMap[key] = append(columnMap[key], table.Name)
			}
		}
	}

	// Print tables with the same column attributes
	for key, tables := range columnMap {
		if len(tables) == lo.Reduce(data.Schemas, func(acc int, schema model.Schema, _ int) int { return acc + len(schema.Tables) }, 0) {
			fmt.Printf("Tables with the same column attributes [%s]: %v\n", key, tables)
		}

	}
}

func generateColumnKey(column model.Column) string {
	// Use reflection to get the column attributes
	v := reflect.ValueOf(column)
	numFields := v.NumField()

	// Slice to store attribute values
	attributes := make([]interface{}, numFields)

	fieldNames := []string{"Name", "DataType", "Identity", "NotNull", "Unique", "Value", "ForeignKey", "Cardinality", "Title", "Index", "Compute"}
	for i := 0; i < len(fieldNames); i++ {
		attributes[i] = v.FieldByName(fieldNames[i])
	}

	// Iterate over struct fields
	// for i := 0; i < numFields; i++ {
	// 	fmt.Println(v.Field(i))
	// 	attributes[i] = v.Field(i).Interface()
	// }

	// Format the attributes as a string and return
	return fmt.Sprintf("%v", attributes)
}
