package dst

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service/txt"
	"github.com/ztrue/tracerr"
)

func RegisterTextCmd() *cli.Command {
	var options config.SettingDef

	return &cli.Command{
		Name:  "txt",
		Usage: "transform using a template into a text file",
		Flags: []cli.Flag{
			flagbuilder.SdfFlag().WithDestination(&options.Sdf).Build(),
			flagbuilder.OutputFileFlag().WithDestination(&options.Output).Build(),
			flagbuilder.TableNameFlag().WithDestination(&options.TableName).Build(),
			flagbuilder.TemplateFileFlag().WithDestination(&options.Template).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := config.LoadConfig(config.TextConf, options); err != nil {
				return tracerr.Wrap(err)
			}
			return txt.Generate()
		},
	}
}
