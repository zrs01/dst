package export

import (
	"slices"

	"github.com/samber/lo"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/ddloader"
	"github.com/zrs01/dst/internal/ddwriter"
	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

func Generate() error {
	dataDef, err := ddloader.NewLoadBuilder(
		ddloader.WithTablePattern(config.Setting.Export.TableFilter)).
		LoadFromDB(config.Setting.Export.Dsn, config.Setting.Export.CommonColumnFile)
	if err != nil {
		return tracerr.Wrap(err)
	}
	if err := removeCommonColumns(dataDef); err != nil {
		return tracerr.Wrap(err)
	}
	return ddwriter.OutputYml(dataDef, config.Setting.Export.Output)
}

func removeCommonColumns(dataDef *model.DataDef) error {
	if config.Setting.Export.CommonColumnFile == "" {
		return nil
	}

	// read common column table
	cct, err := util.UnmarshalYml(config.Setting.Export.CommonColumnFile, &model.Table{})
	if err != nil {
		return tracerr.Wrap(err)
	}
	if len(cct.Columns) == 0 {
		return nil
	}
	// remove common columns in each table in DataDef
	for i := 0; i < len(dataDef.Schemas); i++ {
		schemas := dataDef.Schemas[i]
		for j := 0; j < len(schemas.Tables); j++ {
			table := schemas.Tables[j]
			for k := len(table.Columns) - 1; k >= 0; k-- {
				column := table.Columns[k]
				if lo.ContainsBy(cct.Columns, func(c *model.Column) bool {
					return c.Name == column.Name
				}) {
					table.Columns = slices.Delete(table.Columns, k, k+1)
				}
			}
		}
	}
	return nil
}
