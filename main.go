package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/codeskyblue/go-sh"
	"github.com/rotisserie/eris"
	"github.com/shomali11/util/xstrings"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"

	xmfr "dst/xfmr"
)

var version = "development"

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
		if xstrings.IsBlank(usage) {
			usage = "Output file"
		}
		return &cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       usage,
			Required:    true,
			Destination: file,
		}
	}
	validInOutFile := func(ifile string, iexts []string, ofile string, oexts []string) bool {
		iext := strings.ToLower(filepath.Ext(ifile))
		oext := strings.ToLower(filepath.Ext(ofile))
		if xstrings.IsNotBlank(ifile) && xstrings.IsNotBlank(ofile) {
			return funk.Contains(iexts, iext) && funk.Contains(oexts, oext)
		}
		return funk.Contains(iexts, iext)
	}

	convertCmd := &cli.Command{
		Name:    "convert",
		Aliases: []string{"c"},
		Usage:   "Convert to other format",
	}
	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		return convertCmd
	}())

	convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
		var ifile, ofile string
		return &cli.Command{
			Name:    "yaml",
			Aliases: []string{"y"},
			Flags: []cli.Flag{
				ifileFlag(&ifile),
				ofileFlag(&ofile, ""),
			},
			Action: func(c *cli.Context) error {
				if validInOutFile(ifile, []string{".xlsx"}, "", []string{}) {
					tx := xmfr.NewXMFR()
					if err := tx.LoadXlsx(ifile); err != nil {
						return eris.Wrapf(err, "failed to load %s", ifile)
					}
					if err := tx.SaveToYaml(ofile); err != nil {
						return eris.Wrapf(err, "failed output to %s", ofile)
					}
					return nil
				}
				return eris.Errorf("failed to convert %s to %s", ifile, ofile)
			},
		}
	}())

	convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
		var ifile, ofile string
		return &cli.Command{
			Name:    "excel",
			Aliases: []string{"e"},
			Flags: []cli.Flag{
				ifileFlag(&ifile),
				ofileFlag(&ofile, "Output file (.xlsx)"),
			},
			Action: func(c *cli.Context) error {
				if validInOutFile(ifile, []string{".yml", ".yaml"}, ofile, []string{".xlsx"}) {
					tx := xmfr.NewXMFR()
					tx.LoadYaml(ifile)
					if err := tx.SaveToXlsx(ofile); err != nil {
						return eris.Wrapf(err, "failed output to %s", ofile)
					}
					return nil
				}
				if validInOutFile(ifile, []string{".xlsx"}, ofile, []string{".yml", ".yaml"}) {
					tx := xmfr.NewXMFR()
					tx.LoadXlsx(ifile)
					if err := tx.SaveToYaml(ofile); err != nil {
						return eris.Wrapf(err, "failed output to %s", ofile)
					}
					return nil
				}
				return eris.Errorf("failed to convert %s to %s", ifile, ofile)
			},
		}
	}())

	convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, tfile string
		return &cli.Command{
			Name:    "text",
			Aliases: []string{"t"},
			Flags: []cli.Flag{
				ifileFlag(&ifile),
				ofileFlag(&ofile, "Output file ('stdout' output to console)"),
				&cli.StringFlag{
					Name:        "template",
					Aliases:     []string{"t"},
					Usage:       "Template file",
					Required:    true,
					Destination: &tfile,
				},
			},
			Action: func(c *cli.Context) error {
				if validInOutFile(ifile, []string{".yml", ".yaml"}, "", []string{}) {
					tx := xmfr.NewXMFR()
					tx.LoadYaml(ifile)
					if err := tx.SaveToText(ofile, tfile); err != nil {
						return eris.Wrapf(err, "failed output to %s", ofile)
					}
					return nil
				}
				return eris.Errorf("failed to convert %s to %s", ifile, ofile)
			},
		}
	}())

	convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
		var args xmfr.DiagramArgs
		return &cli.Command{
			Name:    "diagram",
			Aliases: []string{"d"},
			Flags: []cli.Flag{
				ifileFlag(&args.InFile),
				ofileFlag(&args.OutFile, "Output file (.wsd, .pu, .puml, .plantuml, .iuml)"),
				&cli.StringFlag{
					Name:        "type",
					Aliases:     []string{"t"},
					Usage:       "Diagram type (ER)",
					Value:       "ER",
					Required:    false,
					Destination: &args.DigType,
				},
				&cli.StringFlag{
					Name:        "schema",
					Aliases:     []string{"s"},
					Usage:       "Schema name",
					Value:       "",
					Required:    false,
					Destination: &args.Schema,
				},
				&cli.StringFlag{
					Name:        "prefix",
					Aliases:     []string{"p"},
					Usage:       "Table prefix",
					Value:       "",
					Required:    false,
					Destination: &args.TablePrefix,
				},
				&cli.StringFlag{
					Name:        "jar",
					Aliases:     []string{"j"},
					Usage:       "plantuml.jar file",
					Value:       "",
					Required:    false,
					Destination: &args.JarFile,
				},
			},
			Action: func(c *cli.Context) error {
				if validInOutFile(args.InFile, []string{".yml", ".yaml"}, args.OutFile, []string{".wsd", ".pu", ".puml", ".plantuml", ".iuml"}) {
					tx := xmfr.NewXMFR()
					tx.LoadYaml(args.InFile)
					if err := tx.SaveToPlantUML(args); err != nil {
						return eris.Wrapf(err, "failed output to %s", args.OutFile)
					}
					if xstrings.IsNotBlank(args.JarFile) {
						if err := sh.Command("java", "-jar", args.JarFile, args.OutFile).Run(); err != nil {
							return eris.Wrapf(err, "failed to generate the diagram")
						}
					}
					return nil
				}
				return eris.Errorf("failed to convert %s to %s", args.InFile, args.OutFile)
			},
		}
	}())

	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		var ifile string
		return &cli.Command{
			Name:    "verify",
			Aliases: []string{"v"},
			Usage:   "Verify the foreign key",
			Flags: []cli.Flag{
				ifileFlag(&ifile),
			},
			Action: func(c *cli.Context) error {
				tx := xmfr.NewXMFR()
				tx.LoadYaml(ifile)
				return tx.Verify()
			},
		}
	}())

	if err := cliapp.Run(os.Args); err != nil {
		logrus.Error(eris.ToString(err, debug))
	}
}
