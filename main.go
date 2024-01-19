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
	srcData := func(ifile, schema, table string) (*transform.DataDef, error) {
		rawData, err := transform.ReadYml(ifile)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		data := transform.FilterData(rawData, schema, table)
		validateResult := transform.Verify(data)
		if len(validateResult) > 0 {
			lo.ForEach(validateResult, func(v string, _ int) {
				fmt.Println(v)
			})
			return nil, tracerr.Errorf("invalid data")
		}
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
				data, err := srcData(ifile, schema, table)
				if err != nil {
					return tracerr.Wrap(err)
				}
				if tfile != "" {
					return transform.WriteTpl(data, tfile, ofile, table)
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
				data, err := srcData(ifile, schema, table)
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
				oext := lo.Ternary(ofile != "", strings.ToLower(filepath.Ext(ofile)), "")
				data, err := srcData(ifile, schema, table)
				if err != nil {
					return tracerr.Wrap(err)
				}
				switch oext {
				case ".png":
					if err := transform.WriteERD(data, tfile, ofile); err != nil {
						return tracerr.Wrap(err)
					}
				}
				return tracerr.New("Not implemented yet")
			},
		}
	}())

	if err := cliapp.Run(os.Args); err != nil {
		tracerr.Print(err)
	}

	// convert command
	// convertCmd := &cli.Command{
	// 	Name:    "convert",
	// 	Aliases: []string{"c"},
	// 	Usage:   "Convert to other format",
	// }
	// cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
	// 	return convertCmd
	// }())

	// // from yml to xlsx
	// convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
	// 	var ifile, ofile string
	// 	return &cli.Command{
	// 		Name:    "yaml",
	// 		Aliases: []string{"y"},
	// 		Flags: []cli.Flag{
	// 			ifileFlag(&ifile),
	// 			ofileFlag(&ofile, ""),
	// 		},
	// 		Action: func(c *cli.Context) error {
	// 			if validInOutFile(ifile, []string{".xlsx"}, "", []string{}) {
	// 				tx := xmfr.NewXMFR()
	// 				if err := tx.LoadXlsx(ifile); err != nil {
	// 					return eris.Wrapf(err, "failed to load %s", ifile)
	// 				}
	// 				if err := tx.SaveToYaml(ofile); err != nil {
	// 					return eris.Wrapf(err, "failed output to %s", ofile)
	// 				}
	// 				return nil
	// 			}
	// 			return eris.Errorf("failed to convert %s to %s", ifile, ofile)
	// 		},
	// 	}
	// }())

	// convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
	// 	var args xmfr.ExcelArgs
	// 	// var ifile, ofile string
	// 	return &cli.Command{
	// 		Name:    "excel",
	// 		Aliases: []string{"e"},
	// 		Flags: []cli.Flag{
	// 			ifileFlag(&args.InFile),
	// 			ofileFlag(&args.OutFile, "Output file (.xlsx)"),
	// 			&cli.BoolFlag{
	// 				Name:        "simple",
	// 				Usage:       "simple output, only PK & FK",
	// 				Value:       false,
	// 				Required:    false,
	// 				Destination: &args.Simple,
	// 			},
	// 		},
	// 		Action: func(c *cli.Context) error {
	// 			if validInOutFile(args.InFile, []string{".yml", ".yaml"}, args.OutFile, []string{".xlsx"}) {
	// 				tx := xmfr.NewXMFR()
	// 				tx.LoadYaml(args.InFile)
	// 				if err := tx.SaveToXlsx(args); err != nil {
	// 					return eris.Wrapf(err, "failed output to %s", args.OutFile)
	// 				}
	// 				return nil
	// 			}
	// 			if validInOutFile(args.InFile, []string{".xlsx"}, args.OutFile, []string{".yml", ".yaml"}) {
	// 				tx := xmfr.NewXMFR()
	// 				tx.LoadXlsx(args.InFile)
	// 				if err := tx.SaveToYaml(args.OutFile); err != nil {
	// 					return eris.Wrapf(err, "failed output to %s", args.OutFile)
	// 				}
	// 				return nil
	// 			}
	// 			return eris.Errorf("failed to convert %s to %s", args.InFile, args.OutFile)
	// 		},
	// 	}
	// }())

	// convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
	// 	// var ifile, ofile, tfile string
	// 	var args xmfr.TextArgs
	// 	return &cli.Command{
	// 		Name:    "text",
	// 		Aliases: []string{"t"},
	// 		Flags: []cli.Flag{
	// 			ifileFlag(&args.InFile),
	// 			ofileFlag(&args.OutFile, "Output file ('stdout' output to console)"),
	// 			&cli.StringFlag{
	// 				Name:        "template",
	// 				Aliases:     []string{"t"},
	// 				Usage:       "Template file",
	// 				Required:    true,
	// 				Destination: &args.TemplateFile,
	// 			},
	// 			&cli.StringFlag{
	// 				Name:        "pattern",
	// 				Aliases:     []string{"p"},
	// 				Usage:       "Table name pattern, e.g. table*",
	// 				Value:       "",
	// 				Required:    false,
	// 				Destination: &args.Pattern,
	// 			},
	// 		},
	// 		Action: func(c *cli.Context) error {
	// 			if validInOutFile(args.InFile, []string{".yml", ".yaml"}, "", []string{}) {
	// 				tx := xmfr.NewXMFR()
	// 				tx.LoadYaml(args.InFile)
	// 				if err := tx.SaveToText(args); err != nil {
	// 					return eris.Wrapf(err, "failed output to %s", args.OutFile)
	// 				}
	// 				return nil
	// 			}
	// 			return eris.Errorf("failed to convert %s to %s", args.InFile, args.OutFile)
	// 		},
	// 	}
	// }())

	// convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
	// 	var args xmfr.DiagramArgs
	// 	return &cli.Command{
	// 		Name:    "diagram",
	// 		Aliases: []string{"d"},
	// 		Flags: []cli.Flag{
	// 			ifileFlag(&args.InFile),
	// 			ofileFlag(&args.OutFile, "Output file (.wsd, .pu, .puml, .plantuml, .iuml)"),
	// 			&cli.StringFlag{
	// 				Name:        "type",
	// 				Aliases:     []string{"t"},
	// 				Usage:       "Diagram type (ER)",
	// 				Value:       "ER",
	// 				Required:    false,
	// 				Destination: &args.DigType,
	// 			},
	// 			&cli.StringFlag{
	// 				Name:        "schema",
	// 				Aliases:     []string{"s"},
	// 				Usage:       "Schema name",
	// 				Value:       "",
	// 				Required:    false,
	// 				Destination: &args.Schema,
	// 			},
	// 			&cli.StringFlag{
	// 				Name:        "pattern",
	// 				Aliases:     []string{"p"},
	// 				Usage:       "Table name pattern, e.g. table*",
	// 				Value:       "",
	// 				Required:    false,
	// 				Destination: &args.Pattern,
	// 			},
	// 			&cli.StringFlag{
	// 				Name:        "jar",
	// 				Aliases:     []string{"j"},
	// 				Usage:       "plantuml.jar file",
	// 				Value:       "",
	// 				Required:    false,
	// 				Destination: &args.JarFile,
	// 			},
	// 			&cli.BoolFlag{
	// 				Name:        "simple",
	// 				Usage:       "simple output, only PK & FK",
	// 				Value:       false,
	// 				Required:    false,
	// 				Destination: &args.Simple,
	// 			},
	// 			&cli.BoolFlag{
	// 				Name:        "fk",
	// 				Usage:       "include foreign key name in the line",
	// 				Value:       false,
	// 				Required:    false,
	// 				Destination: &args.IncludeFK,
	// 			},
	// 		},
	// 		Action: func(c *cli.Context) error {
	// 			if validInOutFile(args.InFile, []string{".yml", ".yaml"}, args.OutFile, []string{".wsd", ".pu", ".puml", ".plantuml", ".iuml"}) {
	// 				tx := xmfr.NewXMFR()
	// 				tx.LoadYaml(args.InFile)
	// 				if err := tx.SaveToPlantUML(args); err != nil {
	// 					return eris.Wrapf(err, "failed output to %s", args.OutFile)
	// 				}
	// 				if args.JarFile != "" {
	// 					if err := sh.Command("java", "-jar", args.JarFile, args.OutFile).Run(); err != nil {
	// 						return eris.Wrapf(err, "failed to generate the diagram")
	// 					}
	// 				}
	// 				return nil
	// 			}
	// 			return eris.Errorf("failed to convert %s to %s", args.InFile, args.OutFile)
	// 		},
	// 	}
	// }())

	// cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
	// 	var ifile string
	// 	return &cli.Command{
	// 		Name:    "verify",
	// 		Aliases: []string{"v"},
	// 		Usage:   "Verify the foreign key",
	// 		Flags: []cli.Flag{
	// 			ifileFlag(&ifile),
	// 		},
	// 		Action: func(c *cli.Context) error {
	// 			tx := xmfr.NewXMFR()
	// 			tx.LoadYaml(ifile)
	// 			return tx.Verify()
	// 		},
	// 	}
	// }())

	// if err := cliapp.Run(os.Args); err != nil {
	// 	logrus.Error(eris.ToString(err, debug))
	// }
}
