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
	fmt.Printf("Retrieving database definition from '%s' ... ", config.Setting.Dsn)
	dbSchema, err := db.NewService().Load(config.Setting.TableName)
	if err != nil {
		return tracerr.Wrap(err)
	}
	fmt.Printf("done\n")

	// get model from schema definition file
	loadBuilder := service.NewLoadBuilder(
		service.WithTablePattern(config.Setting.TableName),
	)
	fmt.Printf("Retrieving schema definition from '%s' ... ", config.Setting.Sdf)
	schema, err := loadBuilder.LoadFromFile(config.Setting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	fmt.Printf("done\n\n")
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
			for i, column := range table.Columns {
				dbColumn, found := lo.Find(dbTable.Columns, func(c *model.Column) bool {
					return c.Name == column.Name
				})
				if !found {
					// new column
					fmt.Println(ddlBuilder.AddColumn(table.Name, column, table.Columns[i-1].Name))
					indexStmt := createIndexForColumn(ddlBuilder, table.Name, column)
					if indexStmt != "" {
						fmt.Printf("\n%s\n", indexStmt)
					}
				} else {
					// update column
					if dbColumn.DataType != column.DataType ||
						lo.If(dbColumn.NotNull == "N", "").Else(dbColumn.NotNull) != column.NotNull ||
						dbColumn.Desc != column.Desc ||
						lo.If(dbColumn.Unique == "N", "").Else(dbColumn.Unique) != column.Unique {
						fmt.Println(ddlBuilder.AlterTableModifyColumn(table.Name, column))
					}
				}
			}
		}
	}

	for _, table := range dbSchema.Tables {
		_, found := lo.Find(schema.Tables, func(t *model.Table) bool {
			return t.Name == table.Name
		})
		if !found {
			// drop table
			for _, stmt := range ddlBuilder.DropTable(table) {
				fmt.Println(stmt)
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
	if config.Setting.IsAlter {
		stmts = append(stmts, builder.CreateTableWithAlter(table))
	} else {
		stmts = append(stmts, builder.CreateTableWithDirect(table))
	}
	for _, indexStmt := range createIndexForTable(builder, table) {
		stmts = append(stmts, indexStmt)
	}
	stmts = append(stmts, "")
	for _, constraint := range addConstraintForTable(builder, table) {
		stmts = append(stmts, constraint)
	}
	return strings.Join(stmts, "\n")
}

func createIndexForTable(builder dbcm.DDL, table *model.Table) []string {
	var stmts []string
	for _, column := range table.Columns {
		indexStmt := createIndexForColumn(builder, table.Name, column)
		if indexStmt != "" {
			stmts = append(stmts, indexStmt)
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

func createIndexForColumn(builder dbcm.DDL, tableName string, column *model.Column) string {
	var index string
	if util.IsYes(column.Index) {
		indexName := fmt.Sprintf("idx_%s_%s", tableName, column.Name)
		index = builder.CreateIndex(indexName, tableName, []string{column.Name}, util.IsYes(column.Unique))
	}
	return index
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
