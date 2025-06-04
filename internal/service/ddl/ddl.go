package ddl

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/ddloader"
	"github.com/ztrue/tracerr"
)

func Generate() error {
	logrus.Info("Generating DDL...")
	data, err := ddloader.NewLoadBuilder(
		ddloader.WithSchemaPattern(config.Setting.Erd.SchemaFilter),
		ddloader.WithTablePattern(config.Setting.Erd.TableFilter)).
		LoadFromFile(config.Setting.Erd.Input)
	if err != nil {
		return tracerr.Wrap(err)
	}
	spew.Dump(data)
	return nil
}
