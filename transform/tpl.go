package transform

import (
	"os"
	"path/filepath"

	"github.com/CloudyKit/jet/v6"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

// WriteFileTpl generates a template using the provided data and template file,
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
func WriteFileTpl(data *model.DataDef, tplf string, out string) error {
	return writeTpl(data, fileLoader(data, tplf), tplf, out)

	// views := jet.NewSet(loader, jet.InDevelopmentMode())
	// views := jet.NewSet(loader)
	// view, err := views.GetTemplate(filepath.Base(tplf))
	// if err != nil {
	// 	return tracerr.Wrap(err)
	// }

	// // output
	// var fh *os.File
	// if out == "" || out == "stdout" {
	// 	fh = os.Stdout
	// } else {
	// 	fh, err = os.Create(out)
	// 	if err != nil {
	// 		return tracerr.Wrap(err)
	// 	}
	// 	defer fh.Close()
	// }

	// // merge
	// if err := view.Execute(fh, nil, *data); err != nil {
	// 	return tracerr.Wrap(err)
	// }
	// return nil
}

func WriteMemoryTpl(data *model.DataDef, tplf, tplc string, out string) error {
	return writeTpl(data, memoryLoader(data, tplf, tplc), tplf, out)
}

func writeTpl(data *model.DataDef, loader jet.Loader, tplf string, out string) error {
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

func fileLoader(data *model.DataDef, tplf string) jet.Loader {
	return jet.NewOSFileSystemLoader(filepath.Dir(tplf))
}

func memoryLoader(data *model.DataDef, tplf string, tplc string) jet.Loader {
	loader := jet.NewInMemLoader()
	loader.Set(filepath.Base(tplf), tplc)
	return loader
}
