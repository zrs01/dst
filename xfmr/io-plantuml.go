package xfmr

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/shomali11/util/xconditions"
	"github.com/thoas/go-funk"
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

	var isValidSchema = func(name string) bool {
		return (schemaName != "" && strings.ToLower(schemaName) == strings.ToLower(name)) || schemaName == ""
	}
	var addTable = func(builder *strings.Builder, tbNames []string) bool {
		for _, schema := range data.Schemas {
			for _, table := range schema.Tables {
				if funk.Contains(tbNames, strings.ToLower(table.Name)) {
					builder.WriteString(fmt.Sprintf("\nentity %s {", table.Name))
					for _, column := range table.Columns {
						cn := column.Name
						if strings.ToUpper(column.Identity) == "Y" {
							cn = "<u>" + cn + "</u>"
						}
						builder.WriteString(fmt.Sprintf("\n  <size:11>%s %s</size>", cn, column.DataType))
					}
					builder.WriteString("\n  ..")
					builder.WriteString("\n}")
				}
			}
		}
		return false
	}
	var getCardinality = func(cd string, reverse bool) string {
		lcd, lshape, rcd, rshape := cd, "", "", ""
		cdParts := strings.Split(cd, ":")
		if len(cdParts) > 1 {
			if reverse {
				lcd, rcd = cdParts[1], cdParts[0]
			} else {
				lcd, rcd = cdParts[0], cdParts[1]
			}
		}
		lshape = xconditions.IfThenElse(strings.Contains(lcd, "*"), "}", "|").(string)
		lshape += xconditions.IfThenElse(strings.Contains(lcd, "0"), "o", "|").(string)
		rshape = xconditions.IfThenElse(strings.Contains(rcd, "0"), "o", "|").(string)
		rshape += xconditions.IfThenElse(strings.Contains(rcd, "*"), "{", "|").(string)

		var sb strings.Builder
		if lcd != "" {
			sb.WriteString("\"" + lcd + "\" ")
		}
		// sb.WriteString(lshape + "--" + rshape)
		sb.WriteString("--")
		if rcd != "" {
			sb.WriteString(" \"" + rcd + "\"")
		}
		return sb.String()
	}

	/* ---------------------------------- main ---------------------------------- */

	var head strings.Builder
	var body strings.Builder
	var content strings.Builder

	tables := []string{}

	graphName := schemaName
	if graphName == "" {
		graphName = filepath.Base(outfile)
	}
	content.WriteString(fmt.Sprintf("@startuml %s", graphName))
	// content.WriteString("\n\nskinparam linetype ortho")

	for _, schema := range data.Schemas {
		for _, table := range schema.Tables {
			if isValidSchema(schema.Name) {
				tables = append(tables, strings.ToLower(table.Name))

				for _, column := range table.Columns {
					if len(column.ForeignKey) > 0 {
						fkParts := strings.Split(column.ForeignKey, ".")
						fkTable := column.ForeignKey
						if len(fkParts) > 1 {
							fkTable = fkParts[0]
						}
						body.WriteString(fmt.Sprintf("\n%s %s %s", fkTable, getCardinality(column.Cardinality, true), table.Name))
						// body.WriteString(fmt.Sprintf("\n%s %s %s", table.Name, getCardinality(column.Cardinality, false), fkTable))
						tables = append(tables, strings.ToLower(fkTable))
					}
				}
			}
		}
	}
	tables = funk.UniqString(tables) // remove duplicated values
	addTable(&head, tables)

	content.WriteString("\n" + head.String())
	content.WriteString("\n" + body.String())

	content.WriteString("\n\n@enduml")
	return content.String(), nil
}
