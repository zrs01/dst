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
	// var input, output, template, schema, table string

	return &cli.Command{
		Name:  "text",
		Usage: "Transform using a template into a text file",
		Flags: []cli.Flag{
			flagbuilder.SchemaFileFlag().WithDestination(&options.Input).Build(),
			flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&options.Output).Build(),
			flagbuilder.SchemaNameFlag().WithDestination(&options.SchemaFilter).Build(),
			flagbuilder.TableNameFlag().WithDestination(&options.TableFilter).Build(),
			flagbuilder.TemplateFileFlag().WithDestination(&options.Template).Build(),
			// flagbuilder.NewBoolFlag("dump").WithUsage("dump the content in .yml format").WithDestination(&dump).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config.InitSetting(config.TextConf, options)
			return text.Generate()
		},
	}
}
