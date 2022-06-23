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

func (s *Xfmr) mergeTemplate(data *InDB, outfile, outtpl string) error {

	// output
	f, err := os.Create(outfile)
	if err != nil {
		return eris.Wrapf(err, "failed to create the file %s", outfile)
	}
	defer f.Close()

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
	// merge
	if err := view.Execute(f, nil, *data); err != nil {
		return eris.Wrapf(err, "failed to merge the template %s", outtpl)
	}
	return nil
}
