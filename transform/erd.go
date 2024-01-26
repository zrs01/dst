package transform

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codeskyblue/go-sh"
	"github.com/rotisserie/eris"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

//go:embed templates/erd.jet
var tplERD []byte

func WriteERD(data *model.DataDef, tplf string, out string) error {
	if tplf == "" {
		// use default template
		fh, err := os.CreateTemp("", "dst")
		if err != nil {
			return tracerr.Wrap(err)
		}
		tplf = fh.Name()
		fh.Close()

		if err := os.WriteFile(tplf, tplERD, 0644); err != nil {
			return tracerr.Wrap(err)
		}
		defer os.Remove(tplf)
	}

	outPuml := strings.TrimSuffix(out, filepath.Ext(out)) + ".puml"
	if err := WriteFileTpl(data, tplf, outPuml); err != nil {
		return tracerr.Wrap(err)
	}
	// if err := writePlantuml(data, tplf, outPuml); err != nil {
	// 	return tracerr.Wrap(err)
	// }
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

// func writePlantuml(data *model.DataDef, tplf string, out string) error {
// 	loader := jet.NewOSFileSystemLoader(filepath.Dir(tplf))

// 	views := jet.NewSet(loader)
// 	view, err := views.GetTemplate(filepath.Base(tplf))
// 	if err != nil {
// 		return tracerr.Wrap(err)
// 	}

// 	// output
// 	var fh *os.File
// 	if out == "" {
// 		fh = os.Stdout
// 	} else {
// 		fh, err = os.Create(out)
// 		if err != nil {
// 			return tracerr.Wrap(err)
// 		}
// 		defer fh.Close()
// 	}

// 	if err := view.Execute(fh, nil, *data); err != nil {
// 		return tracerr.Wrap(err)
// 	}
// 	return nil
// }
