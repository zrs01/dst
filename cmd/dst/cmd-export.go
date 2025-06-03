package dst

import (
	"context"

	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/export"
	"github.com/zrs01/dst/internal/flagbuilder"
)

func RegisterExportCmd() *cli.Command {
	var dsn, ccf, output, table string

	return &cli.Command{
		Name:  "export",
		Usage: "Export the database schema to schema file",
		Flags: []cli.Flag{
			flagbuilder.NewStringFlag("dsn").WithUsage("database source name").WithDestination(&dsn).Build(),
			flagbuilder.NewStringFlag("common").WithAliases("cc").WithUsage("common column file").WithDestination(&ccf).Build(),
			flagbuilder.OutputFileFlag().WithUsage("output file (file extension must be either .puml or .png)").WithDestination(&output).Build(),
			flagbuilder.TableNameFlag().WithDestination(&table).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config.Setting.Export.Dsn = lo.If(dsn != "", dsn).Else(config.Setting.Export.Dsn)
			config.Setting.Export.CommonColumnFile = lo.If(ccf != "", ccf).Else(config.Setting.Export.CommonColumnFile)
			config.Setting.Export.Output = lo.If(output != "", output).Else(config.Setting.Export.Output)
			config.Setting.Export.TableFilter = lo.If(table != "", table).Else(config.Setting.Export.TableFilter)
			return export.Generate()
		},
	}
}
