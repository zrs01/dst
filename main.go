package main

import (
	"fmt"
	"os"

	"github.com/rotisserie/eris"
	"github.com/urfave/cli/v2"
)

func main() {
	debug := false

	cliapp := cli.NewApp()
	cliapp.Name = "dst"
	cliapp.Usage = "Database schema tool"
	cliapp.Version = "0.0.1-202204"
	cliapp.Commands = []*cli.Command{}

	// global options
	cliapp.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Usage:       "enable debug mode",
			Required:    false,
			Destination: &debug,
		},
	}

	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		var file, exportType, outfile string
		return &cli.Command{
			Name:  "export",
			Usage: "Export to other format",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "file",
					Aliases:     []string{"f"},
					Usage:       "Definition file",
					Required:    true,
					Destination: &file,
				},
				&cli.StringFlag{
					Name:        "type",
					Aliases:     []string{"t"},
					Usage:       "Export type, 'xlsx' (default) or 'sql'",
					Required:    false,
					Value:       "xlsx",
					Destination: &exportType,
				},
				&cli.StringFlag{
					Name:        "output",
					Aliases:     []string{"o"},
					Usage:       "Output file",
					Required:    true,
					Destination: &outfile,
				},
			},
			Action: func(c *cli.Context) error {
				db := NewDatabase()
				if err := db.Load(file); err != nil {
					return eris.Wrapf(err, "failed to load the definition file %s", file)
				}
				if exportType == "xlsx" {
					if err := db.ExportToExcel(outfile); err != nil {
						return eris.Wrapf(err, "failed to output the file %s", outfile)
					}
				} else {
					return eris.Errorf("Export type '%s' is not supported", exportType)
				}
				return nil
			},
		}
	}())
	if err := cliapp.Run(os.Args); err != nil {
		fmt.Println(eris.ToString(err, debug))
	}
}
