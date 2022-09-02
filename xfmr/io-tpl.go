package xfmr

import (
	"embed"
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/CloudyKit/jet/v6"
	"github.com/rotisserie/eris"
)

//go:embed templates
var tpl embed.FS

func (s *Xfmr) SaveToText(outfile, outtpl string) error {

	var loader jet.Loader
	if _, err := os.Stat(outtpl); errors.Is(err, os.ErrNotExist) {
		// read the template from embed file store
		tplFile := path.Join("templates", outtpl)
		content, err := tpl.ReadFile(tplFile)
		if err != nil {
			return eris.Wrapf(err, "failed to read the template %s", outtpl)
		}
		// use memory loader
		inMemloader := jet.NewInMemLoader()
		inMemloader.Set(path.Join("/", outtpl), string(content))
		loader = inMemloader
	} else {
		loader = jet.NewOSFileSystemLoader(filepath.Dir(outtpl))
	}

	views := jet.NewSet(loader, jet.InDevelopmentMode())
	view, err := views.GetTemplate(filepath.Base(outtpl))
	if err != nil {
		return eris.Wrapf(err, "failed to get the template %s", outtpl)
	}

	// output
	var fh *os.File
	if outfile == "stdout" {
		fh = os.Stdout
	} else {
		fh, err = os.Create(outfile)
		if err != nil {
			return eris.Wrapf(err, "failed to create the file %s", outfile)
		}
		defer fh.Close()
	}

	// merge
	if err := view.Execute(fh, nil, *s.Data); err != nil {
		return eris.Wrapf(err, "failed to merge the template %s", outtpl)
	}
	return nil
}
