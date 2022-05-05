package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rotisserie/eris"
	"github.com/urfave/cli/v2"

	xmfr "dst/xfmr"
)

func main() {

	cliapp := cli.NewApp()
	cliapp.Name = "dst"
	cliapp.Usage = "Database schema tool"
	cliapp.Version = "0.0.1-202204"
	cliapp.Commands = []*cli.Command{}

	debug := false
	var infile, outfile, outtpl string

	// global options
	cliapp.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Usage:       "Debug mode",
			Required:    false,
			Destination: &debug,
		},
		&cli.StringFlag{
			Name:        "input",
			Aliases:     []string{"i"},
			Usage:       "Input file (source)",
			Required:    true,
			Destination: &infile,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "Output file",
			Required:    true,
			Destination: &outfile,
		},
		&cli.StringFlag{
			Name:        "template",
			Aliases:     []string{"t"},
			Usage:       "Template file",
			Required:    false,
			Destination: &outtpl,
		},
	}
	cliapp.Action = func(c *cli.Context) error {
		inext := filepath.Ext(infile)
		outext := filepath.Ext(outfile)
		if inext == ".yml" || inext == ".yaml" {
			if outext == ".xlsx" {
				tx := xmfr.NewXMFR()
				if err := tx.YamlToExcel(infile, outfile); err != nil {
					return eris.Wrapf(err, "failed to output the file %s", outfile)
				}
			} else if outtpl != "" {
				tx := xmfr.NewXMFR()
				if err := tx.YamlToText(infile, outfile, outtpl); err != nil {
					return eris.Wrapf(err, "failed to output the file %s with template %s", outfile, outtpl)
				}
			} else {
				return eris.Errorf("Output file type '%s' is not supported", outext)
			}
		} else if inext == ".xlsx" {
			if outext == ".yml" || outext == ".yaml" {
				tx := xmfr.NewXMFR()
				if err := tx.ExcelToYaml(infile, outfile); err != nil {
					return eris.Wrapf(err, "failed to output the file %s", outfile)
				}
			} else {
				return eris.Errorf("Output file type '%s' is not supported", outext)
			}
		} else {
			return eris.Errorf("Input file type '%s' is not supported", inext)
		}
		return nil
	}

	// cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
	// 	var file, exportType, outfile string
	// 	return &cli.Command{
	// 		Name:  "export",
	// 		Usage: "Export to other format",
	// 		Flags: []cli.Flag{
	// 			&cli.StringFlag{
	// 				Name:        "file",
	// 				Aliases:     []string{"f"},
	// 				Usage:       "Definition file",
	// 				Required:    true,
	// 				Destination: &file,
	// 			},
	// 			&cli.StringFlag{
	// 				Name:        "type",
	// 				Aliases:     []string{"t"},
	// 				Usage:       "Export type, 'xlsx' (default) or 'sql'",
	// 				Required:    false,
	// 				Value:       "xlsx",
	// 				Destination: &exportType,
	// 			},
	// 			&cli.StringFlag{
	// 				Name:        "output",
	// 				Aliases:     []string{"o"},
	// 				Usage:       "Output file",
	// 				Required:    true,
	// 				Destination: &outfile,
	// 			},
	// 		},
	// 		Action: func(c *cli.Context) error {
	// 			db := db.NewDatabase()
	// 			if err := db.Load(file); err != nil {
	// 				return eris.Wrapf(err, "failed to load the definition file %s", file)
	// 			}
	// 			if exportType == "xlsx" {
	// 				if err := db.ExportToExcel(outfile); err != nil {
	// 					return eris.Wrapf(err, "failed to output the file %s", outfile)
	// 				}
	// 			} else {
	// 				return eris.Errorf("Export type '%s' is not supported", exportType)
	// 			}
	// 			return nil
	// 		},
	// 	}
	// }())

	// cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
	// 	var file, importType, outfile string
	// 	return &cli.Command{
	// 		Name:  "import",
	// 		Usage: "Import from other format",
	// 		Flags: []cli.Flag{
	// 			&cli.StringFlag{
	// 				Name:        "file",
	// 				Aliases:     []string{"f"},
	// 				Usage:       "Import file",
	// 				Required:    true,
	// 				Destination: &file,
	// 			},
	// 			&cli.StringFlag{
	// 				Name:        "type",
	// 				Aliases:     []string{"t"},
	// 				Usage:       "Import type, 'xlsx' (default) or 'sql'",
	// 				Required:    false,
	// 				Value:       "xlsx",
	// 				Destination: &importType,
	// 			},
	// 			&cli.StringFlag{
	// 				Name:        "output",
	// 				Aliases:     []string{"o"},
	// 				Usage:       "Output file",
	// 				Required:    true,
	// 				Destination: &outfile,
	// 			},
	// 		},
	// 		Action: func(c *cli.Context) error {
	// 			db := db.NewDatabase()
	// 			if importType == "xlsx" {
	// 				if err := db.ImportFromExcel(file, outfile); err != nil {
	// 					return eris.Wrapf(err, "failed to output the file %s", outfile)
	// 				}
	// 			} else {
	// 				return eris.Errorf("Import type '%s' is not supported", importType)
	// 			}
	// 			return nil
	// 		},
	// 	}
	// }())

	if err := cliapp.Run(os.Args); err != nil {
		fmt.Println(eris.ToString(err, debug))
	}
}
