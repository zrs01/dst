package fileloader

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/sanity-io/litter"
	"github.com/zrs01/dst/internal/dbm"
	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/utils"
	"github.com/ztrue/tracerr"
	yamlIn "gopkg.in/yaml.v3"
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
		f, err := guessSchemaFileName()
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		input = f
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
	dataDef, err = utils.Filter(dataDef, schemaPattern, tablePattern, columnPattern)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return dataDef, nil
}

func loadFromYml(file string) (*model.DataDef, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var d model.DataDef
	if err := yamlIn.Unmarshal(yamlFile, &d); err != nil {
		return nil, tracerr.Wrap(err)
	}

	updateReferenceTables(&d)
	expandFixColumns(&d)

	// validate	data
	validateResult := model.Verify(&d)
	if len(validateResult) > 0 {
		lo.ForEach(validateResult, func(v string, _ int) {
			fmt.Println(v)
		})
		return nil, tracerr.Errorf("invalid data")
	}
	return &d, err
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
		return nil, tracerr.Errorf("failed to parse %s", dataSourceName)
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
// If no matching file is found, it returns empty string.
// If multiple matching files are found, it returns an error.
//
// Returns:
// - string: the name of the schema file.
// - error: an error if there are multiple schema files found, or if there is an error during the file search.
func guessSchemaFileName() (string, error) {
	fileName := "schema.yml"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// if schema.yml at current folder, use it as default
		files, err := filepath.Glob("*schema*.yml")
		if err != nil {
			return "", nil
		}
		if len(files) > 0 {
			return "", tracerr.Errorf("multiple schema files found: %v", files)
		}
	}
	return fileName, nil
}

func updateReferenceTables(dataDef *model.DataDef) {
	// the map table to speed up the lookup process
	tableMap := make(map[string]*model.Table)
	for i := 0; i < len(dataDef.Schemas); i++ {
		schema := dataDef.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			table := &schema.Tables[j]
			tableMap[table.Name] = table
		}
	}

	// update the reference table
	for i := 0; i < len(dataDef.Schemas); i++ {
		schema := dataDef.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			table := &schema.Tables[j]
			for k := 0; k < len(table.Columns); k++ {
				column := table.Columns[k]
				if column.ForeignKey != "" {
					fkTableName, fkColumnName, found := strings.Cut(column.ForeignKey, ".")
					if found {
						fkTable, ok := tableMap[fkTableName]
						if ok {
							fkTable.References = append(fkTable.References, model.Reference{
								ColumnName: fkColumnName,
								Foreign:    []*model.ForeignTable{{Table: table.Name, Column: column.Name}},
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
	dataDef.Fixed = []*model.Column{}
}

func DumpYml(dataDef *model.DataDef, outfile string, schemaPattern, tablePattern string) error {
	patternDataDef, err := utils.Filter(dataDef, schemaPattern, tablePattern, "")
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
