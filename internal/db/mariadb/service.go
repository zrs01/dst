package mariadb

import (
	"regexp"
	"strings"

	"github.com/zrs01/dst/internal/db"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

type MariadbService struct {
	dsn string
	ccf string // common column file
}

func NewMariadbService(dsn, ccf string) db.Service {
	return &MariadbService{
		dsn: dsn,
		ccf: ccf,
	}
}

func (s *MariadbService) Load(tableFilter string) (*model.Schema, error) {
	// dataSourceName: mariadb://[username[:password]@][protocol[(address[:port])]]/dbname[?param1=value1&...&paramN=valueN]
	regex := regexp.MustCompile(`\w+\://[^/]*/(\w+)`).FindStringSubmatch(s.dsn)
	if len(regex) == 0 {
		return nil, tracerr.Errorf("Failed to parse %s", s.dsn)
	}

	dsName := s.dsn
	parts := strings.Split(s.dsn, "://")
	if len(parts) > 1 {
		dsName = parts[1]
	}

	databaseName := regex[1]
	dbManager := NewExportBuilder().
		WithDriverName("mysql").
		WithDataSourceName(dsName).
		WithDatabaseName(databaseName).
		WithTableName(tableFilter).
		WithCommonColumnFile(s.ccf)
	schemaModel, err := dbManager.ToSchemaModel()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return schemaModel, nil
}

func (s *MariadbService) DDLBuilder() db.DDL {
	return nil
}
