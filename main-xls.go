package main

import (
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"github.com/zrs01/dst/internal/ddloader"
	"github.com/zrs01/dst/internal/ddwriter"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/ztrue/tracerr"
)

func RegisterXls(cliapp *cli.App) {
	var ifile, ofile, schema, table string
	var simple bool

	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		return &cli.Command{
			Name:    "excel",
			Usage:   "transform from yaml to excel",
			Aliases: []string{"e"},
			Flags: []cli.Flag{
				schemaFileFlagBuilder().WithDestination(&ifile).Build(),
				outputFileFlagBuilder().WithUsage("output file (.xlsx)").WithDestination(&ofile).Build(),
				schemaNameFlagBuilder().WithDestination(&schema).Build(),
				tableNameFlagBuilder().WithDestination(&table).Build(),
				flagbuilder.NewBoolFlag("simple").WithUsage("simple content").WithDestination(&simple).Build(),
			},
			Action: func(c *cli.Context) error {
				oext := lo.Ternary(ofile != "", strings.ToLower(filepath.Ext(ofile)), "")
				// data, err := ddloader.LoadWithFilter(ifile, schema, table, "")
				data, err := ddloader.Load(ifile,
					ddloader.WithSchemaPattern(schema),
					ddloader.WithTablePattern(table))
				if err != nil {
					return tracerr.Wrap(err)
				}
				switch oext {
				case ".xlsx":
					if err := ddwriter.WriteXlsx(data, ofile, simple); err != nil {
						return tracerr.Wrap(err)
					}
					return nil
				}
				return tracerr.New("Not implemented yet")
			},
		}
	}())
}
