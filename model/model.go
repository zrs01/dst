package model

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

type DataDef struct {
	Fixed    []*Column   `yaml:"fixed,omitempty"`
	OutFixed []OutColumn `yaml:"out_fixed,omitempty"`
	Schemas  []*Schema   `yaml:"schemas,omitempty"`
}

type Schema struct {
	Name   string  `yaml:"name,omitempty" default:"Schema"`
	Desc   string  `yaml:"desc,omitempty"`
	Tables []Table `yaml:"tables,omitempty"`
}

type Table struct {
	Name          string         `yaml:"name,omitempty"`
	Title         string         `yaml:"title,omitempty"`
	Desc          string         `yaml:"desc,omitempty"`
	Version       bool           `yaml:"version,omitempty"`
	Columns       []*Column      `yaml:"columns,omitempty"`
	OutColumns    []OutColumn    `yaml:"out_columns,omitempty"`
	References    []Reference    `yaml:"references,omitempty"`
	OutReferences []OutReference `yaml:"out_references,omitempty"`
}

type Column struct {
	Name        string `yaml:"na,flow,omitempty"`
	DataType    string `yaml:"ty,omitempty"`
	Identity    string `yaml:"id,omitempty"`
	NotNull     string `yaml:"nu,omitempty" default:"N"`
	Unique      string `yaml:"un,omitempty"`
	Value       string `yaml:"va,omitempty"`
	ForeignKey  string `yaml:"fk,omitempty"`
	Cardinality string `yaml:"cd,omitempty"`
	Title       string `yaml:"tt,omitempty"`
	Index       string `yaml:"in,omitempty"`
	Desc        string `yaml:"dc,omitempty"`
	Compute     string `yaml:"cm,omitempty"`
}

type Reference struct {
	ColumnName string          `yaml:"column,omitempty"`
	Foreign    []*ForeignTable `yaml:"foreign,omitempty"`
}

type ForeignTable struct {
	Table  string `yaml:"table,omitempty"`
	Column string `yaml:"column,omitempty"`
}

type OutColumn struct {
	Value Column `yaml:"_column,flow,omitempty"`
}

type OutReference struct {
	Value Reference `yaml:"_reference,flow,omitempty"`
}

// MarshalYAML is a method that marshals a Column struct into a YAML Node
func (c *Column) MarshalYAML() (any, error) {
	// toYamlNodes is a function that converts a struct into a slice of YAML Nodes
	content := toYamlNodes(*c, func(name string, value any) bool {
		// exclude below columns if their value is N
		if name == "Identity" && value == "N" {
			return false
		}
		if name == "Unique" && value == "N" {
			return false
		}
		if name == "NotNull" && value == "N" {
			return false
		}
		// if name == "AutoIncrement" && value == "N" {
		// 	return false
		// }
		return true
	}, func(nameNode, valueNode *yaml.Node) {
		// make expression of compute column more readable
		if nameNode.Value == "compute" {
			valueNode.Style = yaml.LiteralStyle
		}
	})

	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Style:   yaml.FlowStyle,
		Content: content,
	}
	return node, nil
}

func (f *ForeignTable) MarshalYAML() (any, error) {
	content := toYamlNodes(*f, func(name string, value any) bool {
		return true
	}, nil)

	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Style:   yaml.FlowStyle,
		Content: content,
	}
	return node, nil
}

func toYamlNodes(source any, isValid func(string, any) bool, touch func(*yaml.Node, *yaml.Node)) []*yaml.Node {
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

		valid := true
		if lo.Contains(tags, "omitempty") {
			valid = fieldValue != nil && fieldValue != "" && fieldValue != "NULL"
		}
		if valid && isValid != nil {
			valid = isValid(fieldName, fieldValue)
		}
		if valid {
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

			if touch != nil {
				touch(nameNode, valueNode)
			}

			outNodes = append(outNodes, nameNode, valueNode)
		}
	}
	return outNodes
}
