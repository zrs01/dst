package dst

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/internal/ddwriter"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service"
	"github.com/ztrue/tracerr"
)

func RegisterExcelCmd() *cli.Command {
	var ifile, ofile, schema, table string
	var simple bool

	return &cli.Command{
		Name:  "excel",
		Usage: "transform from yaml to excel",
		Flags: []cli.Flag{
			flagbuilder.SchemaFileFlag().WithDestination(&ifile).Build(),
			flagbuilder.OutputFileFlag().WithUsage("output file (.xlsx)").WithDestination(&ofile).Build(),
			flagbuilder.SchemaNameFlag().WithDestination(&schema).Build(),
			flagbuilder.TableNameFlag().WithDestination(&table).Build(),
			flagbuilder.NewBoolFlag("simple").WithUsage("simple content").WithDestination(&simple).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			oext := lo.Ternary(ofile != "", strings.ToLower(filepath.Ext(ofile)), "")
			// data, err := ddloader.LoadWithFilter(ifile, schema, table, "")
			data, err := service.NewLoadBuilder(
				service.WithSchemaPattern(schema),
				service.WithTablePattern(table)).
				LoadFromFile(ifile)
			if err != nil {
				return tracerr.Wrap(err)
			}
			switch oext {
			case ".xlsx":
				if err := ddwriter.WriteXlsx(data, ofile, simple); err != nil {
					return tracerr.Wrap(err)
				}
				return nil
			}
			return tracerr.New("Not implemented yet")
		},
	}
}
