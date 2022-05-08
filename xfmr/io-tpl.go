package xfmr

import (
	"os"
	"path/filepath"

	"github.com/CloudyKit/jet/v6"
	"github.com/rotisserie/eris"
)

func (s *Xfmr) mergeTemplate(data *InDB, outfile, outtpl string) error {
	// output
	f, err := os.Create(outfile)
	if err != nil {
		return eris.Wrapf(err, "failed to create the file %s", outfile)
	}
	defer f.Close()

	// template
	views := jet.NewSet(jet.NewOSFileSystemLoader(filepath.Dir(outtpl)), jet.InDevelopmentMode())
	view, err := views.GetTemplate(filepath.Base(outtpl))
	if err != nil {
		return eris.Wrapf(err, "failed to get the template %s", outtpl)
	}
	// merge
	if err := view.Execute(f, nil, *data); err != nil {
		return eris.Wrapf(err, "failed to merge the template %s", outtpl)
	}
	return nil
}
