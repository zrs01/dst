package dst

import (
	"context"

	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service"
	"github.com/zrs01/dst/internal/service/sql"
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

	var ifile, ofile, schema, table, db string
	setOptions := func() {
		config.Setting.DbType = lo.If(db != "", db).Else(config.Setting.DbType)
		// config.Setting.SchemaFilter = lo.If(schema != "", schema).Else(config.Setting.SchemaFilter)
		config.Setting.TableFilter = lo.If(table != "", table).Else(config.Setting.TableFilter)
		config.Setting.Output = lo.If(ofile != "", ofile).Else(config.Setting.Output)
	}

	cmd.Commands = append(cmd.Commands, func() *cli.Command {
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
				setOptions()
				// data, err := ddloader.LoadWithFilter(ifile, schema, table, "")
				data, err := service.NewLoadBuilder(
					// service.WithSchemaPattern(schema),
					service.WithTablePattern(table)).
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
				data, err := service.NewLoadBuilder(
					// service.WithSchemaPattern(schema),
					service.WithTablePattern(table)).
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
				data, err := service.NewLoadBuilder(
					// service.WithSchemaPattern(schema),
					service.WithTablePattern(table),
					service.WithColumnPattern(col)).
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
				data, err := service.NewLoadBuilder(
					// service.WithSchemaPattern(schema),
					service.WithTablePattern(table),
					service.WithColumnPattern(col)).
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
				data, err := service.NewLoadBuilder(
					// service.WithSchemaPattern(schema),
					service.WithTablePattern(table),
					service.WithColumnPattern(col)).
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
				data, err := service.NewLoadBuilder(
					// service.WithSchemaPattern(schema),
					service.WithTablePattern(table),
					service.WithColumnPattern(col)).
					LoadFromFile(ifile)
				if err != nil {
					return tracerr.Wrap(err)
				}
				return sql.ModifyColumn(data, db, ofile)
			},
		}
	}())

	cmd.Commands = append(cmd.Commands, func() *cli.Command {
		var input, output, schema, table, db, col string
		return &cli.Command{
			Name:  "ci",
			Usage: "create index DDL",
			Flags: []cli.Flag{
				flagbuilder.SchemaFileFlag().WithDestination(&input).Build(),
				flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&output).Build(),
				flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
				flagbuilder.TableNameFlag().WithDestination(&table).Build(),
				databaseFlagBuilder().WithDestination(&db).Build(),
				columnFlagBuilder().WithDestination(&col).Build(),
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				config.Setting.Input = lo.If(input != "", input).Else(config.Setting.Input)
				config.Setting.Output = lo.If(output != "", output).Else(config.Setting.Output)
				// config.Setting.SchemaFilter = lo.If(schema != "", schema).Else(config.Setting.SchemaFilter)
				config.Setting.TableFilter = lo.If(table != "", table).Else(config.Setting.TableFilter)
				config.Setting.ColumnFilter = lo.If(col != "", col).Else(config.Setting.ColumnFilter)

				// factory.NewService()

				// data, err := ddloader.LoadWithFilter(ifile, schema, table, col)
				data, err := service.NewLoadBuilder(
					// service.WithSchemaPattern(schema),
					service.WithTablePattern(table),
					service.WithColumnPattern(col)).
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
				data, err := service.NewLoadBuilder(
					// service.WithSchemaPattern(schema),
					service.WithTablePattern(table),
					service.WithColumnPattern(col)).
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
