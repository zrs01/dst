package db

import (
	"strings"

	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/db/dbcm"
	"github.com/zrs01/dst/internal/db/mariadb"
)

func NewService() dbcm.Service {
	switch {
	// case strings.HasPrefix(dsn, MYSQL_PREFIX):
	// 	return mariadb.NewMariadbService(dsn, ccf)
	// case strings.HasPrefix(dsn, MSSQL_PREFIX):
	// 	return mssql.NewMssqlService(dsn)
	case strings.HasPrefix(config.Setting.Dsn, dbcm.MARIADB_PREFIX):
		return mariadb.NewMariadbService()
	default:
		return nil
	}
}
