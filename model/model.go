package model

type DataDef struct {
	Fixed    []Column    `yaml:"fixed,omitempty"`
	OutFixed []OutColumn `yaml:"out_fixed,omitempty"`
	Schemas  []Schema    `yaml:"schemas,omitempty"`
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
	Columns       []Column       `yaml:"columns,omitempty"`
	OutColumns    []OutColumn    `yaml:"out_columns,omitempty"`
	References    []Reference    `yaml:"references,flow,omitempty"`
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
	ColumnName string         `yaml:"column,omitempty"`
	Foreign    []ForeignTable `yaml:"foreign,omitempty"`
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
