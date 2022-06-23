package xfmr

import (
	"github.com/rotisserie/eris"
)

type Xfmr struct{}

func NewXMFR() *Xfmr {
	return &Xfmr{}
}

func (s *Xfmr) YamlToExcel(infile, outfile string) error {
	data, err := s.loadYaml(infile)
	if err != nil {
		return eris.Wrapf(err, "failed to load the file %s", infile)
	}
	if err := s.saveExcel(data, outfile); err != nil {
		return eris.Wrapf(err, "failed to save the file %s", outfile)
	}
	return nil
}

func (s *Xfmr) YamlToText(infile, outfile, outtpl string) error {
	data, err := s.loadYaml(infile)
	if err != nil {
		return eris.Wrapf(err, "failed to load input file %s", infile)
	}
	if err := s.mergeTemplate(data, outfile, outtpl); err != nil {
		return eris.Wrap(err, "failed to generate the output from template")
	}
	return nil
}

func (s *Xfmr) YamlToPlantUML(infile, outfile, inschema string) error {
	data, err := s.loadYaml(infile)
	if err != nil {
		return eris.Wrapf(err, "failed to load input file %s", infile)
	}
	if err := s.savePlantUml(data, inschema, outfile); err != nil {
		return eris.Wrapf(err, "failed to save the file %s", outfile)
	}
	return nil
}

func (s *Xfmr) ExcelToYaml(infile string, outfile string) error {
	inData, err := s.loadExcel(infile)
	if err != nil {
		return eris.Wrapf(err, "failed to load the data from %s", infile)
	}
	if err := s.saveYaml(inData, outfile); err != nil {
		return eris.Wrapf(err, "failed to save %s", outfile)
	}
	return nil
}
