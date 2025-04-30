package main

import (
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/tpl"
)

func RegisterTxt(cliapp *cli.App) {
	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		var input, output, template, schema, table string
		var dump bool

		return &cli.Command{
			Name:    "text",
			Usage:   "transform from yaml to text",
			Aliases: []string{"t"},
			Flags: []cli.Flag{
				schemaFileFlagBuilder().WithDestination(&input).Build(),
				outputFileFlagBuilder().WithUsage("output file (text file)").WithDestination(&output).Build(),
				schemaNameFlagBuilder().WithDestination(&schema).Build(),
				tableNameFlagBuilder().WithDestination(&table).Build(),
				templateFileFlagBuilder().WithDestination(&template).Build(),
				flagbuilder.NewBoolFlag("dump").WithUsage("dump the content in .yml format").WithDestination(&dump).Build(),
			},
			Action: func(c *cli.Context) error {
				config.Setting.Input = lo.If(input != "", input).Else(config.Setting.Input)
				config.Setting.Text.Output = lo.If(output != "", output).Else(config.Setting.Text.Output)
				config.Setting.Text.Schema = lo.If(schema != "", schema).Else(config.Setting.Text.Schema)
				config.Setting.Text.Table = lo.If(table != "", table).Else(config.Setting.Text.Table)
				config.Setting.Text.Template = lo.If(template != "", template).Else(config.Setting.Text.Template)

				// data, err := fileloader.Load(input) // the data does not filter by schema and table
				// if err != nil {
				// 	return tracerr.Wrap(err)
				// }

				// // dump the original content in .yml format to console
				// if dump {
				// 	if err := fileloader.DumpYml(data, output, schema, table); err != nil {
				// 		return tracerr.Wrap(err)
				// 	}
				// 	return nil
				// }
				// if template != "" {
				// 	filteredData, err := utils.Filter(data, schema, table, "")
				// 	// filter the data with pattern
				// 	if err != nil {
				// 		return tracerr.Wrap(err)
				// 	}
				// 	return tpl.WriteWithFileLoader(filteredData, template, output)
				// }
				// if err := filewriter.WriteYml(data, output, schema, table); err != nil {
				// 	return tracerr.Wrap(err)
				// }
				return tpl.Generate()
			},
		}
	}())
}
