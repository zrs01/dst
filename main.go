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

	// inFileFlagBuilder func() *urafvcli.StringFlagBuilder
	ischFlagBuilder   func() *urafvcli.StringFlagBuilder
	outputFlagBuilder func() *urafvcli.StringFlagBuilder
	tmplFlagBuilder   func() *urafvcli.StringFlagBuilder
	schemaFlagBuilder func() *urafvcli.StringFlagBuilder
	tableFlagBuilder  func() *urafvcli.StringFlagBuilder
)

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

	// inFileFlagBuilder = func() *urafvcli.StringFlagBuilder {
	// 	return urafvcli.NewStringFlagBuilder("input").WithAliases("i").WithUsage("input file")
	// }
	outputFlagBuilder = func() *urafvcli.StringFlagBuilder {
		return urafvcli.NewStringFlagBuilder("output").WithAliases("o").WithUsage("output file")
	}
	ischFlagBuilder = func() *urafvcli.StringFlagBuilder {
		return urafvcli.NewStringFlagBuilder("input").WithAliases("i").WithUsage("input schema.yml file")
	}
	tmplFlagBuilder = func() *urafvcli.StringFlagBuilder {
		return urafvcli.NewStringFlagBuilder("template").WithAliases("t").WithUsage("template file")
	}
	schemaFlagBuilder = func() *urafvcli.StringFlagBuilder {
		return urafvcli.NewStringFlagBuilder("schema").WithUsage("schema name pattern, wildcard char: * or %")
	}
	tableFlagBuilder = func() *urafvcli.StringFlagBuilder {
		return urafvcli.NewStringFlagBuilder("table").WithUsage("table name pattern, wildcard char: * or %")
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
