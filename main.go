package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"

	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/urafvcli"
)

var (
	version = "development"

	inFileFlagBuilder   *urafvcli.StringFlagBuilder
	inSchemaFlagBuilder *urafvcli.StringFlagBuilder
	outFileFlagBuilder  *urafvcli.StringFlagBuilder
	tmplFileFlagBuilder *urafvcli.StringFlagBuilder
	schemaFlagBuilder   *urafvcli.StringFlagBuilder
	tableFlagBuilder    *urafvcli.StringFlagBuilder
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

	inFileFlagBuilder = urafvcli.NewStringFlagBuilder("input").WithAliases("i").WithUsage("input file")
	outFileFlagBuilder = urafvcli.NewStringFlagBuilder("output").WithAliases("o").WithUsage("output file")
	inSchemaFlagBuilder = urafvcli.NewStringFlagBuilder("input").WithAliases("i").WithValue("schema.yml")
	tmplFileFlagBuilder = urafvcli.NewStringFlagBuilder("template").WithAliases("t").WithUsage("template file").WithValue("template.yml")
	schemaFlagBuilder = urafvcli.NewStringFlagBuilder("schema").WithUsage("schema name pattern, wildcard char: * or %")
	tableFlagBuilder = urafvcli.NewStringFlagBuilder("table").WithUsage("table name pattern, wildcard char: * or %")

	// ifileFlag = func(file *string, usage string) *cli.StringFlag {
	// 	return &cli.StringFlag{Name: "input", Aliases: []string{"i"}, Usage: lo.Ternary(usage == "", "input file", usage), Required: true, Destination: file}
	// }
	// iSchemaFileFlag = func(file *string) *cli.StringFlag {
	// 	// if schema.yml at current folder, use it as default
	// 	flag := inFileFlagBuilder.WithUsage("input file (.yml)").WithDestination(file).Build()
	// 	if _, err := os.Stat("schema.yml"); !os.IsNotExist(err) {
	// 		flag.Value = "schema.yml"
	// 		flag.Required = false
	// 	}
	// 	return flag
	// }
	// ofileFlag = func(file *string, usage string) *cli.StringFlag {
	// 	return &cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: lo.Ternary(usage == "", "output file", usage), Required: false, Destination: file}
	// }
	// templateFlag = func(file *string) *cli.StringFlag {
	// 	return &cli.StringFlag{Name: "template", Aliases: []string{"t"}, Usage: "template file", Required: false, Destination: file}
	// }
	// schemaFile = func(schema *string) *cli.StringFlag {
	// 	return &cli.StringFlag{Name: "schema", Usage: "schema name pattern, wildcard char: * or %", Required: false, Destination: schema}
	// }
	// tableFlag = func(table *string) *cli.StringFlag {
	// 	return &cli.StringFlag{Name: "table", Usage: "table name pattern, wildcard char: * or %", Required: false, Destination: table}
	// }
	// simpleFlag = func(simple *bool) *cli.BoolFlag {
	// 	return &cli.BoolFlag{Name: "simple", Usage: "simple content", Value: false, Required: false, Destination: simple}
	// }
	// libFlag = func(lib *string) *cli.StringFlag {
	// 	return &cli.StringFlag{Name: "lib", Usage: "plantuml.jar file, used when output format is png", Required: false, Destination: lib}
	// }
	// dumpFlag = func(dump *bool) *cli.BoolFlag {
	// 	return &cli.BoolFlag{Name: "dump", Usage: "dump content", Value: false, Required: false, Destination: dump}
	// }

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
