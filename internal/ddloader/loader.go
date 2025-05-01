package ddloader

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/db/mariadb"
	"github.com/zrs01/dst/internal/dbm"
	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
	yamlIn "gopkg.in/yaml.v3"
)

type LoadBuilder struct {
	options struct {
		schemaPattern         string
		tablePattern          string
		columnPattern         string
		updateReferenceTables bool
		expandFixColumns      bool
	}
}

const (
	MYSQL_PREFIX   = "mysql://"
	MSSQL_PREFIX   = "sqlserver://"
	MARIADB_PREFIX = "mariadb://"
)

func NewLoadBuilder(opts ...func(*LoadBuilder)) *LoadBuilder {
	builder := &LoadBuilder{}
	for _, o := range opts {
		o(builder)
	}
	return builder
}

// Load loads the data from the specified input.
func (s *LoadBuilder) Load(input string) (*model.DataDef, error) {
	if input == "" {
		return nil, tracerr.Errorf("the input source is empty")
	}

	var dataDef *model.DataDef
	var err error
	if strings.HasPrefix(input, MYSQL_PREFIX) {
		dataDef, err = s.loadFromMysql(input)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
	} else if strings.HasPrefix(input, MSSQL_PREFIX) {
		dataDef, err = s.loadFromMssql(input)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
	} else if strings.HasPrefix(input, MARIADB_PREFIX) {
		dataDef, err = s.loadFromMariadb(input)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
	} else {
		dataDef, err = s.loadFromFile(input)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
	}

	if s.options.schemaPattern != "" || s.options.tablePattern != "" || s.options.columnPattern != "" {
		dataDef, err = util.FilterData(dataDef, s.options.schemaPattern, s.options.tablePattern, s.options.columnPattern)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
	}
	return dataDef, nil
}

func (s *LoadBuilder) loadFromFile(file string) (*model.DataDef, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var d model.DataDef
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

// loadFromMysql loads the data from a MySQL database specified by the dataSourceName.
func (s *LoadBuilder) loadFromMysql(dataSourceName string) (*model.DataDef, error) {
	// dataSourceName: mysql://[username[:password]@][protocol[(address[:port])]]/dbname[?param1=value1&...&paramN=valueN]
	regex := regexp.MustCompile(`\w+\://[^/]*/(\w+)`).FindStringSubmatch(dataSourceName)
	if len(regex) == 0 {
		return nil, tracerr.Errorf("Failed to parse %s", dataSourceName)
	}
	schema := regex[1]
	dsName := dataSourceName[len(MYSQL_PREFIX):]
	service := dbm.NewMysqlService(dsName)
	// fmt.Printf("Read from %s\n", dsName)
	result, err := service.Read(strings.Replace(schema, "/", "", -1))
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return result, nil
}

// loadFromMssql loads the data from a Microsoft SQL Server database specified by the dataSourceName.
func (s *LoadBuilder) loadFromMssql(dataSourceName string) (*model.DataDef, error) {
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
	// fmt.Printf("Read from %s\n", dataSourceName)
	result, err := service.Read(strings.Replace(schema, "/", "", -1))
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return result, nil
}

func (s *LoadBuilder) loadFromMariadb(dataSourceName string) (*model.DataDef, error) {
	// dataSourceName: mariadb://[username[:password]@][protocol[(address[:port])]]/dbname[?param1=value1&...&paramN=valueN]
	regex := regexp.MustCompile(`\w+\://[^/]*/(\w+)`).FindStringSubmatch(dataSourceName)
	if len(regex) == 0 {
		return nil, tracerr.Errorf("Failed to parse %s", dataSourceName)
	}
	schema := regex[1]
	dsName := dataSourceName[len(MARIADB_PREFIX):]
	dbManager := mariadb.NewExportManager().
		WithDriverName("mysql").
		WithDataSourceName(dsName).
		WithDatabaseName(schema).
		WithTableName(config.Setting.Text.TableFilter)
	schemaDef, err := dbManager.ToModel()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &model.DataDef{
		Schemas: []*model.Schema{schemaDef},
	}, nil
}

func (s *LoadBuilder) updateReferenceTables(dataDef *model.DataDef) {
	// the map table to speed up the lookup process
	tableMap := make(map[string]*model.Table)
	for i := 0; i < len(dataDef.Schemas); i++ {
		schema := dataDef.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			table := schema.Tables[j]
			tableMap[table.Name] = table
		}
	}

	// update the reference table
	for i := 0; i < len(dataDef.Schemas); i++ {
		schema := dataDef.Schemas[i]
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
}

// Expand fixed columns to each tables.
func (s *LoadBuilder) expandFixColumns(dataDef *model.DataDef) {
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
