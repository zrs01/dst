package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"

	"dst/transform"
)

var version = "development"

// input   output   options
// ------------------------
// yaml    xlsx     simple
// yaml    <text>   template
// yaml    png      plantuml.jar
// xlsx    yaml
// TODO:
// dbase   yaml     orginal yaml

func main() {

	cliapp := cli.NewApp()
	cliapp.Name = "dst"
	cliapp.Usage = "Database schema tool"
	cliapp.Version = version
	cliapp.Commands = []*cli.Command{}

	debug := false

	// global options
	cliapp.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Usage:       "Debug mode",
			Required:    false,
			Destination: &debug,
		},
	}

	/* ------------------------------ Common flags ------------------------------ */

	ifileFlag := func(file *string, usage string) *cli.StringFlag {
		return &cli.StringFlag{Name: "input", Aliases: []string{"i"}, Usage: lo.Ternary(usage == "", "input file", usage), Required: true, Destination: file}
	}
	ofileFlag := func(file *string, usage string) *cli.StringFlag {
		return &cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: lo.Ternary(usage == "", "output file", usage), Required: false, Destination: file}
	}
	templateFlag := func(file *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "template", Aliases: []string{"t"}, Usage: "template file", Required: false, Destination: file}
	}
	schemaFile := func(schema *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "schema", Usage: "schema name pattern, wildcard char: * or %", Required: false, Destination: schema}
	}
	tableFlag := func(table *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "table", Usage: "table name pattern, wildcard char: * or %", Required: false, Destination: table}
	}
	simpleFlag := func(simple *bool) *cli.BoolFlag {
		return &cli.BoolFlag{Name: "simple", Usage: "simple content", Value: false, Required: false, Destination: simple}
	}
	libFlag := func(lib *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "lib", Usage: "plantuml.jar file, used when output format is png", Required: false, Destination: lib}
	}
	selectedData := func(ifile, schema, table string) (*transform.DataDef, error) {
		rawData, err := transform.ReadYml(ifile)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		validateResult := transform.Verify(rawData)
		if len(validateResult) > 0 {
			lo.ForEach(validateResult, func(v string, _ int) {
				fmt.Println(v)
			})
			return nil, tracerr.Errorf("invalid data")
		}
		data := transform.FilterData(rawData, schema, table)
		return data, nil
	}

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
		var simple bool
		return &cli.Command{
			Name:    "text",
			Usage:   "transform from yaml to text",
			Aliases: []string{"t"},
			Flags: []cli.Flag{
				ifileFlag(&ifile, "input file (.yml)"),
				ofileFlag(&ofile, "output file (text file)"),
				schemaFile(&schema),
				tableFlag(&table),
				simpleFlag(&simple),
				templateFlag(&tfile),
			},
			Action: func(c *cli.Context) error {
				oext := lo.Ternary(ofile != "", strings.ToLower(filepath.Ext(ofile)), "")
				data, err := selectedData(ifile, schema, table)
				if err != nil {
					return tracerr.Wrap(err)
				}
				if tfile != "" {
					return transform.WriteTpl(data, tfile, ofile)
				}
				switch oext {
				case ".yml", "":
					if err := transform.WriteYml(data, ofile); err != nil {
						return tracerr.Wrap(err)
					}
				}
				return tracerr.New("Not implemented yet")
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
				ifileFlag(&ifile, "input file (.yml)"),
				ofileFlag(&ofile, "output file (.xlsx)"),
				schemaFile(&schema),
				tableFlag(&table),
				simpleFlag(&simple),
			},
			Action: func(c *cli.Context) error {
				oext := lo.Ternary(ofile != "", strings.ToLower(filepath.Ext(ofile)), "")
				data, err := selectedData(ifile, schema, table)
				if err != nil {
					return tracerr.Wrap(err)
				}
				switch oext {
				case ".xlsx":
					if err := transform.WriteXlsx(data, ofile, simple); err != nil {
						return tracerr.Wrap(err)
					}
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
				ifileFlag(&ifile, "input file (.yml)"),
				ofileFlag(&ofile, "output file (.png)"),
				schemaFile(&schema),
				tableFlag(&table),
				templateFlag(&tfile),
				libFlag(&lib),
			},
			Action: func(c *cli.Context) error {
				if ofile == "" {
					ofile = strings.TrimSuffix(ifile, filepath.Ext(ifile)) + ".png"
				}
				oext := lo.Ternary(ofile != "", strings.ToLower(filepath.Ext(ofile)), "")
				data, err := selectedData(ifile, schema, table)
				if err != nil {
					return tracerr.Wrap(err)
				}
				switch oext {
				case ".png":
					if err := transform.WriteERD(data, tfile, ofile); err != nil {
						return tracerr.Wrap(err)
					}
					return nil
				}
				return tracerr.New(fmt.Sprintf("output file extension '%s' is not supported", ofile))
			},
		}
	}())

	if err := cliapp.Run(os.Args); err != nil {
		tracerr.Print(err)
	}
}
