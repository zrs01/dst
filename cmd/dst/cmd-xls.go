package dst

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service/xls"
	"github.com/ztrue/tracerr"
)

func RegisterXlsCmd() *cli.Command {
	var options config.SettingDef

	return &cli.Command{
		Name:  "xls",
		Usage: "transform from yaml to excel",
		Flags: []cli.Flag{
			flagbuilder.SdfFlag().WithDestination(&options.Sdf).Build(),
			flagbuilder.CommonColumnFlag().WithDestination(&options.Ccf).Build(),
			flagbuilder.OutputFileFlag().WithUsage("output file (.xlsx)").WithDestination(&options.Output).Build(),
			flagbuilder.TableNameFlag().WithDestination(&options.TableName).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := config.LoadConfig(config.ExcelConf, options); err != nil {
				return tracerr.Wrap(err)
			}
			return xls.Generate()
		},
	}
}
