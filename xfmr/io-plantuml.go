package xfmr

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/shomali11/util/xconditions"
	"github.com/shomali11/util/xstrings"
	"github.com/thoas/go-funk"
)

func (s *Xfmr) SaveToPlantUML(args DiagramArgs) error {
	uml, err := s.buildPlantUml(args)
	if err != nil {
		return eris.Wrapf(err, "failed to build UML")
	}
	ioutil.WriteFile(args.OutFile, []byte(uml), 0744)
	return nil
}

func (s *Xfmr) buildPlantUml(args DiagramArgs) (string, error) {

	var isValidSchema = func(name string) bool {
		// return (schemaName != "" && strings.ToLower(schemaName) == strings.ToLower(name)) || schemaName == ""
		return xstrings.IsBlank(args.Schema) || (xstrings.IsNotBlank(args.Schema) && strings.EqualFold(strings.ToLower(args.Schema), strings.ToLower(name)))
	}
	var isValidTable = func(name string) bool {
		if xstrings.IsBlank(args.TablePrefix) {
			return true
		}
		parts := strings.Split(args.TablePrefix, ",")
		for i := 0; i < len(parts); i++ {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(strings.TrimSpace(parts[i]))) {
				return true
			}
		}
		return false
		// return (xstrings.IsNotBlank(args.TablePrefix) && strings.HasPrefix(strings.ToLower(name), strings.ToLower(args.TablePrefix)))
	}
	var addTable = func(builder *strings.Builder, tbNames []string) bool {
		for _, schema := range s.Data.Schemas {
			for _, table := range schema.Tables {
				if funk.Contains(tbNames, strings.ToLower(table.Name)) {
					builder.WriteString(fmt.Sprintf("\nentity %s {", table.Name))
					builder.WriteString("\n  |= |= <size:11>name</size> |= <size:11>type</size> |")
					for _, column := range table.Columns {
						builder.WriteString("\n  | ")
						if strings.ToUpper(column.Identity) == "Y" {
							builder.WriteString("<size:11>PK</size>")
						}
						if xstrings.IsNotBlank(column.ForeignKey) {
							builder.WriteString("<size:11>FK</size>")
						}
						builder.WriteString(fmt.Sprintf(" | <size:11>%s</size> | <size:11>%s</size> |", column.Name, column.DataType))
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
		// if lcd != "" {
		// 	sb.WriteString("\"" + lcd + "\" ")
		// }
		sb.WriteString(lshape + "--" + rshape)
		// sb.WriteString("--")
		// if rcd != "" {
		// 	sb.WriteString(" \"" + rcd + "\"")
		// }
		return sb.String()
	}

	/* ---------------------------------- main ---------------------------------- */

	var head strings.Builder
	var body strings.Builder
	var content strings.Builder

	tables := []string{}

	content.WriteString(fmt.Sprintf("@startuml %s", filepath.Base(args.OutFile)))
	content.WriteString("\n\nskinparam linetype ortho")

	for _, schema := range s.Data.Schemas {
		if isValidSchema(schema.Name) {
			for _, table := range schema.Tables {
				if isValidTable(table.Name) {
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
	}
	tables = funk.UniqString(tables) // remove duplicated values
	addTable(&head, tables)

	content.WriteString("\n" + head.String())
	content.WriteString("\n" + body.String())

	content.WriteString("\n\n@enduml")
	return content.String(), nil
}
