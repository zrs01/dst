package main

import (
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"

	nested "github.com/antonfisher/nested-logrus-formatter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
)

var (
	version = "development"

	schemaFileFlagBuilder   func() *flagbuilder.StringFlag // schedule file
	outputFileFlagBuilder   func() *flagbuilder.StringFlag // output file
	templateFileFlagBuilder func() *flagbuilder.StringFlag // template file
	schemaNameFlagBuilder   func() *flagbuilder.StringFlag // schema name
	tableNameFlagBuilder    func() *flagbuilder.StringFlag // table name
)

func main() {
	initLogrus()

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

	outputFileFlagBuilder = func() *flagbuilder.StringFlag {
		return flagbuilder.NewStringFlag("output").WithAliases("o").WithUsage("output file")
	}
	schemaFileFlagBuilder = func() *flagbuilder.StringFlag {
		return flagbuilder.NewStringFlag("input").WithAliases("i").WithUsage(`input file. It can accept various types of input sources:
	1. A schema file in YAML format (e.g. schemal.yml)
	2. A MySQL database connnection string in the format:
	   mysql://[user[:cred]@][protocol[(address[:port])]]/dbname[?param1=value1&...]
	3. A SQL Server database connection string in the format:
	   sqlserver://user:cred@host[:port][/dbname][?param1=value1&...]`)
	}
	templateFileFlagBuilder = func() *flagbuilder.StringFlag {
		return flagbuilder.NewStringFlag("template").WithAliases("t").WithUsage("Template file")
	}
	schemaNameFlagBuilder = func() *flagbuilder.StringFlag {
		return flagbuilder.NewStringFlag("schema").WithUsage("Schema name pattern; wildcard characters allowed: * or %")
	}
	tableNameFlagBuilder = func() *flagbuilder.StringFlag {
		return flagbuilder.NewStringFlag("table").WithUsage("Table name pattern; wildcard characters allowed: * or %")
	}

	RegisterTxt(cliapp) // Template
	RegisterXls(cliapp) // Excel
	RegisterErd(cliapp) // ER diagram
	RegisterSql(cliapp) // SQL DDL statement

	if err := cliapp.Run(os.Args); err != nil {
		if config.Debug {
			logrus.Error(tracerr.SprintSourceColor(err, 0))
		} else {
			logrus.Errorf("%s", err)
		}
	}
}

func initLogrus() {
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: time.RFC3339,
	})
}
