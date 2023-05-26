package xfmr

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
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
		return xstrings.IsBlank(args.Schema) || (xstrings.IsNotBlank(args.Schema) && strings.EqualFold(args.Schema, name))
	}
	var isValidTable = func(name string) bool {
		if xstrings.IsBlank(args.Pattern) {
			return true
		}
		parts := strings.Split(args.Pattern, ",")
		for i := 0; i < len(parts); i++ {
			if wildCardMatch(strings.ToLower(strings.TrimSpace(parts[i])), strings.ToLower(name)) {
				return true
			}
		}
		return false
	}
	var addTable = func(builder *strings.Builder, tbNames []string) bool {
		for _, schema := range s.Data.Schemas {
			for _, table := range schema.Tables {
				if funk.Contains(tbNames, strings.ToLower(table.Name)) {
					builder.WriteString(fmt.Sprintf("\nentity %s", table.Name))
					if xstrings.IsNotBlank(table.Title) {
						builder.WriteString(fmt.Sprintf(" as \"%s\\n<size:11>(%s)</size>\"", table.Name, table.Title))
					}
					builder.WriteString(" {")

					columnContent := ""
					for _, column := range table.Columns {
						// simple mode, only show PK and FK
						if args.Simple {
							if strings.ToUpper(column.Identity) != "Y" && xstrings.IsBlank(column.ForeignKey) {
								continue
							}
						}
						columnContent += "\n  | "
						if strings.ToUpper(column.Identity) == "Y" {
							columnContent += "<size:11>PK</size>"
						}
						if xstrings.IsNotBlank(column.ForeignKey) {
							columnContent += "<size:11>FK</size>"
						}
						columnContent += fmt.Sprintf(" | <size:11>%s</size> | <size:11>%s</size> |", column.Name, column.DataType)
					}
					if xstrings.IsNotBlank(columnContent) {
						// write the table heading
						builder.WriteString("\n  |= |= <size:11>name</size> |= <size:11>type</size> |")
						builder.WriteString(columnContent)
					}

					builder.WriteString("\n}")
				}
			}
		}
		return false
	}
	var getCardinality = func(lcd, rcd, color string) string {
		lshape, rshape := "", ""
		lshape = xconditions.IfThenElse(strings.Contains(lcd, "*"), "}", "|").(string)
		lshape += xconditions.IfThenElse(strings.Contains(lcd, "0"), "o", "|").(string)
		rshape = xconditions.IfThenElse(strings.Contains(rcd, "0"), "o", "|").(string)
		rshape += xconditions.IfThenElse(strings.Contains(rcd, "*"), "{", "|").(string)

		var sb strings.Builder
		// if lcd != "" {
		// 	sb.WriteString("\"" + lcd + "\" ")
		// }
		sb.WriteString(fmt.Sprintf("%s-[%s]-%s", lshape, color, rshape))
		// sb.WriteString("--")
		// if rcd != "" {
		// 	sb.WriteString(" \"" + rcd + "\"")
		// }
		return sb.String()
	}
	var splitToTwo = func(s, sep string) (string, string) {
		v1, v2 := s, ""
		parts := strings.Split(s, sep)
		if len(parts) > 1 {
			v1 = parts[0]
			v2 = parts[1]
		}
		return v1, v2
	}

	/* ---------------------------------- main ---------------------------------- */

	var head strings.Builder
	var content strings.Builder

	// save all tables involved (include foreign tables)
	tables := []string{}
	cards := []string{}

	content.WriteString(fmt.Sprintf("@startuml %s", filepath.Base(args.OutFile)))
	// content.WriteString("\n\nskinparam linetype ortho")
	content.WriteString(`

skinparam {
	Linetype ortho
	ArrowFontSize 10
}`)

	for _, schema := range s.Data.Schemas {
		// user user selected schema or all
		if isValidSchema(schema.Name) {
			for _, table := range schema.Tables {
				// user selected prefix name or all
				if isValidTable(table.Name) {
					tables = append(tables, strings.ToLower(table.Name))

					for _, column := range table.Columns {
						if xstrings.IsNotBlank(column.ForeignKey) {
							fkTable, _ := splitToTwo(column.ForeignKey, ".")
							// foreign key in '<table>.<field>' format
							tbCd, fkCd := splitToTwo(column.Cardinality, ":")
							// if blank, default use '1 to many' relationship
							if xstrings.IsBlank(column.Cardinality) {
								tbCd = "1"
								fkCd = "*"
							}
							card := fmt.Sprintf("%s %s %s", fkTable, getCardinality(fkCd, tbCd, "#000000"), table.Name)
							if args.IncludeFK {
								card += fmt.Sprintf(" : %s", column.Name)
							}
							cards = append(cards, card)
							tables = append(tables, strings.ToLower(fkTable))
						}
					}
				}
			}
		}
	}
	tables = funk.UniqString(tables) // remove duplicated tables
	addTable(&head, tables)

	if !args.IncludeFK {
		cards = funk.UniqString(cards)
	}
	sort.Strings(cards)

	content.WriteString("\n" + head.String())
	content.WriteString("\n")
	for _, card := range cards {
		content.WriteString("\n" + card)
	}
	content.WriteString("\n\n@enduml")

	return content.String(), nil
}
