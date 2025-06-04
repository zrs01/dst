package dst

import (
	"context"

	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service/erd"
)

func RegisterERDCmd() *cli.Command {
	var input, output, template, schema, table, lib string

	return &cli.Command{
		Name:  "erd",
		Usage: "Generate ERD diagram",
		Flags: []cli.Flag{
			flagbuilder.SchemaFileFlag().WithDestination(&input).Build(),
			flagbuilder.OutputFileFlag().WithUsage("output file (file extension must be either .puml or .png)").WithDestination(&output).Build(),
			flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
			flagbuilder.TableNameFlag().WithDestination(&table).Build(),
			flagbuilder.TemplateFileFlag().WithDestination(&template).Build(),
			flagbuilder.NewStringFlag("umllib").WithUsage("Path of plantuml.jar file (required when output is .png").WithDestination(&lib).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config.Setting.Erd.Input = lo.If(input != "", input).Else(config.Setting.Erd.Input)
			config.Setting.Erd.Output = lo.If(output != "", output).Else(config.Setting.Erd.Output)
			config.Setting.Erd.SchemaFilter = lo.If(schema != "", schema).Else(config.Setting.Erd.SchemaFilter)
			config.Setting.Erd.TableFilter = lo.If(table != "", table).Else(config.Setting.Erd.TableFilter)
			config.Setting.Erd.Template = lo.If(template != "", template).Else(config.Setting.Erd.Template)
			config.Setting.Erd.PlantumlLib = lo.If(lib != "", lib).Else(config.Setting.Erd.PlantumlLib)
			return erd.Generate()
		},
	}
}
