package util

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	"github.com/ztrue/tracerr"
)

func OpenSqlDb(driverName, dataSourceName string) (*sql.DB, error) {
	if driverName == "" || dataSourceName == "" {
		return nil, tracerr.Errorf("The database driver name or source name was not provided in configuration")
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return db, nil
}

func Query(db *sql.DB, query string, args []any, rowFunc func(*sql.Rows) error) error {
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.Debugf("\n----------\nargs: %+v\n%s\n", args, query)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rowFunc(rows); err != nil {
			return tracerr.Wrap(err)
		}
	}
	if err := rows.Err(); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}
