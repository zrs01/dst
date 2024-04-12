package main

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"

	"github.com/zrs01/dst/config"
)

var (
	version         = "development"
	iSchemaFileFlag func(file *string) *cli.StringFlag
	ifileFlag       func(file *string, usage string) *cli.StringFlag
	ofileFlag       func(file *string, usage string) *cli.StringFlag
	templateFlag    func(file *string) *cli.StringFlag
	schemaFile      func(schema *string) *cli.StringFlag
	tableFlag       func(table *string) *cli.StringFlag
	simpleFlag      func(simple *bool) *cli.BoolFlag
	libFlag         func(lib *string) *cli.StringFlag
	dumpFlag        func(dump *bool) *cli.BoolFlag
)

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

	// debug := false

	// global options
	cliapp.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Usage:       "Debug mode",
			Required:    false,
			Destination: &config.Debug,
		},
	}

	/* ------------------------------ Common flags ------------------------------ */

	ifileFlag = func(file *string, usage string) *cli.StringFlag {
		return &cli.StringFlag{Name: "input", Aliases: []string{"i"}, Usage: lo.Ternary(usage == "", "input file", usage), Required: true, Destination: file}
	}
	iSchemaFileFlag = func(file *string) *cli.StringFlag {
		// if schema.yml at current folder, use it as default
		flag := ifileFlag(file, "input file (.yml)")
		if _, err := os.Stat("schema.yml"); !os.IsNotExist(err) {
			flag.Value = "schema.yml"
			flag.Required = false
		}
		return flag
	}
	ofileFlag = func(file *string, usage string) *cli.StringFlag {
		return &cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: lo.Ternary(usage == "", "output file", usage), Required: false, Destination: file}
	}
	templateFlag = func(file *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "template", Aliases: []string{"t"}, Usage: "template file", Required: false, Destination: file}
	}
	schemaFile = func(schema *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "schema", Usage: "schema name pattern, wildcard char: * or %", Required: false, Destination: schema}
	}
	tableFlag = func(table *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "table", Usage: "table name pattern, wildcard char: * or %", Required: false, Destination: table}
	}
	simpleFlag = func(simple *bool) *cli.BoolFlag {
		return &cli.BoolFlag{Name: "simple", Usage: "simple content", Value: false, Required: false, Destination: simple}
	}
	libFlag = func(lib *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "lib", Usage: "plantuml.jar file, used when output format is png", Required: false, Destination: lib}
	}
	dumpFlag = func(dump *bool) *cli.BoolFlag {
		return &cli.BoolFlag{Name: "dump", Usage: "dump content", Value: false, Required: false, Destination: dump}
	}

	registerCmdConvert(cliapp)
	registerCmdSql(cliapp)

	if err := cliapp.Run(os.Args); err != nil {
		if config.Debug {
			tracerr.Print(err)
		} else {
			fmt.Printf("Error: %s\n", err)
		}
	}
}
