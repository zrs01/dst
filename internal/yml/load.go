package yml

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	yamlOut "github.com/goccy/go-yaml"
	"github.com/samber/lo"
	"github.com/sanity-io/litter"
	"github.com/zrs01/dst/internal/dbm"
	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/utils"
	"github.com/ztrue/tracerr"
)

const (
	MYSQL_PREFIX = "mysql://"
	MSSQL_PREFIX = "sqlserver://"
)

// Load loads the data from the specified input.
//
// The input can be a YAML file path, a MySQL database data source name, or a SQL Server database data source name.
// If the input is a MySQL or SQL Server database, the data is loaded from the database.
// If the input is a YAML file path, the data is loaded from the file.
//
// Returns the loaded data and any error encountered.
func Load(input string) (*model.DataDef, error) {
	// Check if the input starts with mysql:// or sqlserver://.
	// If so, load the data from the corresponding database.
	if strings.HasPrefix(input, MYSQL_PREFIX) {
		return loadFromMysql(input)
	} else if strings.HasPrefix(input, MSSQL_PREFIX) {
		return loadFromMssql(input)
	}

	if input == "" {
		input = guessSchemaFileName()
	}

	// If the input does not start with mysql:// or sqlserver://,
	// assume it is a YAML file path and load the data from the file.
	return loadFromYml(input)
}

// LoadWithFilter loads the data from the specified input, filters it based on the provided schema, table, and column patterns, and returns the filtered data.
//
// The input can be a YAML file path, a MySQL database data source name, or a SQL Server database data source name.
// If the input is a MySQL or SQL Server database, the data is loaded from the database and filtered based on the provided schema, table, and column patterns.
// If the input is a YAML file path, the data is loaded from the file and filtered based on the provided schema, table, and column patterns.
//
// The schemaPattern, tablePattern, and columnPattern parameters are used to filter the data.
// The schemaPattern parameter specifies the schemas to include in the filtered data.
// The tablePattern parameter specifies the tables to include in the filtered data.
// The columnPattern parameter specifies the columns to include in the filtered data.
// If the schemaPattern, tablePattern, or columnPattern parameters are empty, all schemas, tables, or columns are included in the filtered data.
//
// Returns the filtered data and any error encountered.
func LoadWithFilter(input string, schemaPattern, tablePattern, columnPattern string) (*model.DataDef, error) {
	dataDef, err := Load(input)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	dataDef, err = Filter(dataDef, schemaPattern, tablePattern, columnPattern)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return dataDef, nil
}

// loadFromMysql loads the data from a MySQL database specified by the dataSourceName.
//
// The dataSourceName should be in the format: mysql://[username[:password]@][protocol[(address[:port])]]/dbname[?param1=value1&...&paramN=valueN].
// The value after '/' in the dataSourceName is used as the schema name.
//
// Returns the loaded data and any error encountered.
func loadFromMysql(dataSourceName string) (*model.DataDef, error) {
	// dataSourceName: mysql://[username[:password]@][protocol[(address[:port])]]/dbname[?param1=value1&...&paramN=valueN]
	regex := regexp.MustCompile(`\w+\://[^/]*/(\w+)`).FindStringSubmatch(dataSourceName)
	if len(regex) == 0 {
		return nil, tracerr.Errorf("Failed to parse %s", dataSourceName)
	}
	schema := regex[1]
	dsName := dataSourceName[len(MYSQL_PREFIX):]
	service := dbm.NewMysqlService(dsName)
	fmt.Printf("Read from %s\n", dsName)
	result, err := service.Read(strings.Replace(schema, "/", "", -1))
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return result, nil
}

