package db

import (
	"io/ioutil"
	"strings"

	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v2"
)

type InDB struct {
	Fixed   []InColumn `yaml:"fixed,omitempty"`
	Schemas []InSchema `yaml:"schemas,omitempty"`
}

type InSchema struct {
	Name   string    `yaml:"name,omitempty" default:"Schema"`
	Desc   string    `yaml:"description,omitempty"`
	Tables []InTable `yaml:"tables,omitempty"`
}

type InTable struct {
	Name    string     `yaml:"name,omitempty"`
	Desc    string     `yaml:"desc,omitempty"`
	Columns []InColumn `yaml:"columns,omitempty"`
}

type InColumn struct {
	Name           string `yaml:"na,omitempty"`
	DataType       string `yaml:"ty,omitempty"`
	Identity       string `yaml:"id,omitempty"`
	NotNull        string `yaml:"nu,omitempty" default:"N"`
	Value          string `yaml:"va,omitempty"`
	ForeignKeyHint string `yaml:"fk,omitempty"`
	Desc           string `yaml:"dc,omitempty"`
}

type OutDB struct {
	Schemas []OutSchema `yaml:"schemas"`
}

type OutSchema struct {
	Name   string     `yaml:"name,omitempty" default:"Schema"`
	Desc   string     `yaml:"description,omitempty"`
	Tables []OutTable `yaml:"tables,omitempty"`
}

type OutTable struct {
	Name    string      `yaml:"name,omitempty"`
	Desc    string      `yaml:"description,omitempty"`
	Columns []OutColumn `yaml:"columns,omitempty"`
}

type OutColumn struct {
	Values []string `yaml:"_column_values,flow,omitempty"`
}

func (s *Database) loadYaml(infile string) (*InDB, error) {
	yamlFile, err := ioutil.ReadFile(infile)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to read the file %s", infile)
	}
	var d InDB
	if err := yaml.Unmarshal(yamlFile, &d); err != nil {
		return nil, eris.Wrapf(err, "failed to unmarshal the file %s", infile)
	}
	return &d, nil
}

func (s *Database) saveYaml(data *OutDB, outfile string) error {
	yaml.FutureLineWrap()
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return eris.Wrap(err, "failed to transform to yaml")
	}
	output := string(bytes)
	// remove useless key
	output = strings.ReplaceAll(output, "_column_values: ", "")
	ioutil.WriteFile(outfile, []byte(output), 0744)
	return nil
}
