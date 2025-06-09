package dbcm

import "github.com/zrs01/dst/model"

type Service interface {
	DDLBuilder() DDL
	Load(tableFilter string) (*model.Schema, error)
}

type DDL interface {
	CreateDatabase(schemaModel *model.Schema) (string, error)
	CreateTableWithDirect(table *model.Table) string
	CreateTableWithAlter(table *model.Table) string
	AddColumn(tableName string, column *model.Column, previousColumnName string) string
	AddConstraint(tableName string, column *model.Column) string
	// CreateView()
	// CreateUser()
	// CreateTrigger()
	// CreateFunction()
	CreateIndex(indexName, tableName string, fields []string, isUnique bool) string
	// CreateProcedure()

	AlterTableModifyColumn(tableName string, column *model.Column) string
	// AlterDatabase()
	// AlterTable()
	// AlterView()
	// AlterUser()
	// AlterTrigger()
	// AlterFunction()
	// AlterIndex()
	// AlterProcedure()

	// DropDatabase()
	DropTable(table *model.Table) []string
	// DropView()
	// DropUser()
	// DropTrigger()
	// DropFunction()
	DropIndex(tableName, columnName string) string
	DropConstraint(tableName, columnName string) string
	// DropProcedure()
}
