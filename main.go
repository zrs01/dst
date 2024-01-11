package main

import (
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

	ifileFlag := func(file *string) *cli.StringFlag {
		return &cli.StringFlag{
			Name:        "input",
			Aliases:     []string{"i"},
			Usage:       "Input file",
			Required:    true,
			Destination: file,
		}
	}
	ofileFlag := func(file *string, usage string) *cli.StringFlag {
		if usage == "" {
			usage = "Output file"
		}
		return &cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       usage,
			Destination: file,
		}
	}
	// validInOutFile := func(ifile string, iexts []string, ofile string, oexts []string) bool {
	// 	iext := strings.ToLower(filepath.Ext(ifile))
	// 	oext := strings.ToLower(filepath.Ext(ofile))
	// 	if ifile != "" && ofile != "" {
	// 		return lo.Contains(iexts, iext) && lo.Contains(oexts, oext)
	// 	}
	// 	return lo.Contains(iexts, iext)
	// }

	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		var ifile, ofile, tfile string
		return &cli.Command{
			Name:    "transform",
			Aliases: []string{"t"},
			Usage:   "Transform to other format",
			Flags: []cli.Flag{
				ifileFlag(&ifile),
				ofileFlag(&ofile, ""),

				&cli.StringFlag{
					Name:        "template",
					Aliases:     []string{"t"},
					Usage:       "Template file",
					Destination: &tfile,
				},
			},
			Action: func(c *cli.Context) error {
				iext := strings.ToLower(filepath.Ext(ifile))
				oext := lo.Ternary(ofile != "", strings.ToLower(filepath.Ext(ofile)), "")

				switch iext {
				case ".yml":
					data, err := transform.ReadYml(ifile)
					if err != nil {
						return tracerr.Wrap(err)
					}

					// text output from template
					if tfile != "" {

					}

					switch oext {
					case ".yml", "":
						if err := transform.WriteYml(data, ofile); err != nil {
							return tracerr.Wrap(err)
						}
					}
				}
				return nil
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
