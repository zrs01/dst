package db

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v2"
)

type Database struct {
	Data Data
}

type Data struct {
	Schemas []Schema `yaml:"schemas"`
}

type Schema struct {
	Name   string  `yaml:"name" default:"Schema"`
	Desc   string  `yaml:"description"`
	Tables []Table `yaml:"tables"`
}

type Table struct {
	Name    string          `yaml:"name"`
	Desc    string          `yaml:"description"`
	Columns [][]interface{} `yaml:"columns"`
}

type Column struct {
	Name        string
	DataType    string
	IsPK        string
	IsUnique    string
	Nullable    string
	ForeignHint string
	Desc        string
}

func NewDatabase() *Database {
	return &Database{}
}

func (s *Database) Load(f string) error {
	yamlFile, err := ioutil.ReadFile(f)
	if err != nil {
		return eris.Wrapf(err, "failed to read the file %s", f)
	}
	var d Data
	if err := yaml.Unmarshal(yamlFile, &d); err != nil {
		return eris.Wrapf(err, "failed to unmarshal the file %s", f)
	}
	s.Data = d
	return nil
}

func (s *Database) convertToColumn(icols [][]interface{}) []Column {
	var trueFalse = func(value string) string {
		if value == "true" {
			return "Y"
		}
		return "N"
	}
	var createColumn = func(value []interface{}) Column {
		return Column{
			Name:        fmt.Sprintf("%v", value[0]),
			DataType:    fmt.Sprintf("%v", value[1]),
			IsPK:        trueFalse(fmt.Sprintf("%v", value[2])),
			IsUnique:    trueFalse(fmt.Sprintf("%v", value[3])),
			Nullable:    trueFalse(fmt.Sprintf("%v", value[4])),
			ForeignHint: fmt.Sprintf("%v", value[5]),
			Desc:        fmt.Sprintf("%v", value[6]),
		}
	}

	var columns []Column
	for _, icol := range icols {
		// due to merge array in yaml will cause increase array one level, so need check is array of array
		if reflect.TypeOf(icol[0]).Kind() == reflect.Slice {
			for _, sub := range icol {
				columns = append(columns, createColumn(sub.([]interface{})))
			}
		} else {
			columns = append(columns, createColumn(icol))
		}
	}
	return columns
}
