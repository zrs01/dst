package db

import (
	"github.com/zrs01/dst/model"
)

const (
	MYSQL_PREFIX   = "mysql://"
	MSSQL_PREFIX   = "sqlserver://"
	MARIADB_PREFIX = "mariadb://"
)

type Service interface {
	DDLBuilder() DDL
	Load(tableFilter string) (*model.Schema, error)
}

type DDL interface {
	CreateDatabase() (string, error)
	CreateTable()
	CreateView()
	CreateUser()
	CreateTrigger()
	CreateFunction()
	CreateIndex() ([]string, error)
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
