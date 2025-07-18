package model

import "gopkg.in/yaml.v3"

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
	ComputeType string `yaml:"computetype,omitempty"`
}

// MarshalYAML is a method that marshals a Column struct into a YAML Node
func (c *Column) MarshalYAML() (any, error) {
	// toYamlNodes is a function that converts a struct into a slice of YAML Nodes
	content := getContent(*c, func(name string, value any) bool {
		// check whether below columns should be excluded
		if name == "Identity" && value == "N" {
			return false
		}
		if name == "Unique" && value == "N" {
			return false
		}
		if name == "NotNull" && value == "N" {
			return false
		}
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
