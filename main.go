package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"

	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/cliflag"
)

var (
	version = "development"

	schemaFileFlagBuilder   func() *cliflag.StringFlagBuilder // schedule file
	outputFileFlagBuilder   func() *cliflag.StringFlagBuilder // output file
	templateFileFlagBuilder func() *cliflag.StringFlagBuilder // template file
	schemaNameFlagBuilder   func() *cliflag.StringFlagBuilder // schema name
	tableNameFlagBuilder    func() *cliflag.StringFlagBuilder // table name
)

func main() {
	cliapp := cli.NewApp()
	cliapp.Name = "dst"
	cliapp.Usage = "Database schema tool"
	cliapp.Version = version
	cliapp.Commands = []*cli.Command{}

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

	outputFileFlagBuilder = func() *cliflag.StringFlagBuilder {
		return cliflag.NewStringFlagBuilder("output").WithAliases("o").WithUsage("output file")
	}
	schemaFileFlagBuilder = func() *cliflag.StringFlagBuilder {
		return cliflag.NewStringFlagBuilder("input").WithAliases("i").WithUsage("schema file")
	}
	templateFileFlagBuilder = func() *cliflag.StringFlagBuilder {
		return cliflag.NewStringFlagBuilder("template").WithAliases("t").WithUsage("template file")
	}
	schemaNameFlagBuilder = func() *cliflag.StringFlagBuilder {
		return cliflag.NewStringFlagBuilder("schema").WithUsage("schema name pattern, wildcard char: * or %")
	}
	tableNameFlagBuilder = func() *cliflag.StringFlagBuilder {
		return cliflag.NewStringFlagBuilder("table").WithUsage("table name pattern, wildcard char: * or %")
	}

	registerTemplateRenderer(cliapp)
	registerDDLRenderer(cliapp)

	if err := cliapp.Run(os.Args); err != nil {
		if config.Debug {
			tracerr.Print(err)
		} else {
			fmt.Printf("Error: %s\n", err)
		}
	}
}
