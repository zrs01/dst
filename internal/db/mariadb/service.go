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

func (s *MariadbService) Load(tableFilter string) (*model.DataDef, error) {
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

	schema := regex[1]
	dbManager := NewExportService().
		WithDriverName("mysql").
		WithDataSourceName(dsName).
		WithDatabaseName(schema).
		WithTableName(tableFilter).
		WithCommonColumnFile(s.ccf)
	schemaDef, err := dbManager.ToModel()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &model.DataDef{
		Schemas: []*model.Schema{schemaDef},
	}, nil
}
