package ddl

import (
	"fmt"

	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/db"
	"github.com/zrs01/dst/internal/service"
	"github.com/ztrue/tracerr"
)

func GenerateCreateDatabase() error {
	schemaModel, err := service.NewLoadBuilder().LoadFromFile(config.Setting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	builder := db.NewService().DDLBuilder()
	output, err := builder.CreateDatabase(schemaModel)
	if err != nil {
		return tracerr.Wrap(err)
	}
	fmt.Println(output)
	return nil
}

func GenerateCreateIndex() error {
	loadBuilder := service.NewLoadBuilder(
		service.WithTablePattern(config.Setting.TableName),
	)
	schemaModel, err := loadBuilder.LoadFromFile(config.Setting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	schemaModel, err = loadBuilder.Filter(schemaModel)
	if err != nil {
		return tracerr.Wrap(err)
	}

	builder := db.NewService().DDLBuilder()
	output, err := builder.CreateIndex(schemaModel)
	if err != nil {
		return tracerr.Wrap(err)
	}
	for _, o := range output {
		fmt.Println(o)
	}
	return nil
}
