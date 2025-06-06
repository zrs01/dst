package dst

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service/exp"
)

func RegisterExportCmd() *cli.Command {
	var options config.SettingDef

	return &cli.Command{
		Name:  "export",
		Usage: "export the database schema",
		Flags: []cli.Flag{
			flagbuilder.DsnFlag().WithDestination(&options.Dsn).Build(),
			flagbuilder.CommonColumnFlag().WithDestination(&options.CommonColumnFile).Build(),
			flagbuilder.OutputFileFlag().WithDestination(&options.Output).Build(),
			flagbuilder.TableNameFlag().WithDestination(&options.TableName).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config.InitSetting(config.ExportConf, options)
			return exp.Generate()
		},
	}
}
