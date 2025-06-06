package ddl

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/service"
	"github.com/ztrue/tracerr"
)

func GenerateCreateDatabase() error {
	data, err := service.NewLoadBuilder(
		// service.WithSchemaPattern(config.Setting.SchemaFilter),
		service.WithTablePattern(config.Setting.TableName)).
		LoadFromFile(config.Setting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// builder := factory.NewService(config.Setting.Dsn, config.Setting.Ccf).DDLBuilder()
	// builder.WithDataDef(data).CreateDatabase()
	spew.Dump(data)
	return nil
}
