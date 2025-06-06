package dst

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service/text"
)

func RegisterTextCmd() *cli.Command {
	var options config.SettingDef

	return &cli.Command{
		Name:  "text",
		Usage: "transform using a template into a text file",
		Flags: []cli.Flag{
			flagbuilder.InputFileFlag().WithDestination(&options.Input).Build(),
			flagbuilder.OutputFileFlag().WithDestination(&options.Output).Build(),
			flagbuilder.TableNameFlag().WithDestination(&options.TableName).Build(),
			flagbuilder.TemplateFileFlag().WithDestination(&options.Template).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config.InitSetting(config.TextConf, options)
			return text.Generate()
		},
	}
}
