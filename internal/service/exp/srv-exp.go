package exp

import (
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/db"
	"github.com/zrs01/dst/internal/ddwriter"
	"github.com/zrs01/dst/internal/service"
	"github.com/ztrue/tracerr"
)

func Generate() error {
	dbService := db.NewService()
	schemaDef, err := dbService.Load(config.Setting.TableName)
	if err != nil {
		return tracerr.Wrap(err)
	}
	if err := service.RemoveCommonColumns(schemaDef); err != nil {
		return tracerr.Wrap(err)
	}
	return ddwriter.OutputYml(schemaDef, config.Setting.Output)
}
