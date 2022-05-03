package db

import (
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
