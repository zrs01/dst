package dst

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service/xls"
)

func RegisterXlsCmd() *cli.Command {
	var options config.SettingDef
	var simple bool

	return &cli.Command{
		Name:  "xls",
		Usage: "transform from yaml to excel",
		Flags: []cli.Flag{
			flagbuilder.InputFileFlag().WithDestination(&options.Input).Build(),
			flagbuilder.OutputFileFlag().WithUsage("output file (.xlsx)").WithDestination(&options.Output).Build(),
			flagbuilder.TableNameFlag().WithDestination(&options.TableName).Build(),
			flagbuilder.NewBoolFlag("simple").WithUsage("simple content").WithDestination(&simple).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config.InitSetting(config.ExcelConf, options)
			return xls.Generate()
		},
	}
}
