package yml

import (
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
	yamlIn "gopkg.in/yaml.v3"
)

func readFromFile(file string) (*model.DataDef, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var d model.DataDef
	if err := yamlIn.Unmarshal(yamlFile, &d); err != nil {
		return nil, tracerr.Wrap(err)
	}

	/* ------------------------- update reference tables ------------------------ */
	// the map table to speed up the lookup process
	var tableMap = make(map[string]*model.Table)
	for i := 0; i < len(d.Schemas); i++ {
		schema := &d.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			table := &schema.Tables[j]
			tableMap[table.Name] = table
		}
	}

	// update the reference table
	for i := 0; i < len(d.Schemas); i++ {
		schema := &d.Schemas[i]
		for j := 0; j < len(schema.Tables); j++ {
			table := &schema.Tables[j]
			for k := 0; k < len(table.Columns); k++ {
				column := &table.Columns[k]
				if column.ForeignKey != "" {
					fkTableName, fkColumnName, found := strings.Cut(column.ForeignKey, ".")
					if found {
						fkTable, ok := tableMap[fkTableName]
						if ok {
							fkTable.References = append(fkTable.References, model.Reference{
								ColumnName: fkColumnName,
								Foreign:    []model.ForeignTable{{Table: table.Name, Column: column.Name}},
							})
						} else {
							fmt.Printf("failed to find table '%s'", fkTableName)
						}
					}
				}
			}
		}
	}

	expandFixColumns(&d)

	// validate	data
	validateResult := model.Verify(&d)
	if len(validateResult) > 0 {
		lo.ForEach(validateResult, func(v string, _ int) {
			fmt.Println(v)
		})
		return nil, tracerr.Errorf("invalid data")
	}

	return &d, err
}
