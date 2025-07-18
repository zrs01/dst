package model

import "gopkg.in/yaml.v3"

type ForeignTable struct {
	Table  string `yaml:"table,omitempty"`
	Column string `yaml:"column,omitempty"`
}

func (f *ForeignTable) MarshalYAML() (any, error) {
	content := getContent(*f, nil, nil)

	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Style:   yaml.FlowStyle,
		Content: content,
	}
	return node, nil
}
