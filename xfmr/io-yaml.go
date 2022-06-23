package xfmr

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v2"
)

const (
	CName           = 0
	CDataType       = 1
	CIdentity       = 2
	CNotNull        = 3
	CUnique         = 4
	CValue          = 5
	CForeignKeyHint = 6
	CDesc           = 7
)

type OutColumn struct {
	Values Column `yaml:"_column_values,flow,omitempty"`
}

func (s *Xfmr) loadYaml(infile string) (*InDB, error) {
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

func (s *Xfmr) saveYaml(data *InDB, outfile string) error {
	yaml.FutureLineWrap()
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return eris.Wrap(err, "failed to transform to yaml")
	}
	output := string(bytes)
	// correct the names
	output = strings.ReplaceAll(output, "_column_values: ", "")
	output = strings.ReplaceAll(output, "out_columns", "columns")
	output = strings.ReplaceAll(output, "\"N\"", "N")
	output = strings.ReplaceAll(output, "\"n\"", "n")
	output = strings.ReplaceAll(output, "\"Y\"", "Y")
	output = strings.ReplaceAll(output, "\"y\"", "y")
	ioutil.WriteFile(outfile, []byte(output), 0744)
	return nil
}

func (s *Xfmr) VerifyForeignKey(infile string) error {
	data, err := s.loadYaml(infile)
	if err != nil {
		return eris.Wrapf(err, "failed to load the data from %s", infile)
	}

	tables := make(map[string][]Column)

	// convert to map for easy searching
	for _, schema := range data.Schemas {
		for _, table := range schema.Tables {
			tables[table.Name] = table.Columns
		}
	}
	tables["fixed"] = data.Fixed

	isFKExist := func(fk string) bool {
		// fk format: table.field
		tf := strings.Split(fk, ".")
		if columns, found := tables[tf[0]]; found {
			for _, column := range columns {
				if column.Name == tf[1] {
					return true
				}
			}
		}
		return false
	}

	// check the foreign key whether exists
	for k, columns := range tables {
		for _, column := range columns {
			if column.ForeignKey != "" {
				if !isFKExist(column.ForeignKey) {
					fmt.Printf("Warning: '%s', FK '%s' cannot be found\n", k, column.ForeignKey)
				}
			}
		}
	}
	return nil
}
