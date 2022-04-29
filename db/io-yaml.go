package db

import (
	"io/ioutil"
	"strings"

	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v2"
)

type InDB struct {
	Schemas []InSchema `yaml:"schemas"`
}

type InSchema struct {
	Name   string    `yaml:"name" default:"Schema"`
	Desc   string    `yaml:"description"`
	Tables []InTable `yaml:"tables"`
}

type InTable struct {
	Name    string          `yaml:"name"`
	Desc    string          `yaml:"description"`
	Columns [][]interface{} `yaml:"columns"`
}

type InColumn struct {
	Name        string
	DataType    string
	IsPK        string
	IsUnique    string
	Nullable    string
	ForeignHint string
	Desc        string
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
