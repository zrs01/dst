package transform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/codeskyblue/go-sh"
	"github.com/rotisserie/eris"
	"github.com/ztrue/tracerr"
)

func WriteERD(data *DataDef, tplf string, out string) error {
	// plantuml file name = output file name + .puml
	outPuml := strings.TrimSuffix(out, filepath.Ext(out)) + ".puml"
	if err := writePlantuml(data, tplf, outPuml); err != nil {
		return tracerr.Wrap(err)
	}
	if out != "" {
		lib, err := SearchPathFiles("plantuml*.jar")
		if err != nil {
			return tracerr.Wrap(err)
		}
		if len(lib) == 0 {
			return tracerr.New("plantuml*.jar not found in PATH environment variable")
		}
		fmt.Printf("> use plantuml library found in '%s'\n", lib[0])
		if err := sh.Command("java", "-jar", lib[0], outPuml).Run(); err != nil {
			return eris.Wrapf(err, "failed to generate the diagram")
		}
	}

	return nil
}

func writePlantuml(data *DataDef, tplf string, out string) error {
	loader := jet.NewOSFileSystemLoader(filepath.Dir(tplf))

	views := jet.NewSet(loader)
	view, err := views.GetTemplate(filepath.Base(tplf))
	if err != nil {
		return tracerr.Wrap(err)
	}

	// output
	var fh *os.File
	if out == "" {
		fh = os.Stdout
	} else {
		fh, err = os.Create(out)
		if err != nil {
			return tracerr.Wrap(err)
		}
		defer fh.Close()
	}

	if err := view.Execute(fh, nil, *data); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}
