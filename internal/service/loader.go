package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
	yamlIn "gopkg.in/yaml.v3"
)

type LoadBuilder struct {
	options struct {
		// schemaPattern         string
		tablePattern          string
		columnPattern         string
		updateReferenceTables bool
		expandFixColumns      bool
	}
}

// func WithSchemaPattern(pattern string) func(*LoadBuilder) {
// 	return func(s *LoadBuilder) {
// 		s.options.schemaPattern = pattern
// 	}
// }

func WithTablePattern(pattern string) func(*LoadBuilder) {
	return func(s *LoadBuilder) {
		s.options.tablePattern = pattern
	}
}

func WithColumnPattern(pattern string) func(*LoadBuilder) {
	return func(s *LoadBuilder) {
		s.options.columnPattern = pattern
	}
}

func WithUpdateReferenceTables() func(*LoadBuilder) {
	return func(s *LoadBuilder) {
		s.options.updateReferenceTables = true
	}
}

func WithExpandFixColumns() func(*LoadBuilder) {
	return func(s *LoadBuilder) {
		s.options.expandFixColumns = true
	}
}

func NewLoadBuilder(opts ...func(*LoadBuilder)) *LoadBuilder {
	builder := &LoadBuilder{}
	for _, o := range opts {
		o(builder)
	}
	return builder
}

