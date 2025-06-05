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
		Usage: "Export the database schema",
		Flags: []cli.Flag{
			flagbuilder.NewStringFlag("dsn").WithUsage("database source name").WithDestination(&options.Dsn).Build(),
			flagbuilder.NewStringFlag("common").WithAliases("cc").WithUsage("common column file").WithDestination(&options.CommonColumnFile).Build(),
			flagbuilder.OutputFileFlag().WithUsage("output file (file extension must be either .puml or .png)").WithDestination(&options.Output).Build(),
			flagbuilder.TableNameFlag().WithDestination(&options.TableFilter).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config.InitSetting(config.ExportConf, options)
			return exp.Generate()
		},
	}
}
