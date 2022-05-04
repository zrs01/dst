package xfmr

import (
	"html/template"
	"os"

	"github.com/rotisserie/eris"
	// "github.com/CloudyKit/jet/v6"
)

func (s *Xfmr) MergeTemplate(infile, outfile, outtpl string) error {
	in, err := s.loadYaml(infile)
	if err != nil {
		return eris.Wrapf(err, "failed to load input file %s", infile)
	}

	// template
	t, err := template.ParseFiles(outtpl)
	if err != nil {
		return eris.Wrapf(err, "failed to parse the template file %s", outtpl)
	}
	// output
	f, err := os.Create(outfile)
	if err != nil {
		return eris.Wrapf(err, "failed to create the file %s", outfile)
	}
	defer f.Close()
	// merge
	if err := t.Execute(f, *in); err != nil {
		return eris.Wrapf(err, "failed to merge the template %s and data %+v", outtpl, *in)
	}
	return nil
}