func (s *LoadBuilder) LoadFromFile(file string) (*model.Schema, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var d model.Schema
	if err := yamlIn.Unmarshal(yamlFile, &d); err != nil {
		return nil, tracerr.Wrap(err)
	}

	if s.options.updateReferenceTables {
		s.updateReferenceTables(&d)
	}
	if s.options.expandFixColumns {
		s.expandFixColumns(&d)
	}

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

// func (s *LoadBuilder) LoadFromDB(dsn string, ccf string) (*model.DataDef, error) {
// 	var dataDef *model.DataDef
// 	var err error
// 	if strings.HasPrefix(dsn, MYSQL_PREFIX) {
// 		dataDef, err = s.loadFromMysql(dsn)
// 		if err != nil {
// 			return nil, tracerr.Wrap(err)
// 		}
// 	} else if strings.HasPrefix(dsn, MSSQL_PREFIX) {
// 		dataDef, err = s.loadFromMssql(dsn)
// 		if err != nil {
// 			return nil, tracerr.Wrap(err)
// 		}
// 	} else if strings.HasPrefix(dsn, MARIADB_PREFIX) {
// 		dataDef, err = s.loadFromMariadb(dsn, ccf)
// 		if err != nil {
// 			return nil, tracerr.Wrap(err)
// 		}
// 	} else {
// 		return nil, tracerr.New("Unsupported database type")
// 	}

// 	// if s.options.schemaPattern != "" || s.options.tablePattern != "" || s.options.columnPattern != "" {
// 	// 	dataDef, err = util.FilterData(dataDef, s.options.schemaPattern, s.options.tablePattern, s.options.columnPattern)
// 	// 	if err != nil {
// 	// 		return nil, tracerr.Wrap(err)
// 	// 	}
// 	// }
// 	return dataDef, nil
// }

// // loadFromMysql loads the data from a MySQL database specified by the dataSourceName.
// func (s *LoadBuilder) loadFromMysql(dsn string) (*model.DataDef, error) {
// 	// dataSourceName: mysql://[username[:password]@][protocol[(address[:port])]]/dbname[?param1=value1&...&paramN=valueN]
// 	regex := regexp.MustCompile(`\w+\://[^/]*/(\w+)`).FindStringSubmatch(dsn)
// 	if len(regex) == 0 {
// 		return nil, tracerr.Errorf("Failed to parse %s", dsn)
// 	}
// 	schema := regex[1]
// 	dsName := dsn[len(MYSQL_PREFIX):]
// 	service := dbm.NewMysqlService(dsName)
// 	// fmt.Printf("Read from %s\n", dsName)
// 	result, err := service.Read(strings.Replace(schema, "/", "", -1))
// 	if err != nil {
// 		return nil, tracerr.Wrap(err)
// 	}
// 	return result, nil
// }

// // loadFromMssql loads the data from a Microsoft SQL Server database specified by the dataSourceName.
// func (s *LoadBuilder) loadFromMssql(dsn string) (*model.DataDef, error) {
// 	// dataSourceName: sqlserver: //username:password@host[:port][/instance][?param1=value1&...&paramN=valueN]
// 	schema := ""
// 	parts := strings.Split(dsn, "?")
// 	for _, part := range parts {
// 		pair := strings.Split(part, "=")
// 		if (len(pair) == 2) && (pair[0] == "database") {
// 			schema = pair[1]
// 		}
// 	}
// 	if schema == "" {
// 		return nil, tracerr.Errorf("failed to parse %s", dsn)
// 	}
// 	service := dbm.NewMssqlService(dsn)
// 	// fmt.Printf("Read from %s\n", dataSourceName)
// 	result, err := service.Read(strings.Replace(schema, "/", "", -1))
// 	if err != nil {
// 		return nil, tracerr.Wrap(err)
// 	}
// 	return result, nil
// }

// func (s *LoadBuilder) loadFromMariadb(dsn, ccf string) (*model.DataDef, error) {
// 	// dataSourceName: mariadb://[username[:password]@][protocol[(address[:port])]]/dbname[?param1=value1&...&paramN=valueN]
// 	regex := regexp.MustCompile(`\w+\://[^/]*/(\w+)`).FindStringSubmatch(dsn)
// 	if len(regex) == 0 {
// 		return nil, tracerr.Errorf("Failed to parse %s", dsn)
// 	}
// 	schema := regex[1]
// 	dsName := dsn[len(MARIADB_PREFIX):]
// 	dbManager := mariadb.NewExportService().
// 		WithDriverName("mysql").
// 		WithDataSourceName(dsName).
// 		WithDatabaseName(schema).
// 		WithTableName(config.Setting.Text.TableFilter).
// 		WithCommonColumnFile(ccf)
// 	schemaDef, err := dbManager.ToModel()
// 	if err != nil {
// 		return nil, tracerr.Wrap(err)
// 	}
// 	return &model.DataDef{
// 		Schemas: []*model.Schema{schemaDef},
// 	}, nil
// }

func (s *LoadBuilder) updateReferenceTables(schema *model.Schema) {
	// the map table to speed up the lookup process
	tableMap := make(map[string]*model.Table)
	for j := 0; j < len(schema.Tables); j++ {
		table := schema.Tables[j]
		tableMap[table.Name] = table
	}

	// update the reference table
	for j := 0; j < len(schema.Tables); j++ {
		table := schema.Tables[j]
		for k := 0; k < len(table.Columns); k++ {
			column := table.Columns[k]
			if column.ForeignKey != "" {
				if fkTableName, fkColumnName, found := strings.Cut(column.ForeignKey, "."); found {
					if fkTable, ok := tableMap[fkTableName]; ok {
						fkTable.References = append(fkTable.References, &model.Reference{
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

// Expand fixed columns to each tables.
func (s *LoadBuilder) expandFixColumns(schema *model.Schema) {
	// copy fixed columns to each tables
	// for j := 0; j < len(schema.Tables); j++ {
	// 	schema.Tables[j].Columns = append(schema.Tables[j].Columns, dataDef.Fixed...)
	// }
}

func (s *LoadBuilder) Filter(srcSchema *model.Schema) (*model.Schema, error) {
	var isPatternMatch = func(pattern, value string) bool {
		return pattern == "" || util.WildCardMatchWithCommaPattern(pattern, value)
	}

	dstSchema := model.Schema{}
	var dstTables []*model.Table

	// Filter tables that match the table pattern specified in settings
	// Uses wildcard matching to include/exclude tables based on their names
	matchedTables := lo.Filter(srcSchema.Tables, func(t *model.Table, _ int) bool {
		return isPatternMatch(s.options.tablePattern, t.Name)
	})

	// return empty schema if no tables match the filter pattern
	if len(matchedTables) == 0 {
		return &model.Schema{}, tracerr.New(fmt.Sprintf("no tables matched the filter pattern '%s'", s.options.tablePattern))
	}

	for j := 0; j < len(matchedTables); j++ {
		columns := lo.Filter(matchedTables[j].Columns, func(c *model.Column, _ int) bool {
			return isPatternMatch(s.options.columnPattern, c.Name)
		})
		if len(columns) > 0 {
			matchedTables[j].Columns = columns
			dstTables = append(dstTables, matchedTables[j])
		}
	}

	if len(dstTables) > 0 {
		dstSchema.Tables = dstTables
	} else {
		return &model.Schema{}, tracerr.New(fmt.Sprintf("no tables with matching columns found for filter pattern '%s'", s.options.columnPattern))
	}
	return &dstSchema, nil
}
