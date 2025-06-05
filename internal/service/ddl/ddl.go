package ddl

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/service"
	"github.com/ztrue/tracerr"
)

func Generate() error {
	logrus.Info("Generating DDL...")
	data, err := service.NewLoadBuilder(
		// service.WithSchemaPattern(config.Setting.SchemaFilter),
		service.WithTablePattern(config.Setting.TableFilter)).
		LoadFromFile(config.Setting.Input)
	if err != nil {
		return tracerr.Wrap(err)
	}
	spew.Dump(data)
	return nil
}
