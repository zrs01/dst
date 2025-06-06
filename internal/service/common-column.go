package service

import (
	"slices"

	"github.com/samber/lo"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

func RemoveCommonColumns(schema *model.Schema) error {
	if config.Setting.Ccf == "" {
		return nil
	}

	// read common column table
	cct, err := util.UnmarshalYml(config.Setting.Ccf, &model.Table{})
	if err != nil {
		return tracerr.Wrap(err)
	}
	if len(cct.Columns) == 0 {
		return nil
	}
	// remove common columns in each table
	for j := 0; j < len(schema.Tables); j++ {
		table := schema.Tables[j]
		for k := len(table.Columns) - 1; k >= 0; k-- {
			column := table.Columns[k]
			if lo.ContainsBy(cct.Columns, func(c *model.Column) bool {
				return c.Name == column.Name
			}) {
				table.Columns = slices.Delete(table.Columns, k, k+1)
			}
		}
	}
	return nil
}

func AppendCommonColumns(schema *model.Schema) error {
	if config.Setting.Ccf == "" {
		return nil
	}

	// read common column table
	cct, err := util.UnmarshalYml(config.Setting.Ccf, &model.Table{})
	if err != nil {
		return tracerr.Wrap(err)
	}
	if len(cct.Columns) == 0 {
		return nil
	}
	// append common columns in each table
	for j := 0; j < len(schema.Tables); j++ {
		table := schema.Tables[j]
		table.Columns = append(table.Columns, cct.Columns...)
	}
	return nil
}
