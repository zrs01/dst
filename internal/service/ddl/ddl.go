package ddl

import (
	"fmt"

	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/db"
	"github.com/zrs01/dst/internal/db/dbcm"
	"github.com/zrs01/dst/internal/service"
	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

func GenerateCreateDatabase() error {
	schemaModel, err := service.NewLoadBuilder().LoadFromFile(config.Setting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	builder := db.NewService().DDLBuilder()
	output, err := builder.CreateDatabase(schemaModel)
	if err != nil {
		return tracerr.Wrap(err)
	}
	fmt.Println(output)
	return nil
}

func GenerateCreateIndex() error {
	loadBuilder := service.NewLoadBuilder(
		service.WithTablePattern(config.Setting.TableName),
	)
	schemaModel, err := loadBuilder.LoadFromFile(config.Setting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	schemaModel, err = loadBuilder.Filter(schemaModel)
	if err != nil {
		return tracerr.Wrap(err)
	}
	builder := db.NewService().DDLBuilder()
	var stmts []string
	for _, table := range schemaModel.Tables {
		stmts = append(stmts, buildCreateIndexForTable(builder, table)...)
	}
	for _, stmt := range stmts {
		fmt.Println(stmt)
	}
	return nil
}

func buildCreateIndexForTable(builder dbcm.DDL, table *model.Table) []string {
	getIndexName := func(tableName, columnName string) string {
		return fmt.Sprintf("idx_%s_%s", tableName, columnName)
	}
	var stmts []string
	for _, column := range table.Columns {
		if util.IsYes(column.Index) {
			indexName := getIndexName(table.Name, column.Name)
			index := builder.CreateIndex(indexName, table.Name, []string{column.Name}, util.IsYes(column.Unique))
			if index != "" {
				stmts = append(stmts, index)
			}
		}
	}
	// additional index
	for i, index := range table.Indexes {
		indexName := index.Name
		if indexName == "" {
			indexName = fmt.Sprintf("idx_%s_%s", table.Name, fmt.Sprintf("%03d", (i+1)))
		}
		stmts = append(stmts, builder.CreateIndex(indexName, table.Name, index.Columns, util.IsYes(index.Unique)))
	}
	return stmts
}
