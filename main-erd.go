package main

import (
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/erd"
	"github.com/zrs01/dst/internal/flagbuilder"
)

func RegisterErd(cliapp *cli.App) {
	var input, output, template, schema, table, lib string

	userText := `Example:
  dst erd --input schema.sql --output schema.png
	`
	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		return &cli.Command{
			Name:      "erd",
			Aliases:   []string{"d"},
			Usage:     "Generate ERD diagram",
			UsageText: userText,
			Flags: []cli.Flag{
				schemaFileFlagBuilder().WithDestination(&input).Build(),
				outputFileFlagBuilder().WithUsage("output file (file extension must be either .puml or .png)").WithDestination(&output).Build(),
				schemaNameFlagBuilder().WithDestination(&schema).Build(),
				tableNameFlagBuilder().WithDestination(&table).Build(),
				templateFileFlagBuilder().WithDestination(&template).Build(),
				flagbuilder.NewStringFlag("umllib").WithUsage("Path of plantuml.jar file (required when output is .png").WithDestination(&lib).Build(),
			},
			Action: func(c *cli.Context) error {
				config.Setting.Input = lo.If(input != "", input).Else(config.Setting.Input)
				config.Setting.Erd.Output = lo.If(output != "", output).Else(config.Setting.Erd.Output)
				config.Setting.Erd.SchemaFilter = lo.If(schema != "", schema).Else(config.Setting.Erd.SchemaFilter)
				config.Setting.Erd.TableFilter = lo.If(table != "", table).Else(config.Setting.Erd.TableFilter)
				config.Setting.Erd.Template = lo.If(template != "", template).Else(config.Setting.Erd.Template)
				config.Setting.Erd.UmlLib = lo.If(lib != "", lib).Else(config.Setting.Erd.UmlLib)
				return erd.Generate()
			},
		}
	}())
}
