package yml

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
)

func ReadYml(in string) (*model.DataDef, error) {
	if strings.Contains(in, "://") {
		return readFromDb(in)
	}
	return readFromFile(in)
}

func ReadPatternYml(ifile, schemaPattern, tablePattern, columnPattern string) (*model.DataDef, error) {
	dataDef, err := ReadYml(ifile)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return PatternDataDef(dataDef, schemaPattern, tablePattern, columnPattern)
}

func PatternDataDef(dataDef *model.DataDef, schemaPattern, tablePattern, columnPattern string) (*model.DataDef, error) {
	data, err := model.FilterData(dataDef, schemaPattern, tablePattern, columnPattern)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return data, nil
}

// Expand fixed columns to each tables.
func expandFixColumns(dataDef *model.DataDef) {
	// copy fixed columns to each tables
	for i := 0; i < len(dataDef.Schemas); i++ {
		schema := dataDef.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			schema.Tables[j].Columns = append(schema.Tables[j].Columns, dataDef.Fixed...)
		}
	}
	// clean the fixed column
	dataDef.Fixed = []model.Column{}
}

func DumpYml(dataDef *model.DataDef, outfile string, schemaPattern, tablePattern string) error {
	patternDataDef, err := PatternDataDef(dataDef, schemaPattern, tablePattern, "")
	if err != nil {
		return tracerr.Wrap(err)
	}

	// modify the columns to flow style
	schemas := &patternDataDef.Schemas
	for i := 0; i < len(*schemas); i++ {
		tables := &(*schemas)[i].Tables
		for j := 0; j < len(*tables); j++ {
			(*tables)[j].OutColumns = outColumns(&(*tables)[j].Columns)
			(*tables)[j].Columns = nil
			(*tables)[j].OutReferences = outReferences(&(*tables)[j].References)
			(*tables)[j].References = nil
		}
	}

	bytes, err := yamlOut.Marshal(patternDataDef)
	if err != nil {
		return tracerr.Wrap(err)
	}
	output := string(bytes)
	// correct the names
	output = strings.ReplaceAll(output, "out_columns", "columns")
	output = strings.ReplaceAll(output, "_column: ", "")
	output = strings.ReplaceAll(output, "out_references", "references")
	output = strings.ReplaceAll(output, "_reference: ", "")
	fmt.Println(output)
	return nil
}

// WriteYml writes data to yml file with pattern
func WriteYml(dataDef *model.DataDef, outfile string, schemaPattern, tablePattern string) error {
	restoreFixColumns(dataDef)

	patternDataDef, err := PatternDataDef(dataDef, schemaPattern, tablePattern, "")
	if err != nil {
		return tracerr.Wrap(err)
	}

	// modify the columns in fixed to flow style
	(*patternDataDef).OutFixed = outColumns(&patternDataDef.Fixed)
	patternDataDef.Fixed = nil

	// modify the columns to flow style
	schemas := &patternDataDef.Schemas
	for i := 0; i < len(*schemas); i++ {
		tables := &(*schemas)[i].Tables
		for j := 0; j < len(*tables); j++ {
			// references is a runtime content, should not show in output
			(*tables)[j].References = nil
			(*tables)[j].OutColumns = outColumns(&(*tables)[j].Columns)
			(*tables)[j].Columns = nil
		}
	}

	bytes, err := yamlOut.Marshal(patternDataDef)
	if err != nil {
		return tracerr.Wrap(err)
	}

	output := string(bytes)
	// correct the names
	output = strings.ReplaceAll(output, "out_fixed", "fixed")
	output = strings.ReplaceAll(output, "out_columns", "columns")
	output = strings.ReplaceAll(output, "_column: ", "")
	// remove the quote for boolean
	output = strings.ReplaceAll(output, "\"N\"", "N")
	output = strings.ReplaceAll(output, "\"n\"", "n")
	output = strings.ReplaceAll(output, "\"Y\"", "Y")
	output = strings.ReplaceAll(output, "\"y\"", "y")

	if outfile == "" || outfile == "stdout" {
		fmt.Println(output)
	} else {
		if err := os.WriteFile(outfile, []byte(output), fs.FileMode(0o744)); err != nil {
			return tracerr.Wrap(err)
		}
	}
	return nil
}

func outColumns(columns *[]model.Column) []model.OutColumn {
	// lowercase the column type
	for k := 0; k < len(*columns); k++ {
		(*columns)[k].DataType = strings.ToLower((*columns)[k].DataType)
	}

	outColumns := make([]model.OutColumn, len(*columns))
	for k, column := range *columns {
		outColumns[k].Value = column
	}
	return outColumns
}

func outReferences(references *[]model.Reference) []model.OutReference {
	outReferences := make([]model.OutReference, len(*references))
	for k, reference := range *references {
		outReferences[k].Value = reference
	}
	return outReferences
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
