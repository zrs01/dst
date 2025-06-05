package dst

import (
	"context"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/ddwriter"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service"
	"github.com/ztrue/tracerr"
)

func RegisterExcelCmd() *cli.Command {
	var options config.SettingDef
	var simple bool

	return &cli.Command{
		Name:  "excel",
		Usage: "transform from yaml to excel",
		Flags: []cli.Flag{
			flagbuilder.SchemaFileFlag().WithDestination(&options.Input).Build(),
			flagbuilder.OutputFileFlag().WithUsage("output file (.xlsx)").WithDestination(&options.Output).Build(),
			flagbuilder.TableNameFlag().WithDestination(&options.TableFilter).Build(),
			flagbuilder.NewBoolFlag("simple").WithUsage("simple content").WithDestination(&simple).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config.InitSetting(config.ExcelConf, options)
			schemaDef, err := service.NewLoadBuilder(
				service.WithTablePattern(config.Setting.TableFilter)).
				LoadFromFile(config.Setting.Input)
			if err != nil {
				return tracerr.Wrap(err)
			}
			if !strings.HasSuffix(config.Setting.Output, ".xlsx") {
				config.Setting.Output = config.Setting.Output + ".xlsx"
			}

			if err := ddwriter.WriteXlsx(schemaDef, config.Setting.Output, simple); err != nil {
				return tracerr.Wrap(err)
			}
			return nil
		},
	}
}
