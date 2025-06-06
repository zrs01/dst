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
	schema, err := service.NewLoadBuilder().LoadFromFile(config.Setting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	builder := db.NewService().DDLBuilder()
	output, err := builder.CreateDatabase(schema)
	if err != nil {
		return tracerr.Wrap(err)
	}
	fmt.Println(output)
	return nil
}

func GenerateCreateTable() error {
	loadBuilder := service.NewLoadBuilder(
		service.WithTablePattern(config.Setting.TableName),
	)
	schema, err := loadBuilder.LoadFromFile(config.Setting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	schema, err = loadBuilder.Filter(schema)
	if err != nil {
		return tracerr.Wrap(err)
	}
	builder := db.NewService().DDLBuilder()
	for _, table := range schema.Tables {
		fmt.Println(builder.CreateTable(table))
		for _, index := range buildCreateIndexForTable(builder, table) {
			fmt.Println(index)
		}
		for i, constraint := range buildAddConstraintForTable(builder, table) {
			if i == 0 {
				fmt.Println()
			}
			fmt.Println(constraint)
		}
		fmt.Println()
	}
	return nil
}

func GenerateCreateIndex() error {
	loadBuilder := service.NewLoadBuilder(
		service.WithTablePattern(config.Setting.TableName),
	)
	schema, err := loadBuilder.LoadFromFile(config.Setting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	schema, err = loadBuilder.Filter(schema)
	if err != nil {
		return tracerr.Wrap(err)
	}
	builder := db.NewService().DDLBuilder()
	var stmts []string
	for _, table := range schema.Tables {
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

func buildAddConstraintForTable(builder dbcm.DDL, table *model.Table) []string {
	var stmts []string
	for _, column := range table.Columns {
		if column.ForeignKey != "" {
			stmts = append(stmts, builder.AddConstraint(table.Name, column))
		}
	}
	return stmts
}
