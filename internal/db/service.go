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
	// DDLBuilder()
	Load(tableFilter string) (*model.DataDef, error)
}
