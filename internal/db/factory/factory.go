package factory

import (
	"strings"

	"github.com/zrs01/dst/internal/db"
	"github.com/zrs01/dst/internal/db/mariadb"
)

func Service(dsn, ccf string) db.Service {
	switch {
	// case strings.HasPrefix(dsn, MYSQL_PREFIX):
	// 	return mariadb.NewMariadbService(dsn, ccf)
	// case strings.HasPrefix(dsn, MSSQL_PREFIX):
	// 	return mssql.NewMssqlService(dsn)
	case strings.HasPrefix(dsn, db.MARIADB_PREFIX):
		return mariadb.NewMariadbService(dsn, ccf)
	default:
		return nil
	}
}
