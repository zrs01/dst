package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"

	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/sql"
	"github.com/zrs01/dst/transform"
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

	ifileFlag := func(file *string, usage string) *cli.StringFlag {
		return &cli.StringFlag{Name: "input", Aliases: []string{"i"}, Usage: lo.Ternary(usage == "", "input file", usage), Required: true, Destination: file}
	}
	ofileFlag := func(file *string, usage string) *cli.StringFlag {
		return &cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: lo.Ternary(usage == "", "output file", usage), Required: false, Destination: file}
	}
	templateFlag := func(file *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "template", Aliases: []string{"t"}, Usage: "template file", Required: false, Destination: file}
	}
	schemaFile := func(schema *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "schema", Usage: "schema name pattern, wildcard char: * or %", Required: false, Destination: schema}
	}
	tableFlag := func(table *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "table", Usage: "table name pattern, wildcard char: * or %", Required: false, Destination: table}
	}
	simpleFlag := func(simple *bool) *cli.BoolFlag {
		return &cli.BoolFlag{Name: "simple", Usage: "simple content", Value: false, Required: false, Destination: simple}
	}
	libFlag := func(lib *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "lib", Usage: "plantuml.jar file, used when output format is png", Required: false, Destination: lib}
	}
	selectedData := func(ifile, schema, table string, column string) (*model.DataDef, error) {
		rawData, err := transform.ReadYml(ifile)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		validateResult := transform.Verify(rawData)
		if len(validateResult) > 0 {
			lo.ForEach(validateResult, func(v string, _ int) {
				fmt.Println(v)
			})
			return nil, tracerr.Errorf("invalid data")
		}
		data, err := transform.FilterData(rawData, schema, table, column)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		return data, nil
	}

	// convert command
	convertCmd := &cli.Command{
		Name:    "convert",
		Aliases: []string{"c"},
		Usage:   "Convert to other format",
	}
	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		return convertCmd
	}())

	// transform to text
	convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, tfile, schema, table string
		var simple bool
		return &cli.Command{
			Name:    "text",
			Usage:   "transform from yaml to text",
			Aliases: []string{"t"},
			Flags: []cli.Flag{
				ifileFlag(&ifile, "input file (.yml)"),
				ofileFlag(&ofile, "output file (text file)"),
				schemaFile(&schema),
				tableFlag(&table),
				simpleFlag(&simple),
				templateFlag(&tfile),
			},
			Action: func(c *cli.Context) error {
				oext := lo.Ternary(ofile != "", strings.ToLower(filepath.Ext(ofile)), "")
				data, err := selectedData(ifile, schema, table, "")
				if err != nil {
					return tracerr.Wrap(err)
				}
				if tfile != "" {
					return transform.WriteFileTpl(data, tfile, ofile)
				}
				switch oext {
				case ".yml", "":
					if err := transform.WriteYml(data, ofile); err != nil {
						return tracerr.Wrap(err)
					}
					return nil
				}
				return tracerr.New("Not implemented yet")
			},
		}
	}())

	// transform to excel
	convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table string
		var simple bool
		return &cli.Command{
			Name:    "excel",
			Usage:   "transform from yaml to excel",
			Aliases: []string{"e"},
			Flags: []cli.Flag{
				ifileFlag(&ifile, "input file (.yml)"),
				ofileFlag(&ofile, "output file (.xlsx)"),
				schemaFile(&schema),
				tableFlag(&table),
				simpleFlag(&simple),
			},
			Action: func(c *cli.Context) error {
				oext := lo.Ternary(ofile != "", strings.ToLower(filepath.Ext(ofile)), "")
				data, err := selectedData(ifile, schema, table, "")
				if err != nil {
					return tracerr.Wrap(err)
				}
				switch oext {
				case ".xlsx":
					if err := transform.WriteXlsx(data, ofile, simple); err != nil {
						return tracerr.Wrap(err)
					}
					return nil
				}
				return tracerr.New("Not implemented yet")
			},
		}
	}())

	// transform to diagram
	convertCmd.Subcommands = append(convertCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, tfile, schema, table, lib string
		return &cli.Command{
			Name:    "diagram",
			Usage:   "transform from yaml to diagram",
			Aliases: []string{"d"},
			Flags: []cli.Flag{
				ifileFlag(&ifile, "input file (.yml)"),
				ofileFlag(&ofile, "output file (.png)"),
				schemaFile(&schema),
				tableFlag(&table),
				templateFlag(&tfile),
				libFlag(&lib),
			},
			Action: func(c *cli.Context) error {
				if ofile == "" {
					ofile = strings.TrimSuffix(ifile, filepath.Ext(ifile)) + ".png"
				}
				oext := lo.Ternary(ofile != "", strings.ToLower(filepath.Ext(ofile)), "")
				data, err := selectedData(ifile, schema, table, "")
				if err != nil {
					return tracerr.Wrap(err)
				}
				switch oext {
				case ".png":
					if err := transform.WriteERD(data, tfile, ofile); err != nil {
						return tracerr.Wrap(err)
					}
					return nil
				}
				return tracerr.New(fmt.Sprintf("output file extension '%s' is not supported", ofile))
			},
		}
	}())

	sqlCmd := &cli.Command{
		Name:    "sql",
		Aliases: []string{"s"},
		Usage:   "Convert to SQL DDL",
	}
	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		return sqlCmd
	}())

	/* -------------------------------------------------------------------------- */
	/*                                     SQL                                    */
	/* -------------------------------------------------------------------------- */

	dbFlag := func(db *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "database", Aliases: []string{"d"}, Usage: "database (mssql)", Required: true, Destination: db}
	}
	colFlag := func(col *string) *cli.StringFlag {
		return &cli.StringFlag{Name: "column", Aliases: []string{"c"}, Usage: "column", Required: false, Destination: col}
	}
	ifileWithDefaultFlag := func(file *string) *cli.StringFlag {
		// if schema.yml at current folder, use it as default
		flag := ifileFlag(file, "input file (.yml)")
		if _, err := os.Stat("schema.yml"); !os.IsNotExist(err) {
			flag.Value = "schema.yml"
			flag.Required = false
		}
		return flag
	}

	sqlCmd.Subcommands = append(sqlCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table, db string
		return &cli.Command{
			Name:    "create-table",
			Usage:   "create table DDL",
			Aliases: []string{"ct"},
			Flags: []cli.Flag{
				ifileWithDefaultFlag(&ifile),
				ofileFlag(&ofile, "output file"),
				schemaFile(&schema),
				tableFlag(&table),
				dbFlag(&db),
			},
			Action: func(c *cli.Context) error {
				data, err := selectedData(ifile, schema, table, "")
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.CreateTable(data, db, ofile)
			},
		}
	}())

	sqlCmd.Subcommands = append(sqlCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table, db string
		return &cli.Command{
			Name:    "drop-table",
			Usage:   "drop table DDL",
			Aliases: []string{"dt"},
			Flags: []cli.Flag{
				ifileWithDefaultFlag(&ifile),
				ofileFlag(&ofile, "output file"),
				schemaFile(&schema),
				tableFlag(&table),
				dbFlag(&db),
			},
			Action: func(c *cli.Context) error {
				data, err := selectedData(ifile, schema, table, "")
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.DropTable(data, db, ofile)
			},
		}
	}())

	sqlCmd.Subcommands = append(sqlCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:    "add-column",
			Usage:   "add column DDL",
			Aliases: []string{"ac"},
			Flags: []cli.Flag{
				ifileWithDefaultFlag(&ifile),
				ofileFlag(&ofile, "output file"),
				schemaFile(&schema),
				tableFlag(&table),
				dbFlag(&db),
				colFlag(&col),
			},
			Action: func(c *cli.Context) error {
				data, err := selectedData(ifile, schema, table, col)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.AddColumn(data, db, ofile)
			},
		}
	}())

	sqlCmd.Subcommands = append(sqlCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:    "drop-column",
			Usage:   "drop column DDL",
			Aliases: []string{"dc"},
			Flags: []cli.Flag{
				ifileWithDefaultFlag(&ifile),
				ofileFlag(&ofile, "output file"),
				schemaFile(&schema),
				tableFlag(&table),
				dbFlag(&db),
				colFlag(&col),
			},
			Action: func(c *cli.Context) error {
				data, err := selectedData(ifile, schema, table, col)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.DropColumn(data, db, ofile)
			},
		}
	}())

	sqlCmd.Subcommands = append(sqlCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col, newcol string
		cf := colFlag(&col)
		cf.Required = true
		return &cli.Command{
			Name:    "rename-column",
			Usage:   "rename column DDL",
			Aliases: []string{"rc"},
			Flags: []cli.Flag{
				ifileWithDefaultFlag(&ifile),
				ofileFlag(&ofile, "output file"),
				schemaFile(&schema),
				tableFlag(&table),
				dbFlag(&db),
				cf,
				&cli.StringFlag{Name: "new-column", Aliases: []string{"n"}, Usage: "new column name", Required: true, Destination: &newcol},
			},
			Action: func(c *cli.Context) error {
				data, err := selectedData(ifile, schema, table, col)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.RenameColumn(data, db, ofile, newcol)
			},
		}
	}())

	sqlCmd.Subcommands = append(sqlCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:    "modify-column",
			Usage:   "modify column type DDL",
			Aliases: []string{"mc"},
			Flags: []cli.Flag{
				ifileWithDefaultFlag(&ifile),
				ofileFlag(&ofile, "output file"),
				schemaFile(&schema),
				tableFlag(&table),
				dbFlag(&db),
				colFlag(&col),
			},
			Action: func(c *cli.Context) error {
				data, err := selectedData(ifile, schema, table, col)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.ModifyColumn(data, db, ofile)
			},
		}
	}())

	if err := cliapp.Run(os.Args); err != nil {
		if debug {
			tracerr.Print(err)
		} else {
			fmt.Printf("Error: %s\n", err)
		}
	}
}
