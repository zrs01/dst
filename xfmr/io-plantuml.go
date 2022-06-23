package xfmr

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rotisserie/eris"
)

func (s *Xfmr) savePlantUml(data *InDB, schemaName, outfile string) error {
	uml, err := s.buildPlantUml(data, schemaName, outfile)
	if err != nil {
		return eris.Wrapf(err, "failed to build UML")
	}
	ioutil.WriteFile(outfile, []byte(uml), 0744)
	return nil
}

func (s *Xfmr) buildPlantUml(data *InDB, schemaName, outfile string) (string, error) {
	var head strings.Builder
	var conn strings.Builder
	var body strings.Builder

	graphName := schemaName
	if graphName == "" {
		graphName = filepath.Base(outfile)
	}
	body.WriteString(fmt.Sprintf("@startuml %s\n\nskinparam linetype ortho", graphName))

	for _, schema := range data.Schemas {
		for _, table := range schema.Tables {
			if (schemaName != "" && strings.ToLower(schemaName) == strings.ToLower(schema.Name)) || schemaName == "" {
				head.WriteString(fmt.Sprintf("\nentity %s {", table.Name))

				for _, column := range table.Columns {
					head.WriteString(fmt.Sprintf("\n  {field} %s %s", column.Name, column.DataType))

					if len(column.ForeignKey) > 0 {
						fkParts := strings.Split(column.ForeignKey, ".")
						fkTable := column.ForeignKey
						if len(fkParts) > 1 {
							fkTable = fkParts[0]
						}
						conn.WriteString(fmt.Sprintf("\n%s ||--|{ %s", fkTable, table.Name))
					}
				}
				head.WriteString("\n}")
			}
		}
	}

	body.WriteString("\n\n" + head.String())
	body.WriteString("\n\n" + conn.String())

	body.WriteString("\n\n@enduml")
	return body.String(), nil
}
