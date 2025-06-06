package mariadb

import (
	"regexp"
	"strings"

	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/db/dbcm"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

type MariadbService struct{}

func NewMariadbService() dbcm.Service {
	return &MariadbService{}
}

func (s *MariadbService) Load(tableFilter string) (*model.Schema, error) {
	// dataSourceName: mariadb://[username[:password]@][protocol[(address[:port])]]/dbname[?param1=value1&...&paramN=valueN]
	regex := regexp.MustCompile(`\w+\://[^/]*/(\w+)`).FindStringSubmatch(config.Setting.Dsn)
	if len(regex) == 0 {
		return nil, tracerr.Errorf("Failed to parse %s", config.Setting.Dsn)
	}

	dsName := config.Setting.Dsn
	parts := strings.Split(config.Setting.Dsn, "://")
	if len(parts) > 1 {
		dsName = parts[1]
	}

	databaseName := regex[1]
	dbManager := NewExportBuilder().
		WithDriverName("mysql").
		WithDataSourceName(dsName).
		WithDatabaseName(databaseName).
		WithTableName(tableFilter).
		WithCommonColumnFile(config.Setting.Ccf)
	schemaModel, err := dbManager.ToSchemaModel()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return schemaModel, nil
}

func (s *MariadbService) DDLBuilder() dbcm.DDL {
	return NewDDBuilder()
}
