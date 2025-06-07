package dst

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service/ddl"
)

func RegisterDDLCmd() *cli.Command {
	cmd := &cli.Command{
		Name:  "ddl",
		Usage: "Generate SQL DDL",
	}

	// databaseFlagBuilder := func() *flagbuilder.StringFlag {
	// 	return flagbuilder.NewStringFlag("database").WithAliases("d").WithRequired(true).WithUsage("database (mariadb, mssql)")
	// }
	// columnFlagBuilder := func() *flagbuilder.StringFlag {
	// 	return flagbuilder.NewStringFlag("column").WithAliases("c").WithUsage("column")
	// }

	// var ifile, ofile string
	// setOptions := func() {
	// 	// config.Setting.DbType = lo.If(db != "", db).Else(config.Setting.DbType)
	// 	// config.Setting.SchemaFilter = lo.If(schema != "", schema).Else(config.Setting.SchemaFilter)
	// 	config.Setting.TableName = lo.If(table != "", table).Else(config.Setting.TableName)
	// 	config.Setting.Output = lo.If(ofile != "", ofile).Else(config.Setting.Output)
	// }

	// cmd.Commands = append(cmd.Commands, func() *cli.Command {
	// 	var options config.SettingDef
	// 	return &cli.Command{
	// 		Name:  "cd",
	// 		Usage: "create database",
	// 		Flags: []cli.Flag{
	// 			flagbuilder.InputFileFlag().WithDestination(&options.Sdf).Build(),
	// 		},
	// 		Action: func(ctx context.Context, cmd *cli.Command) error {
	// 			config.InitSetting(config.DataDefConf, options)
	// 			return ddl.GenerateCreateDatabase()
	// 		},
	// 	}
	// }())

	cmd.Commands = append(cmd.Commands, func() *cli.Command {
		var options config.SettingDef
		return &cli.Command{
			Name:  "ct",
			Usage: "create table",
			Flags: []cli.Flag{
				flagbuilder.SdfFlag().WithDestination(&options.Sdf).Build(),
				flagbuilder.TableNameFlag().WithDestination(&options.TableName).Build(),
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				config.InitSetting(config.DataDefConf, options)
				return ddl.GenerateCreateTable()
				// setOptions()
				// // data, err := ddloader.LoadWithFilter(ifile, schema, table, "")
				// data, err := service.NewLoadBuilder(
				// 	// service.WithSchemaPattern(schema),
				// 	service.WithTablePattern(table)).
				// 	LoadFromFile(ifile)
				// if err != nil {
				// 	return tracerr.Wrap(err)
				// }
				// return sql.CreateTable(data, db, ofile)
			},
		}
	}())

	// cmd.Commands = append(cmd.Commands, func() *cli.Command {
	// 	var options config.SettingDef
	// 	return &cli.Command{
	// 		Name:  "ci",
	// 		Usage: "create index",
	// 		Flags: []cli.Flag{
	// 			flagbuilder.InputFileFlag().WithDestination(&options.Sdf).Build(),
	// 			flagbuilder.TableNameFlag().WithDestination(&options.TableName).Build(),
	// 		},
	// 		Action: func(ctx context.Context, cmd *cli.Command) error {
	// 			config.InitSetting(config.DataDefConf, options)
	// 			return ddl.GenerateCreateIndex()
	// 		},
	// 	}
	// }())

	cmd.Commands = append(cmd.Commands, func() *cli.Command {
		var options config.SettingDef
		return &cli.Command{
			Name:  "dt",
			Usage: "drop table",
			Flags: []cli.Flag{
				flagbuilder.SdfFlag().WithDestination(&options.Sdf).Build(),
				flagbuilder.TableNameFlag().WithDestination(&options.TableName).Build(),
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				config.InitSetting(config.DataDefConf, options)
				return ddl.GenerateDropTable()
			},
		}
	}())

	// cmd.Commands = append(cmd.Commands, func() *cli.Command {
	// 	var ifile, ofile, schema, table, db, col string
	// 	return &cli.Command{
	// 		Name:  "ac",
	// 		Usage: "add column",
	// 		Flags: []cli.Flag{
	// 			flagbuilder.InputFileFlag().WithDestination(&ifile).Build(),
	// 			flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
	// 			flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
	// 			flagbuilder.TableNameFlag().WithDestination(&table).Build(),
	// 			databaseFlagBuilder().WithDestination(&db).Build(),
	// 			columnFlagBuilder().WithDestination(&col).Build(),
	// 		},
	// 		Action: func(ctx context.Context, cmd *cli.Command) error {
	// 			// data, err := ddloader.LoadWithFilter(ifile, schema, table, col)
	// 			data, err := service.NewLoadBuilder(
	// 				// service.WithSchemaPattern(schema),
	// 				service.WithTablePattern(table),
	// 				service.WithColumnPattern(col)).
	// 				LoadFromFile(ifile)
	// 			if err != nil {
	// 				return tracerr.Wrap(err)
	// 			}
	// 			return sql.AddColumn(data, db, ofile)
	// 		},
	// 	}
	// }())

	// cmd.Commands = append(cmd.Commands, func() *cli.Command {
	// 	var ifile, ofile, schema, table, db, col string
	// 	return &cli.Command{
	// 		Name:  "dc",
	// 		Usage: "drop column",
	// 		Flags: []cli.Flag{
	// 			flagbuilder.InputFileFlag().WithDestination(&ifile).Build(),
	// 			flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
	// 			flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
	// 			flagbuilder.TableNameFlag().WithDestination(&table).Build(),
	// 			databaseFlagBuilder().WithDestination(&db).Build(),
	// 			columnFlagBuilder().WithDestination(&col).Build(),
	// 		},
	// 		Action: func(ctx context.Context, cmd *cli.Command) error {
	// 			// data, err := ddloader.LoadWithFilter(ifile, schema, table, col)
	// 			data, err := service.NewLoadBuilder(
	// 				// service.WithSchemaPattern(schema),
	// 				service.WithTablePattern(table),
	// 				service.WithColumnPattern(col)).
	// 				LoadFromFile(ifile)
	// 			if err != nil {
	// 				return tracerr.Wrap(err)
	// 			}
	// 			return sql.DropColumn(data, db, ofile)
	// 		},
	// 	}
	// }())

	// cmd.Commands = append(cmd.Commands, func() *cli.Command {
	// 	var ifile, ofile, schema, table, db, col string
	// 	return &cli.Command{
	// 		Name:  "rc",
	// 		Usage: "rename column",
	// 		Flags: []cli.Flag{
	// 			flagbuilder.InputFileFlag().WithDestination(&ifile).Build(),
	// 			flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
	// 			flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
	// 			flagbuilder.TableNameFlag().WithDestination(&table).Build(),
	// 			databaseFlagBuilder().WithDestination(&db).Build(),
	// 			columnFlagBuilder().WithDestination(&col).Build(),
	// 		},
	// 		Action: func(ctx context.Context, cmd *cli.Command) error {
	// 			data, err := service.NewLoadBuilder(
	// 				// service.WithSchemaPattern(schema),
	// 				service.WithTablePattern(table),
	// 				service.WithColumnPattern(col)).
	// 				LoadFromFile(ifile)
	// 			if err != nil {
	// 				return tracerr.Wrap(err)
	// 			}
	// 			return sql.RenameColumn(data, db, ofile)
	// 		},
	// 	}
	// }())

	// cmd.Commands = append(cmd.Commands, func() *cli.Command {
	// 	var ifile, ofile, schema, table, db, col string
	// 	return &cli.Command{
	// 		Name:  "mc",
	// 		Usage: "modify column type",
	// 		Flags: []cli.Flag{
	// 			flagbuilder.InputFileFlag().WithDestination(&ifile).Build(),
	// 			flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
	// 			flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
	// 			flagbuilder.TableNameFlag().WithDestination(&table).Build(),
	// 			databaseFlagBuilder().WithDestination(&db).Build(),
	// 			columnFlagBuilder().WithDestination(&col).Build(),
	// 		},
	// 		Action: func(ctx context.Context, cmd *cli.Command) error {
	// 			// data, err := ddloader.LoadWithFilter(ifile, schema, table, col)
	// 			data, err := service.NewLoadBuilder(
	// 				// service.WithSchemaPattern(schema),
	// 				service.WithTablePattern(table),
	// 				service.WithColumnPattern(col)).
	// 				LoadFromFile(ifile)
	// 			if err != nil {
	// 				return tracerr.Wrap(err)
	// 			}
	// 			return sql.ModifyColumn(data, db, ofile)
	// 		},
	// 	}
	// }())

	// cmd.Commands = append(cmd.Commands, func() *cli.Command {
	// 	var input, output, schema, table, db, col string
	// 	return &cli.Command{
	// 		Name:  "ci",
	// 		Usage: "create index DDL",
	// 		Flags: []cli.Flag{
	// 			flagbuilder.InputFileFlag().WithDestination(&input).Build(),
	// 			flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&output).Build(),
	// 			flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
	// 			flagbuilder.TableNameFlag().WithDestination(&table).Build(),
	// 			databaseFlagBuilder().WithDestination(&db).Build(),
	// 			columnFlagBuilder().WithDestination(&col).Build(),
	// 		},
	// 		Action: func(ctx context.Context, cmd *cli.Command) error {
	// 			config.Setting.Sdf = lo.If(input != "", input).Else(config.Setting.Sdf)
	// 			config.Setting.Output = lo.If(output != "", output).Else(config.Setting.Output)
	// 			// config.Setting.SchemaFilter = lo.If(schema != "", schema).Else(config.Setting.SchemaFilter)
	// 			config.Setting.TableName = lo.If(table != "", table).Else(config.Setting.TableName)
	// 			config.Setting.ColumnName = lo.If(col != "", col).Else(config.Setting.ColumnName)

	// 			// factory.NewService()

	// 			// data, err := ddloader.LoadWithFilter(ifile, schema, table, col)
	// 			data, err := service.NewLoadBuilder(
	// 				// service.WithSchemaPattern(schema),
	// 				service.WithTablePattern(table),
	// 				service.WithColumnPattern(col)).
	// 				LoadFromFile(ifile)
	// 			if err != nil {
	// 				return tracerr.Wrap(err)
	// 			}
	// 			return sql.CreateIndex(data, db, ofile)
	// 		},
	// 	}
	// }())

	// cmd.Commands = append(cmd.Commands, func() *cli.Command {
	// 	var ifile, ofile, schema, table, db, col string
	// 	return &cli.Command{
	// 		Name:  "di",
	// 		Usage: "drop index",
	// 		Flags: []cli.Flag{
	// 			flagbuilder.InputFileFlag().WithDestination(&ifile).Build(),
	// 			flagbuilder.OutputFileFlag().WithUsage("output file (text file)").WithDestination(&ofile).Build(),
	// 			flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
	// 			flagbuilder.TableNameFlag().WithDestination(&table).Build(),
	// 			databaseFlagBuilder().WithDestination(&db).Build(),
	// 			columnFlagBuilder().WithDestination(&col).Build(),
	// 		},
	// 		Action: func(ctx context.Context, cmd *cli.Command) error {
	// 			// data, err := ddloader.LoadWithFilter(ifile, schema, table, col)
	// 			data, err := service.NewLoadBuilder(
	// 				// service.WithSchemaPattern(schema),
	// 				service.WithTablePattern(table),
	// 				service.WithColumnPattern(col)).
	// 				LoadFromFile(ifile)
	// 			if err != nil {
	// 				return tracerr.Wrap(err)
	// 			}
	// 			return sql.DropIndex(data, db, ofile)
	// 		},
	// 	}
	// }())

	return cmd
}
