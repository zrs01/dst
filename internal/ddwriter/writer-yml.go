package ddwriter

import (
	"fmt"
	"io/fs"
	"os"
	"reflect"

	// yamlOut "github.com/goccy/go-yaml"
	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
	"gopkg.in/yaml.v3"
)

func restoreFixColumns(data *model.Root) {
	type tb struct {
		tableName  string
		columnName string
	}
	// Map to store column attributes as keys and a list of tables as values
	tbMap := make(map[string][]tb)

	// create a list of tables with the same column attributes
	for i := 0; i < len(data.Schemas); i++ {
		schema := data.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			table := schema.Tables[j]
			for k := 0; k < len(table.Columns); k++ {
				column := table.Columns[k]
				// Generate a unique key based on column attributes
				key := generateColumnKey(*column)
				// Append the current table to the list of tables with the same attributes
				tbMap[key] = append(tbMap[key], tb{tableName: table.Name, columnName: column.Name})
			}
		}
	}

	// Check if all tables have the same column attributes
	for _, items := range tbMap {
		if len(items) == lo.Reduce(data.Schemas, func(acc int, schema *model.Schema, _ int) int { return acc + len(schema.Tables) }, 0) {
			// Remove the column from the table
			for _, item := range items {
				for j := 0; j < len(data.Schemas); j++ {
					schema := data.Schemas[j]
					for k := 0; k < len(schema.Tables); k++ {
						table := schema.Tables[k]
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

func OutputYml(dataDef *model.Schema, outfile string) error {
	bytes, err := yaml.Marshal(dataDef)
	if err != nil {
		return tracerr.Wrap(err)
	}

	if outfile == "" {
		fmt.Println(string(bytes))
	} else {
		if err := os.WriteFile(outfile, bytes, fs.FileMode(0o664)); err != nil {
			return tracerr.Wrap(err)
		}
	}
	return nil
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

// func outReferences(references *[]model.Reference) []model.OutReference {
// 	outReferences := make([]model.OutReference, len(*references))
// 	for k, reference := range *references {
// 		outReferences[k].Value = reference
// 	}
// 	return outReferences
// }
