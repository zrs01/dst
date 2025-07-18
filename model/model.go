package model

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// getContent converts a source object into a slice of YAML nodes.
func getContent(source any, isColumnVisible func(string, any) bool, finalTouch func(*yaml.Node, *yaml.Node)) []*yaml.Node {
	outNodes := make([]*yaml.Node, 0)

	v := reflect.ValueOf(source)
	typeOfColumn := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := typeOfColumn.Field(i)
		fieldName := fieldType.Name
		fieldValue := field.Interface()

		tags := strings.Split(fieldType.Tag.Get("yaml"), ",")
		tags = lo.Map(tags, func(s string, _ int) string {
			return strings.TrimSpace(s)
		})

		columnVisible := true
		if lo.Contains(tags, "omitempty") {
			columnVisible = fieldValue != nil && fieldValue != "" && fieldValue != "NULL"
		}
		if columnVisible && isColumnVisible != nil {
			columnVisible = isColumnVisible(fieldName, fieldValue)
		}
		if columnVisible {
			// replace the name of field if found
			if len(tags) >= 1 && tags[0] != "" {
				fieldName = tags[0]
			} else {
				fieldName = string(unicode.ToLower(rune(fieldName[0]))) + fieldName[1:]
			}

			nameNode := &yaml.Node{}
			nameNode.SetString(fieldName)

			valueNode := &yaml.Node{}
			valueNode.SetString(fieldValue.(string))

			if finalTouch != nil {
				finalTouch(nameNode, valueNode)
			}

			outNodes = append(outNodes, nameNode, valueNode)
		}
	}
	return outNodes
}
