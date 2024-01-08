package xfmr

import (
	"fmt"
	"strings"

	"io/fs"
	"os"

	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
)

type OutColumn struct {
	Values Column `yaml:"_column_values,flow,omitempty"`
}

func (s *Xfmr) LoadYaml(infile string) {
	yamlFile, err := os.ReadFile(infile)
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
	os.WriteFile(outfile, []byte(output), fs.FileMode(0744))
	return nil
}

func (s *Xfmr) Verify() error {
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

	for k, columns := range tables {
		for i, column := range columns {
			if funk.IsEmpty(column.Name) {
				logrus.Warnf("'%s', missing column name at line %d", k, i)
			}
			if funk.IsEmpty(column.DataType) {
				logrus.Warnf("'%s', missing data type of the column '%s' at line %d", k, column.Name, i)
			}
			if !funk.IsEmpty(column.ForeignKey) {
				// check the foreign key whether exists
				if !isFKExist(column.ForeignKey) {
					logrus.Warnf("'%s', FK '%s' cannot be found at line %d\n", k, column.ForeignKey, i)
				}
			}
		}
	}
	return nil
}
