package model

type Index struct {
	Name    string   `yaml:"name,omitempty"`
	Unique  string   `yaml:"unique,omitempty"`
	Columns []string `yaml:"columns,flow,omitempty"`
}
