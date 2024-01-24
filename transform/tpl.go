package transform

import (
	"os"
	"path/filepath"

	"github.com/CloudyKit/jet/v6"
	"github.com/ztrue/tracerr"
)

// WriteTpl generates a template using the provided data and template file,
// and writes the output to the specified file or standard output.
//
// Parameters:
//   - data: A pointer to a DataDef struct containing the data for the template.
//   - tplf: A string specifying the path to the template file.
//   - out: A string specifying the path to the output file. If empty, the output
//     will be written to standard output.
//   - pattern: A string specifying a pattern to filter the desired tables from
//     the data.
//
// Return:
// - An error if any occurred during the execution of the function.
func WriteTpl(data *DataDef, tplf string, out string) error {

	// var loader jet.Loader
	// if _, err := os.Stat(outtpl); errors.Is(err, os.ErrNotExist) {
	// 	// read the template from embed file store
	// 	tplFile := path.Join("templates", outtpl)
	// 	content, err := tpl.ReadFile(tplFile)
	// 	if err != nil {
	// 		return eris.Wrapf(err, "failed to read the template %s", outtpl)
	// 	}
	// 	// use memory loader
	// 	inMemloader := jet.NewInMemLoader()
	// 	inMemloader.Set(path.Join("/", outtpl), string(content))
	// 	loader = inMemloader
	// } else {
	// 	loader = jet.NewOSFileSystemLoader(filepath.Dir(outtpl))
	// }
	loader := jet.NewOSFileSystemLoader(filepath.Dir(tplf))

	// views := jet.NewSet(loader, jet.InDevelopmentMode())
	views := jet.NewSet(loader)
	view, err := views.GetTemplate(filepath.Base(tplf))
	if err != nil {
		return tracerr.Wrap(err)
	}

	// output
	var fh *os.File
	if out == "" || out == "stdout" {
		fh = os.Stdout
	} else {
		fh, err = os.Create(out)
		if err != nil {
			return tracerr.Wrap(err)
		}
		defer fh.Close()
	}

	// merge
	if err := view.Execute(fh, nil, *data); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}
