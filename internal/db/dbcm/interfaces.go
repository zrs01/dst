package dbcm

import "github.com/zrs01/dst/model"

type Service interface {
	DDLBuilder() DDL
	Load(tableFilter string) (*model.Schema, error)
}

type DDL interface {
	CreateDatabase(schemaModel *model.Schema) (string, error)
	CreateTable()
	CreateView()
	CreateUser()
	CreateTrigger()
	CreateFunction()
	CreateIndex(indexName, tableName string, fields []string, isUnique bool) string
	CreateProcedure()

	AlterDatabase()
	AlterTable()
	AlterView()
	AlterUser()
	AlterTrigger()
	AlterFunction()
	AlterIndex()
	AlterProcedure()

	DropDatabase()
	DropTable()
	DropView()
	DropUser()
	DropTrigger()
	DropFunction()
	DropIndex()
	DropProcedure()
}
