package dst

import (
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
)

func RegisterExportCmd(cliapp *cli.App) {
	var dsn, output, table string

	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		return &cli.Command{
			Name:  "export",
			Usage: "Export the database schema to schema file",
			Flags: []cli.Flag{
				flagbuilder.NewStringFlag("dsn").WithUsage("database source name").WithDestination(&dsn).Build(),
				outputFileFlagBuilder().WithUsage("output file (file extension must be either .puml or .png)").WithDestination(&output).Build(),
				tableNameFlagBuilder().WithDestination(&table).Build(),
			},
			Action: func(c *cli.Context) error {
				config.Setting.Export.Dsn = lo.If(dsn != "", dsn).Else(config.Setting.Export.Dsn)
				config.Setting.Export.Output = lo.If(output != "", output).Else(config.Setting.Export.Output)
				config.Setting.Export.TableFilter = lo.If(table != "", table).Else(config.Setting.Export.TableFilter)
				// return erd.Generate()
				return nil
			},
		}
	}())
}
