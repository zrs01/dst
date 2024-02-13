package transform

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"

	// "gopkg.in/yaml.v3"
	"github.com/goccy/go-yaml"
)

func ReadYml(file string) (*model.DataDef, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var d model.DataDef
	if err := yaml.Unmarshal(yamlFile, &d); err != nil {
		return nil, tracerr.Wrap(err)
	}

	/* ------------------------- update reference tables ------------------------ */
	findTable := func(table string) *model.Table {
		for i := 0; i < len(d.Schemas); i++ {
			for j := 0; j < len(d.Schemas[i].Tables); j++ {
				if d.Schemas[i].Tables[j].Name == table {
					return &d.Schemas[i].Tables[j]
				}
			}
		}
		return nil
	}

	for i := 0; i < len(d.Schemas); i++ {
		for j := 0; j < len(d.Schemas[i].Tables); j++ {
			for k := 0; k < len(d.Schemas[i].Tables[j].Columns); k++ {
				if d.Schemas[i].Tables[j].Columns[k].ForeignKey != "" {
					fkTableName, fkColumnName, found := strings.Cut(d.Schemas[i].Tables[j].Columns[k].ForeignKey, ".")
					if found {
						fkTable := findTable(fkTableName)
						if fkTable != nil {
							fkTable.References = append(d.Schemas[i].Tables[j].References, model.Reference{
								ColumnName: fkColumnName,
								Source: []model.ReferenceTable{
									{
										TableName:  d.Schemas[i].Tables[j].Name,
										ColumnName: d.Schemas[i].Tables[j].Columns[k].Name,
									},
								},
							})
						}
					}
				}
			}
		}
	}
	return &d, err
}

func FilterYml(ifile, schema, table string, column string) (*model.DataDef, error) {
	rawData, err := ReadYml(ifile)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	// copy fixed columns to each tables
	for i := 0; i < len(rawData.Schemas); i++ {
		schema := rawData.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			schema.Tables[j].Columns = append(schema.Tables[j].Columns, rawData.Fixed...)
		}
	}
	// clean the fixed column
	rawData.Fixed = []model.Column{}

	validateResult := model.Verify(rawData)
	if len(validateResult) > 0 {
		lo.ForEach(validateResult, func(v string, _ int) {
			fmt.Println(v)
		})
		return nil, tracerr.Errorf("invalid data")
	}
	data, err := model.FilterData(rawData, schema, table, column)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return data, nil
}

func WriteYml(data *model.DataDef, outfile string) error {
	// modify the columns in fixed to flow style
	pFixed := &data.Fixed
	for i := 0; i < len(*pFixed); i++ {
		(*data).OutFixed = make([]model.OutColumn, len(*pFixed))
		for j, column := range *pFixed {
			(*data).OutFixed[j].Value = column
		}
	}
	data.Fixed = nil

	// modify the columns to flow style
	pSchemas := &data.Schemas
	for i := 0; i < len(*pSchemas); i++ {
		var pTables = &(*pSchemas)[i].Tables
		for j := 0; j < len(*pTables); j++ {
			(*pTables)[j].OutColumns = make([]model.OutColumn, len((*pTables)[j].Columns))
			for k, column := range (*pTables)[j].Columns {
				(*pTables)[j].OutColumns[k].Value = column
				(*pTables)[j].Columns = nil
			}
		}
	}

	bytes, err := yaml.Marshal(data)
	if err != nil {
		return tracerr.Wrap(err)
	}

	output := string(bytes)
	// correct the names
	output = strings.ReplaceAll(output, "out_fixed", "fixed")
	output = strings.ReplaceAll(output, "out_columns", "columns")
	output = strings.ReplaceAll(output, "_column_values: ", "")
	output = strings.ReplaceAll(output, "\"N\"", "N")
	output = strings.ReplaceAll(output, "\"n\"", "n")
	output = strings.ReplaceAll(output, "\"Y\"", "Y")
	output = strings.ReplaceAll(output, "\"y\"", "y")

	if outfile == "" || outfile == "stdout" {
		fmt.Println(output)
	} else {
		os.WriteFile(outfile, []byte(output), fs.FileMode(0744))
	}
	return nil
}
