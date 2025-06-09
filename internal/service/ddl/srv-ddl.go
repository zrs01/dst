package ddl

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/db"
	"github.com/zrs01/dst/internal/db/dbcm"
	"github.com/zrs01/dst/internal/service"
	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

func CreateDatabase() error {
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

func CreateTable() error {
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
		fmt.Println(createTableForTable(builder, table))
	}
	return nil
}

func DropTable() error {
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
		stmts = append(stmts, dropIndexForTable(builder, table)...)
		stmts = append(stmts, dropConstraintForTable(builder, table)...)
		stmts = append(stmts, builder.DropTable(table)...)
	}
	for _, stmt := range stmts {
		fmt.Println(stmt)
	}
	fmt.Println()
	return nil
}

func Diff() error {
	// get model from database
	dbSchema, err := db.NewService().Load(config.Setting.TableName)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// get model from schema definition file
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

	ddlBuilder := db.NewService().DDLBuilder()
	for _, table := range schema.Tables {
		dbTable, found := lo.Find(dbSchema.Tables, func(t *model.Table) bool {
			return t.Name == table.Name
		})
		if !found {
			// new table
			fmt.Println(createTableForTable(ddlBuilder, table))
		} else {
			// compare columns of table
			for _, column := range table.Columns {
				dbColumn, found := lo.Find(dbTable.Columns, func(c *model.Column) bool {
					return c.Name == column.Name
				})
				if !found {
					// new column
					fmt.Println(dbColumn)
					// fmt.Println(ddlBuilder.AddColumn(table, column))
				} else {
					// compare column
					// if dbColumn.DataType != column.DataType || dbColumn.NotNull != column.NotNull || dbColumn.Value != column.Value || dbColumn.Desc != column.Desc || dbColumn.Unique != column.Unique {
					// 	fmt.Println(ddlBuilder.ModifyColumn(table, column))
					// }
				}
			}
		}
	}

	return nil
}

// func GenerateCreateIndex() error {
// 	loadBuilder := service.NewLoadBuilder(
// 		service.WithTablePattern(config.Setting.TableName),
// 	)
// 	schema, err := loadBuilder.LoadFromFile(config.Setting.Sdf)
// 	if err != nil {
// 		return tracerr.Wrap(err)
// 	}
// 	schema, err = loadBuilder.Filter(schema)
// 	if err != nil {
// 		return tracerr.Wrap(err)
// 	}
// 	builder := db.NewService().DDLBuilder()
// 	var stmts []string
// 	for _, table := range schema.Tables {
// 		stmts = append(stmts, createIndexForTable(builder, table)...)
// 	}
// 	for _, stmt := range stmts {
// 		fmt.Println(stmt)
// 	}
// 	return nil
// }

func createTableForTable(builder dbcm.DDL, table *model.Table) string {
	var stmts []string
	stmts = append(stmts, builder.CreateTable(table))
	for _, index := range createIndexForTable(builder, table) {
		stmts = append(stmts, index)
	}
	stmts = append(stmts, "")
	for _, constraint := range addConstraintForTable(builder, table) {
		stmts = append(stmts, constraint)
	}
	return strings.Join(stmts, "\n")
}

func createIndexForTable(builder dbcm.DDL, table *model.Table) []string {
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

func addConstraintForTable(builder dbcm.DDL, table *model.Table) []string {
	var stmts []string
	for _, column := range table.Columns {
		if column.ForeignKey != "" {
			stmts = append(stmts, builder.AddConstraint(table.Name, column))
		}
	}
	return stmts
}

func dropIndexForTable(builder dbcm.DDL, table *model.Table) []string {
	var stmts []string
	for _, column := range table.Columns {
		if util.IsYes(column.Index) {
			stmts = append(stmts, builder.DropIndex(table.Name, column.Name))
		}
	}
	return stmts
}

func dropConstraintForTable(builder dbcm.DDL, table *model.Table) []string {
	var stmts []string
	for _, column := range table.Columns {
		if column.ForeignKey != "" {
			stmts = append(stmts, builder.DropConstraint(table.Name, column.Name))
		}
	}
	return stmts
}
