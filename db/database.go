package db

import (
	"fmt"
	"reflect"

	"github.com/rotisserie/eris"
)

type Database struct{}

func NewDatabase() *Database {
	return &Database{}
}

func (s *Database) YamlToExcel(infile, outfile string) error {
	data, err := s.loadYaml(infile)
	if err != nil {
		return eris.Wrapf(err, "failed to load the file %s", infile)
	}
	if err := s.saveExcel(data, outfile); err != nil {
		return eris.Wrapf(err, "failed to save the file %s", outfile)
	}
	return nil
}

func (s *Database) ExcelToYaml(infile string, outfile string) error {
	inData, err := s.loadExcel(infile)
	if err != nil {
		return eris.Wrapf(err, "failed to load the data from %s", infile)
	}
	if err := s.saveYaml(inData, outfile); err != nil {
		return eris.Wrapf(err, "failed to save %s", outfile)
	}
	return nil
}

func (s *Database) convertToColumn(icols [][]interface{}) []InColumn {
	var trueFalse = func(value string) string {
		if value == "true" {
			return "Y"
		}
		return "N"
	}
	var createColumn = func(value []interface{}) InColumn {
		return InColumn{
			Name:        fmt.Sprintf("%v", value[0]),
			DataType:    fmt.Sprintf("%v", value[1]),
			IsPK:        trueFalse(fmt.Sprintf("%v", value[2])),
			IsUnique:    trueFalse(fmt.Sprintf("%v", value[3])),
			Nullable:    trueFalse(fmt.Sprintf("%v", value[4])),
			ForeignHint: fmt.Sprintf("%v", value[5]),
			Desc:        fmt.Sprintf("%v", value[6]),
		}
	}

	var columns []InColumn
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
