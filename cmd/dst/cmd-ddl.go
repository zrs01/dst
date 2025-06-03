package dst

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/internal/ddloader"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/sql"
	"github.com/ztrue/tracerr"
)

/* -------------------------------------------------------------------------- */
/*                                     SQL                                    */
/* -------------------------------------------------------------------------- */
func RegisterDLLCmd() *cli.Command {
	cmd := &cli.Command{
		Name:  "ddl",
		Usage: "Generate SQL DDL",
	}

	databaseFlagBuilder := func() *flagbuilder.StringFlag {
		return flagbuilder.NewStringFlag("database").WithAliases("d").Required(true).WithUsage("database (mariadb, mssql)")
	}
	columnFlagBuilder := func() *flagbuilder.StringFlag {
		return flagbuilder.NewStringFlag("column").WithAliases("c").WithUsage("column")
	}

	cmd.Commands = append(cmd.Commands, func() *cli.Command {
		var ifile, ofile, schema, table, db string
		return &cli.Command{
			Name:  "ct",
			Usage: "create table",
			Flags: []cli.Flag{
				flagbuilder.SchemaFileFlag().WithDestination(&ifile).Build(),
				flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
				flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
				flagbuilder.TableNameFlag().WithDestination(&table).Build(),
				databaseFlagBuilder().WithDestination(&db).Build(),
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				// data, err := ddloader.LoadWithFilter(ifile, schema, table, "")
				data, err := ddloader.NewLoadBuilder(
					ddloader.WithSchemaPattern(schema),
					ddloader.WithTablePattern(table)).
					LoadFromFile(ifile)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.CreateTable(data, db, ofile)
			},
		}
	}())

	cmd.Commands = append(cmd.Commands, func() *cli.Command {
		var ifile, ofile, schema, table, db string
		return &cli.Command{
			Name:  "dt",
			Usage: "drop table",
			Flags: []cli.Flag{
				flagbuilder.SchemaFileFlag().WithDestination(&ifile).Build(),
				flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
				flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
				flagbuilder.TableNameFlag().WithDestination(&table).Build(),
				databaseFlagBuilder().WithDestination(&db).Build(),
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				// data, err := ddloader.LoadWithFilter(ifile, schema, table, "")
				data, err := ddloader.NewLoadBuilder(
					ddloader.WithSchemaPattern(schema),
					ddloader.WithTablePattern(table)).
					LoadFromFile(ifile)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.DropTable(data, db, ofile)
			},
		}
	}())

	cmd.Commands = append(cmd.Commands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:  "ac",
			Usage: "add column",
			Flags: []cli.Flag{
				flagbuilder.SchemaFileFlag().WithDestination(&ifile).Build(),
				flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
				flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
				flagbuilder.TableNameFlag().WithDestination(&table).Build(),
				databaseFlagBuilder().WithDestination(&db).Build(),
				columnFlagBuilder().WithDestination(&col).Build(),
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				// data, err := ddloader.LoadWithFilter(ifile, schema, table, col)
				data, err := ddloader.NewLoadBuilder(
					ddloader.WithSchemaPattern(schema),
					ddloader.WithTablePattern(table),
					ddloader.WithColumnPattern(col)).
					LoadFromFile(ifile)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.AddColumn(data, db, ofile)
			},
		}
	}())

	cmd.Commands = append(cmd.Commands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:  "dc",
			Usage: "drop column",
			Flags: []cli.Flag{
				flagbuilder.SchemaFileFlag().WithDestination(&ifile).Build(),
				flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
				flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
				flagbuilder.TableNameFlag().WithDestination(&table).Build(),
				databaseFlagBuilder().WithDestination(&db).Build(),
				columnFlagBuilder().WithDestination(&col).Build(),
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				// data, err := ddloader.LoadWithFilter(ifile, schema, table, col)
				data, err := ddloader.NewLoadBuilder(
					ddloader.WithSchemaPattern(schema),
					ddloader.WithTablePattern(table),
					ddloader.WithColumnPattern(col)).
					LoadFromFile(ifile)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.DropColumn(data, db, ofile)
			},
		}
	}())

	cmd.Commands = append(cmd.Commands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:  "rc",
			Usage: "rename column",
			Flags: []cli.Flag{
				flagbuilder.SchemaFileFlag().WithDestination(&ifile).Build(),
				flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
				flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
				flagbuilder.TableNameFlag().WithDestination(&table).Build(),
				databaseFlagBuilder().WithDestination(&db).Build(),
				columnFlagBuilder().WithDestination(&col).Build(),
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				data, err := ddloader.NewLoadBuilder(
					ddloader.WithSchemaPattern(schema),
					ddloader.WithTablePattern(table),
					ddloader.WithColumnPattern(col)).
					LoadFromFile(ifile)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.RenameColumn(data, db, ofile)
			},
		}
	}())

	cmd.Commands = append(cmd.Commands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:  "mc",
			Usage: "modify column type",
			Flags: []cli.Flag{
				flagbuilder.SchemaFileFlag().WithDestination(&ifile).Build(),
				flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
				flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
				flagbuilder.TableNameFlag().WithDestination(&table).Build(),
				databaseFlagBuilder().WithDestination(&db).Build(),
				columnFlagBuilder().WithDestination(&col).Build(),
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				// data, err := ddloader.LoadWithFilter(ifile, schema, table, col)
				data, err := ddloader.NewLoadBuilder(
					ddloader.WithSchemaPattern(schema),
					ddloader.WithTablePattern(table),
					ddloader.WithColumnPattern(col)).
					LoadFromFile(ifile)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.ModifyColumn(data, db, ofile)
			},
		}
	}())

	cmd.Commands = append(cmd.Commands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:  "ci",
			Usage: "create index DDL",
			Flags: []cli.Flag{
				flagbuilder.SchemaFileFlag().WithDestination(&ifile).Build(),
				flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
				flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
				flagbuilder.TableNameFlag().WithDestination(&table).Build(),
				databaseFlagBuilder().WithDestination(&db).Build(),
				columnFlagBuilder().WithDestination(&col).Build(),
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				// data, err := ddloader.LoadWithFilter(ifile, schema, table, col)
				data, err := ddloader.NewLoadBuilder(
					ddloader.WithSchemaPattern(schema),
					ddloader.WithTablePattern(table),
					ddloader.WithColumnPattern(col)).
					LoadFromFile(ifile)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.CreateIndex(data, db, ofile)
			},
		}
	}())

	cmd.Commands = append(cmd.Commands, func() *cli.Command {
		var ifile, ofile, schema, table, db, col string
		return &cli.Command{
			Name:  "di",
			Usage: "drop index",
			Flags: []cli.Flag{
				flagbuilder.SchemaFileFlag().WithDestination(&ifile).Build(),
				flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
				flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
				flagbuilder.TableNameFlag().WithDestination(&table).Build(),
				databaseFlagBuilder().WithDestination(&db).Build(),
				columnFlagBuilder().WithDestination(&col).Build(),
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				// data, err := ddloader.LoadWithFilter(ifile, schema, table, col)
				data, err := ddloader.NewLoadBuilder(
					ddloader.WithSchemaPattern(schema),
					ddloader.WithTablePattern(table),
					ddloader.WithColumnPattern(col)).
					LoadFromFile(ifile)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.DropIndex(data, db, ofile)
			},
		}
	}())

	return cmd
}
