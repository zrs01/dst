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
	cliapp.Version = "0.0.1-202206"
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

	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		var infile, outfile, outtpl string
		return &cli.Command{
			Name:  "convert",
			Usage: "Convert to other format",
			Flags: []cli.Flag{
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
			},
			Action: func(c *cli.Context) error {
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
			},
		}
	}())

	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		var infile string
		return &cli.Command{
			Name:  "verify",
			Usage: "Verify the foreign key",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "input",
					Aliases:     []string{"i"},
					Usage:       "Input file (source)",
					Required:    true,
					Destination: &infile,
				},
			},
			Action: func(c *cli.Context) error {
				tx := xmfr.NewXMFR()
				return tx.VerifyForeignKey(infile)
			},
		}
	}())

	if err := cliapp.Run(os.Args); err != nil {
		fmt.Println(eris.ToString(err, debug))
	}
}
