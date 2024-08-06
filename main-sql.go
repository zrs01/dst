package main

import (
	"github.com/urfave/cli/v2"
	"github.com/zrs01/dst/internal/sql"
	"github.com/zrs01/dst/internal/urafvcli"
	"github.com/zrs01/dst/internal/yml"
	"github.com/ztrue/tracerr"
)

/* -------------------------------------------------------------------------- */
/*                                     SQL                                    */
/* -------------------------------------------------------------------------- */
func registerCmdSql(cliapp *cli.App) *cli.Command {
	sqlCmd := &cli.Command{
		Name:    "sql",
		Aliases: []string{"s"},
		Usage:   "Generate SQL DDL",
	}
	cliapp.Commands = append(cliapp.Commands, func() *cli.Command {
		return sqlCmd
	}())

	databaseFlagBuilder := urafvcli.NewStringFlagBuilder("database").WithAliases("d").WithUsage("database (mssql)")
	// dbFlag := func(db *string) *cli.StringFlag {
	// 	return &cli.StringFlag{Name: "database", Aliases: []string{"d"}, Usage: "database (mssql)", Required: true, Destination: db}
	// }
	columnFlagBuilder := urafvcli.NewStringFlagBuilder("column").WithAliases("c").WithUsage("column")
	// colFlag := func(col *string) *cli.StringFlag {
	// 	return &cli.StringFlag{Name: "column", Aliases: []string{"c"}, Usage: "column", Required: false, Destination: col}
	// }

	sqlCmd.Subcommands = append(sqlCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table, db string
		return &cli.Command{
			Name:    "create_table",
			Usage:   "create table DDL",
			Aliases: []string{"ct"},
			Flags: []cli.Flag{
				urafvcli.NewStringFlagBuilder("input").WithAliases("i").WithValue("schema.yml").WithDestination(&ifile).Build(),
				urafvcli.NewStringFlagBuilder("output").WithAliases("o").WithUsage("output file (.png)").WithDestination(&ofile).Build(),
				urafvcli.NewStringFlagBuilder("schema").WithUsage("schema name pattern, wildcard char: * or %").WithDestination(&schema).Build(),
				urafvcli.NewStringFlagBuilder("table").WithUsage("table name pattern, wildcard char: * or %").WithDestination(&table).Build(),
				urafvcli.NewStringFlagBuilder("database").WithAliases("d").WithUsage("database (mssql)").WithDestination(&db).Build(),

				// iSchemaFileFlag(&ifile),
				// ofileFlag(&ofile, "output file"),
				// schemaFile(&schema),
				// tableFlag(&table),
				// dbFlag(&db),
			},
			Action: func(c *cli.Context) error {
				data, err := yml.ReadSelectedYml(ifile, schema, table, "")
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
			Name:    "drop_table",
			Usage:   "drop table DDL",
			Aliases: []string{"dt"},
			Flags: []cli.Flag{
				urafvcli.NewStringFlagBuilder("input").WithAliases("i").WithValue("schema.yml").WithDestination(&ifile).Build(),
				urafvcli.NewStringFlagBuilder("output").WithAliases("o").WithUsage("output file (.png)").WithDestination(&ofile).Build(),
				urafvcli.NewStringFlagBuilder("schema").WithUsage("schema name pattern, wildcard char: * or %").WithDestination(&schema).Build(),
				urafvcli.NewStringFlagBuilder("table").WithUsage("table name pattern, wildcard char: * or %").WithDestination(&table).Build(),
				urafvcli.NewStringFlagBuilder("database").WithAliases("d").WithUsage("database (mssql)").WithDestination(&db).Build(),
				// iSchemaFileFlag(&ifile),
				// ofileFlag(&ofile, "output file"),
				// schemaFile(&schema),
				// tableFlag(&table),
				// dbFlag(&db),
			},
			Action: func(c *cli.Context) error {
				data, err := yml.ReadSelectedYml(ifile, schema, table, "")
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
			Name:    "add_column",
			Usage:   "add column DDL",
			Aliases: []string{"ac"},
			Flags: []cli.Flag{
				urafvcli.NewStringFlagBuilder("input").WithAliases("i").WithValue("schema.yml").WithDestination(&ifile).Build(),
				urafvcli.NewStringFlagBuilder("output").WithAliases("o").WithUsage("output file (.png)").WithDestination(&ofile).Build(),
				urafvcli.NewStringFlagBuilder("schema").WithUsage("schema name pattern, wildcard char: * or %").WithDestination(&schema).Build(),
				urafvcli.NewStringFlagBuilder("table").WithUsage("table name pattern, wildcard char: * or %").WithDestination(&table).Build(),
				urafvcli.NewStringFlagBuilder("database").WithAliases("d").WithUsage("database (mssql)").WithDestination(&db).Build(),
				urafvcli.NewStringFlagBuilder("column").WithAliases("c").WithUsage("column").WithDestination(&col).Build(),
				// iSchemaFileFlag(&ifile),
				// ofileFlag(&ofile, "output file"),
				// schemaFile(&schema),
				// tableFlag(&table),
				// dbFlag(&db),
				// colFlag(&col),
			},
			Action: func(c *cli.Context) error {
				data, err := yml.ReadSelectedYml(ifile, schema, table, col)
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
			Name:    "drop_column",
			Usage:   "drop column DDL",
			Aliases: []string{"dc"},
			Flags: []cli.Flag{
				urafvcli.NewStringFlagBuilder("input").WithAliases("i").WithValue("schema.yml").WithDestination(&ifile).Build(),
				urafvcli.NewStringFlagBuilder("output").WithAliases("o").WithUsage("output file (.png)").WithDestination(&ofile).Build(),
				urafvcli.NewStringFlagBuilder("schema").WithUsage("schema name pattern, wildcard char: * or %").WithDestination(&schema).Build(),
				urafvcli.NewStringFlagBuilder("table").WithUsage("table name pattern, wildcard char: * or %").WithDestination(&table).Build(),
				urafvcli.NewStringFlagBuilder("database").WithAliases("d").WithUsage("database (mssql)").WithDestination(&db).Build(),
				urafvcli.NewStringFlagBuilder("column").WithAliases("c").WithUsage("column").WithDestination(&col).Build(),
				// iSchemaFileFlag(&ifile),
				// ofileFlag(&ofile, "output file"),
				// schemaFile(&schema),
				// tableFlag(&table),
				// dbFlag(&db),
				// colFlag(&col),
			},
			Action: func(c *cli.Context) error {
				data, err := yml.ReadSelectedYml(ifile, schema, table, col)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.DropColumn(data, db, ofile)
			},
		}
	}())

	sqlCmd.Subcommands = append(sqlCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:    "rename_column",
			Usage:   "rename column DDL",
			Aliases: []string{"rc"},
			Flags: []cli.Flag{
				urafvcli.NewStringFlagBuilder("input").WithAliases("i").WithValue("schema.yml").WithDestination(&ifile).Build(),
				urafvcli.NewStringFlagBuilder("output").WithAliases("o").WithUsage("output file (.png)").WithDestination(&ofile).Build(),
				urafvcli.NewStringFlagBuilder("schema").WithUsage("schema name pattern, wildcard char: * or %").WithDestination(&schema).Build(),
				urafvcli.NewStringFlagBuilder("table").WithUsage("table name pattern, wildcard char: * or %").WithDestination(&table).Build(),
				urafvcli.NewStringFlagBuilder("database").WithAliases("d").WithUsage("database (mssql)").WithDestination(&db).Build(),
				urafvcli.NewStringFlagBuilder("column").WithAliases("c").WithUsage("column").WithDestination(&col).Build(),
				// iSchemaFileFlag(&ifile),
				// ofileFlag(&ofile, "output file"),
				// schemaFile(&schema),
				// tableFlag(&table),
				// dbFlag(&db),
				// colFlag(&col),
			},
			Action: func(c *cli.Context) error {
				data, err := yml.ReadSelectedYml(ifile, schema, table, col)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.RenameColumn(data, db, ofile)
			},
		}
	}())

	sqlCmd.Subcommands = append(sqlCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:    "modify_column",
			Usage:   "modify column type DDL",
			Aliases: []string{"mc"},
			Flags: []cli.Flag{
				urafvcli.NewStringFlagBuilder("input").WithAliases("i").WithValue("schema.yml").WithDestination(&ifile).Build(),
				urafvcli.NewStringFlagBuilder("output").WithAliases("o").WithUsage("output file (.png)").WithDestination(&ofile).Build(),
				urafvcli.NewStringFlagBuilder("schema").WithUsage("schema name pattern, wildcard char: * or %").WithDestination(&schema).Build(),
				urafvcli.NewStringFlagBuilder("table").WithUsage("table name pattern, wildcard char: * or %").WithDestination(&table).Build(),
				urafvcli.NewStringFlagBuilder("database").WithAliases("d").WithUsage("database (mssql)").WithDestination(&db).Build(),
				urafvcli.NewStringFlagBuilder("column").WithAliases("c").WithUsage("column").WithDestination(&col).Build(),
				// iSchemaFileFlag(&ifile),
				// ofileFlag(&ofile, "output file"),
				// schemaFile(&schema),
				// tableFlag(&table),
				// dbFlag(&db),
				// colFlag(&col),
			},
			Action: func(c *cli.Context) error {
				data, err := yml.ReadSelectedYml(ifile, schema, table, col)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.ModifyColumn(data, db, ofile)
			},
		}
	}())

	sqlCmd.Subcommands = append(sqlCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:    "create_index",
			Usage:   "create index DDL",
			Aliases: []string{"ci"},
			Flags: []cli.Flag{
				urafvcli.NewStringFlagBuilder("input").WithAliases("i").WithValue("schema.yml").WithDestination(&ifile).Build(),
				urafvcli.NewStringFlagBuilder("output").WithAliases("o").WithUsage("output file (.png)").WithDestination(&ofile).Build(),
				urafvcli.NewStringFlagBuilder("schema").WithUsage("schema name pattern, wildcard char: * or %").WithDestination(&schema).Build(),
				urafvcli.NewStringFlagBuilder("table").WithUsage("table name pattern, wildcard char: * or %").WithDestination(&table).Build(),
				urafvcli.NewStringFlagBuilder("database").WithAliases("d").WithUsage("database (mssql)").WithDestination(&db).Build(),
				urafvcli.NewStringFlagBuilder("column").WithAliases("c").WithUsage("column").WithDestination(&col).Build(),
				// iSchemaFileFlag(&ifile),
				// ofileFlag(&ofile, "output file"),
				// schemaFile(&schema),
				// tableFlag(&table),
				// dbFlag(&db),
				// colFlag(&col),
			},
			Action: func(c *cli.Context) error {
				data, err := yml.ReadSelectedYml(ifile, schema, table, col)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.CreateIndex(data, db, ofile)
			},
		}
	}())

	sqlCmd.Subcommands = append(sqlCmd.Subcommands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:    "drop_index",
			Usage:   "drop index DDL",
			Aliases: []string{"di"},
			Flags: []cli.Flag{
				inSchemaFlagBuilder.WithDestination(&ifile).Build(),
				outFileFlagBuilder.WithDestination(&ofile).Build(),
				schemaFlagBuilder.WithDestination(&schema).Build(),
				tableFlagBuilder.WithDestination(&table).Build(),
				databaseFlagBuilder.Required(true).WithDestination(&db).Build(),
				columnFlagBuilder.WithDestination(&col).Build(),
				// iSchemaFileFlag(&ifile),
				// ofileFlag(&ofile, "output file"),
				// schemaFile(&schema),
				// tableFlag(&table),
				// dbFlag(&db),
				// colFlag(&col),
			},
			Action: func(c *cli.Context) error {
				data, err := yml.ReadSelectedYml(ifile, schema, table, col)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.DropIndex(data, db, ofile)
			},
		}
	}())

	return sqlCmd
}
