package dst

import (
	"context"

	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service/text"
)

func RegisterTextCmd() *cli.Command {
	var input, output, template, schema, table string

	return &cli.Command{
		Name:  "text",
		Usage: "Transform using a template into a text file",
		Flags: []cli.Flag{
			flagbuilder.SchemaFileFlag().WithDestination(&input).Build(),
			flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&output).Build(),
			flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
			flagbuilder.TableNameFlag().WithDestination(&table).Build(),
			flagbuilder.TemplateFileFlag().WithDestination(&template).Build(),
			// flagbuilder.NewBoolFlag("dump").WithUsage("dump the content in .yml format").WithDestination(&dump).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config.Setting.Text.Input = lo.If(input != "", input).Else(config.Setting.Text.Input)
			config.Setting.Text.Output = lo.If(output != "", output).Else(config.Setting.Text.Output)
			config.Setting.Text.SchemaFilter = lo.If(schema != "", schema).Else(config.Setting.Text.SchemaFilter)
			config.Setting.Text.TableFilter = lo.If(table != "", table).Else(config.Setting.Text.TableFilter)
			config.Setting.Text.Template = lo.If(template != "", template).Else(config.Setting.Text.Template)
			return text.Generate()
		},
	}
}
