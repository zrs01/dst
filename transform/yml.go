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

	// update the reference table
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

func SelectedYml(ifile, schema, tablePattern string, columnPattern string) (*model.DataDef, error) {
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
	data, err := model.FilterData(rawData, schema, tablePattern, columnPattern)
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
	schemas := &data.Schemas
	for i := 0; i < len(*schemas); i++ {
		var tables = &(*schemas)[i].Tables
		for j := 0; j < len(*tables); j++ {
			// references is a runtime content, should not show in output
			(*tables)[j].References = nil

			// lowercase the column type
			for k := 0; k < len((*tables)[j].Columns); k++ {
				(*tables)[j].Columns[k].DataType = strings.ToLower((*tables)[j].Columns[k].DataType)
			}

			// move columns to out_columns for flow style output
			(*tables)[j].OutColumns = make([]model.OutColumn, len((*tables)[j].Columns))
			for k, column := range (*tables)[j].Columns {
				(*tables)[j].OutColumns[k].Value = column
				(*tables)[j].Columns = nil
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
	type tb struct {
		tableName  string
		columnName string
	}
	// Map to store column attributes as keys and a list of tables as values
	tbMap := make(map[string][]tb)

	// create a list of tables with the same column attributes
	for i := 0; i < len(data.Schemas); i++ {
		schema := &data.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			table := &schema.Tables[j]
			for k := 0; k < len(table.Columns); k++ {
				column := &table.Columns[k]
				// Generate a unique key based on column attributes
				key := generateColumnKey(*column)
				// Append the current table to the list of tables with the same attributes
				tbMap[key] = append(tbMap[key], tb{tableName: table.Name, columnName: column.Name})
			}
		}
	}

	// Check if all tables have the same column attributes
	for _, items := range tbMap {
		if len(items) == lo.Reduce(data.Schemas, func(acc int, schema model.Schema, _ int) int { return acc + len(schema.Tables) }, 0) {

			// Remove the column from the table
			for _, item := range items {
				for j := 0; j < len(data.Schemas); j++ {
					schema := &data.Schemas[j]
					for k := 0; k < len(schema.Tables); k++ {
						table := &schema.Tables[k]
						if table.Name == item.tableName {
							for l := 0; l < len(table.Columns); l++ {
								if table.Columns[l].Name == item.columnName {

									// add to fixed columns if it is not duplicated
									if !lo.Contains(data.Fixed, table.Columns[l]) {
										data.Fixed = append(data.Fixed, table.Columns[l])
									}

									// remove the column
									table.Columns = append(table.Columns[:l], table.Columns[l+1:]...)
									break
								}
							}
							break
						}
					}
				}
			}
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