// loadFromMssql loads the data from a Microsoft SQL Server database specified by the dataSourceName.
//
// The dataSourceName should be in the format: sqlserver://username:password@host[:port][/instance][?param1=value1&...&paramN=valueN].
// The value of the 'database' parameter in the dataSourceName is used as the schema name.
//
// Returns the loaded data and any error encountered.
func loadFromMssql(dataSourceName string) (*model.DataDef, error) {
	// dataSourceName: sqlserver: //username:password@host[:port][/instance][?param1=value1&...&paramN=valueN]
	schema := ""
	parts := strings.Split(dataSourceName, "?")
	for _, part := range parts {
		pair := strings.Split(part, "=")
		if (len(pair) == 2) && (pair[0] == "database") {
			schema = pair[1]
		}
	}
	if schema == "" {
		return nil, tracerr.Errorf("Failed to parse %s", dataSourceName)
	}
	service := dbm.NewMssqlService(dataSourceName)
	fmt.Printf("Read from %s\n", dataSourceName)
	result, err := service.Read(strings.Replace(schema, "/", "", -1))
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return result, nil
}

// guessSchemaFileName tries to guess the name of the schema file by looking for a file matching the pattern "*schema*.yml" in the current directory.
//
// If no matching file is found, it returns the default schema file name "schema.yml".
//
// Returns:
// - string: the name of the schema file.
func guessSchemaFileName() string {
	fileName := "schema.yml"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// if schema.yml at current folder, use it as default
		files, err := filepath.Glob("*schema*.yml")
		if err != nil {
			return ""
		}
		if len(files) > 0 {
			fileName = files[0]
		}
	}
	return fileName
}

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
		if schemaPattern == "" || utils.WildCardMatchs(strings.Split(schemaPattern, ","), schema.Name) {
			var tables []model.Table
			filteredTables := lo.Filter(schema.Tables, func(t model.Table, _ int) bool {
				return tablePattern == "" || utils.WildCardMatchs(strings.Split(tablePattern, ","), t.Name)
			})
			for j := 0; j < len(filteredTables); j++ {
				columns := lo.Filter(filteredTables[j].Columns, func(c model.Column, _ int) bool {
					return columnPattern == "" || utils.WildCardMatchs(strings.Split(columnPattern, ","), c.Name)
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

func updateReferenceTables(dataDef *model.DataDef) {
	// the map table to speed up the lookup process
	tableMap := make(map[string]*model.Table)
	for i := 0; i < len(dataDef.Schemas); i++ {
		schema := &dataDef.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			table := &schema.Tables[j]
			tableMap[table.Name] = table
		}
	}

	// update the reference table
	for i := 0; i < len(dataDef.Schemas); i++ {
		schema := &dataDef.Schemas[i]
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
	patternDataDef, err := Filter(dataDef, schemaPattern, tablePattern, "")
	if err != nil {
		return tracerr.Wrap(err)
	}
	litter.Config.StripPackageNames = true
	litter.Config.FieldExclusions = regexp.MustCompile(`^(Out.*)$`)
	litter.Dump(patternDataDef)
	// spew.Dump(patternDataDef)

	// // modify the columns to flow style
	// schemas := &patternDataDef.Schemas
	// for i := 0; i < len(*schemas); i++ {
	// 	tables := &(*schemas)[i].Tables
	// 	for j := 0; j < len(*tables); j++ {
	// 		(*tables)[j].OutColumns = outColumns(&(*tables)[j].Columns)
	// 		(*tables)[j].Columns = nil
	// 		(*tables)[j].OutReferences = outReferences(&(*tables)[j].References)
	// 		(*tables)[j].References = nil
	// 	}
	// }

	// bytes, err := yamlOut.Marshal(patternDataDef)
	// if err != nil {
	// 	return tracerr.Wrap(err)
	// }
	// output := string(bytes)
	// // correct the names
	// output = strings.ReplaceAll(output, "out_columns", "columns")
	// output = strings.ReplaceAll(output, "_column: ", "")
	// output = strings.ReplaceAll(output, "out_references", "references")
	// output = strings.ReplaceAll(output, "_reference: ", "")
	// fmt.Println(output)
	return nil
}

// WriteYml writes data to yml file with pattern
func WriteYml(dataDef *model.DataDef, outfile string, schemaPattern, tablePattern string) error {
	restoreFixColumns(dataDef)

	patternDataDef, err := Filter(dataDef, schemaPattern, tablePattern, "")
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
