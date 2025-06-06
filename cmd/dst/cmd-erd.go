package dst

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service/erd"
)

func RegisterERDCmd() *cli.Command {
	var options config.SettingDef

	return &cli.Command{
		Name:  "erd",
		Usage: "Generate ERD diagram",
		Flags: []cli.Flag{
			flagbuilder.InputFileFlag().WithDestination(&options.Sdf).Build(),
			flagbuilder.OutputFileFlag().WithRequired(true).
				WithUsage("output file (file extension must be either .puml or .png)").WithDestination(&options.Output).Build(),
			flagbuilder.TableNameFlag().WithDestination(&options.TableName).Build(),
			flagbuilder.TemplateFileFlag().WithDestination(&options.Template).Build(),
			flagbuilder.NewStringFlag("plantuml").
				WithUsage("Path of plantuml.jar file (required when output is .png").WithDestination(&options.Plantuml).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config.InitSetting(config.ERDConf, options)
			return erd.Generate()
		},
	}
}
