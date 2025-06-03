package dst

import (
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/text"
)

func RegisterTemplateCmd(cliapp *cli.App) {
	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		var input, output, template, schema, table string
		// var dump bool

		return &cli.Command{
			Name:  "template",
			Usage: "Transform YAML data using a template into a text file",
			Flags: []cli.Flag{
				schemaFileFlagBuilder().WithDestination(&input).Build(),
				outputFileFlagBuilder().WithUsage("output file (text file)").WithDestination(&output).Build(),
				schemaNameFlagBuilder().WithDestination(&schema).Build(),
				tableNameFlagBuilder().WithDestination(&table).Build(),
				templateFileFlagBuilder().WithDestination(&template).Build(),
				// flagbuilder.NewBoolFlag("dump").WithUsage("dump the content in .yml format").WithDestination(&dump).Build(),
			},
			Action: func(c *cli.Context) error {
				config.Setting.Text.Input = lo.If(input != "", input).Else(config.Setting.Text.Input)
				config.Setting.Text.Output = lo.If(output != "", output).Else(config.Setting.Text.Output)
				config.Setting.Text.SchemaFilter = lo.If(schema != "", schema).Else(config.Setting.Text.SchemaFilter)
				config.Setting.Text.TableFilter = lo.If(table != "", table).Else(config.Setting.Text.TableFilter)
				config.Setting.Text.Template = lo.If(template != "", template).Else(config.Setting.Text.Template)
				return text.Generate()
			},
		}
	}())
}
