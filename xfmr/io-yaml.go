package xfmr

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
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

func (s *Xfmr) LoadYaml(infile string) {
	yamlFile, err := ioutil.ReadFile(infile)
	if err != nil {
		panic(fmt.Sprintf("failed to read the file %s", infile))
	}
	var d InDB
	if err := yaml.Unmarshal(yamlFile, &d); err != nil {
		panic(fmt.Sprintf("failed to unmarshal the file %s", infile))
	}
	s.Data = &d
}

func (s *Xfmr) SaveToYaml(outfile string) error {
	yaml.FutureLineWrap()
	bytes, err := yaml.Marshal(s.Data)
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

func (s *Xfmr) VerifyForeignKey() error {
	tables := make(map[string][]Column)

	// convert to map for easy searching
	for _, schema := range s.Data.Schemas {
		for _, table := range schema.Tables {
			tables[table.Name] = table.Columns
		}
	}
	tables["fixed"] = s.Data.Fixed

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
					logrus.Warnf("'%s', FK '%s' cannot be found\n", k, column.ForeignKey)
				}
			}
		}
	}
	return nil
}
