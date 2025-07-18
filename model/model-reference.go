package model

type Reference struct {
	ColumnName string          `yaml:"column,omitempty"`
	Foreign    []*ForeignTable `yaml:"foreign,omitempty"`
}
