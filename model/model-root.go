package model

type Root struct {
	Fixed   []*Column `yaml:"fixed,omitempty"`
	Schemas []*Schema `yaml:"schemas,omitempty"`
}
