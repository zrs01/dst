package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"github.com/zrs01/dst/internal/cliflag"
	"github.com/zrs01/dst/internal/erd"
	"github.com/zrs01/dst/internal/fileloader"
	"github.com/zrs01/dst/internal/filewriter"
	"github.com/zrs01/dst/internal/tpl"
	"github.com/ztrue/tracerr"

	"github.com/zrs01/dst/utils"
)

func registerTemplateRenderer(cliapp *cli.App) *cli.Command {
	// convert command
	convertCmd := &cli.Command{
		Name:    "convert",
		Aliases: []string{"c"},
		Usage:   "Convert to other format",
	}
	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		return convertCmd
	}())

	// transform to text
	convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, tfile, schema, table string
		var dump bool
		return &cli.Command{
			Name:    "text",
			Usage:   "transform from yaml to text",
			Aliases: []string{"t"},
			Flags: []cli.Flag{
				schemaFileFlagBuilder().WithDestination(&ifile).Build(),
				outputFileFlagBuilder().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
				schemaNameFlagBuilder().WithDestination(&schema).Build(),
				tableNameFlagBuilder().WithDestination(&table).Build(),
				templateFileFlagBuilder().WithDestination(&tfile).Build(),
				cliflag.NewBoolFlagBuilder("dump").WithUsage("dump the content in .yml format").WithDestination(&dump).Build(),
			},
			Action: func(c *cli.Context) error {
				data, err := fileloader.Load(ifile) // the data does not filter by schema and table
				if err != nil {
					return tracerr.Wrap(err)
				}

				// dump the original content in .yml format to console
				if dump {
					if err := fileloader.DumpYml(data, ofile, schema, table); err != nil {
						return tracerr.Wrap(err)
					}
					return nil
				}
				if tfile != "" {
					filteredData, err := utils.Filter(data, schema, table, "")
					// filter the data with pattern
					if err != nil {
						return tracerr.Wrap(err)
					}
					return tpl.WriteFileTpl(filteredData, tfile, ofile)
				}
				if err := filewriter.WriteYml(data, ofile, schema, table); err != nil {
					return tracerr.Wrap(err)
				}
				return nil
			},
		}
	}())

	// transform to excel
	convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table string
		var simple bool
		return &cli.Command{
			Name:    "excel",
			Usage:   "transform from yaml to excel",
			Aliases: []string{"e"},
			Flags: []cli.Flag{
				schemaFileFlagBuilder().WithDestination(&ifile).Build(),
				outputFileFlagBuilder().WithUsage("output file (.xlsx)").WithDestination(&ofile).Build(),
				schemaNameFlagBuilder().WithDestination(&schema).Build(),
				tableNameFlagBuilder().WithDestination(&table).Build(),
				cliflag.NewBoolFlagBuilder("simple").WithUsage("simple content").WithDestination(&simple).Build(),
				// iSchemaFileFlag(&ifile),
				// ofileFlag(&ofile, "output file (.xlsx)"),
				// schemaFile(&schema),
				// tableFlag(&table),
				// simpleFlag(&simple),
			},
			Action: func(c *cli.Context) error {
				oext := lo.Ternary(ofile != "", strings.ToLower(filepath.Ext(ofile)), "")
				data, err := fileloader.LoadWithFilter(ifile, schema, table, "")
				if err != nil {
					return tracerr.Wrap(err)
				}
				switch oext {
				case ".xlsx":
					if err := filewriter.WriteXlsx(data, ofile, simple); err != nil {
						return tracerr.Wrap(err)
					}
					return nil
				}
				return tracerr.New("Not implemented yet")
			},
		}
	}())

	// transform to diagram
	convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, tfile, schema, table, lib string
		return &cli.Command{
			Name:    "diagram",
			Usage:   "transform from yaml to diagram",
			Aliases: []string{"d"},
			Flags: []cli.Flag{
				schemaFileFlagBuilder().WithDestination(&ifile).Build(),
				outputFileFlagBuilder().WithUsage("output file (.png)").WithDestination(&ofile).Build(),
				schemaNameFlagBuilder().WithDestination(&schema).Build(),
				tableNameFlagBuilder().WithDestination(&table).Build(),
				templateFileFlagBuilder().WithDestination(&tfile).Build(),
				cliflag.NewStringFlagBuilder("lib").WithUsage("plantuml.jar file, used when output format is png").WithDestination(&lib).Build(),
				// iSchemaFileFlag(&ifile),
				// ofileFlag(&ofile, "output file (.png)"),
				// schemaFile(&schema),
				// tableFlag(&table),
				// templateFlag(&tfile),
				// libFlag(&lib),
			},
			Action: func(c *cli.Context) error {
				if ofile == "" {
					ofile = strings.TrimSuffix(ifile, filepath.Ext(ifile)) + ".png"
				}
				oext := lo.Ternary(ofile != "", strings.ToLower(filepath.Ext(ofile)), "")
				data, err := fileloader.LoadWithFilter(ifile, schema, table, "")
				if err != nil {
					return tracerr.Wrap(err)
				}
				switch oext {
				case ".png":
					if err := erd.WriteERD(data, tfile, ofile); err != nil {
						return tracerr.Wrap(err)
					}
					return nil
				}
				return tracerr.New(fmt.Sprintf("output file extension '%s' is not supported", ofile))
			},
		}
	}())
	return convertCmd
}
