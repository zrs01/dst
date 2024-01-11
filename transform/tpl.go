package transform

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/ztrue/tracerr"
)

// //go:embed templates
// var tpl embed.FS

// func (s *Xfmr) SaveToText(outfile, outtpl string) error {
func WriteTpl(data *DataDef, tplf string, out string, pattern string) error {

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
	if out == "" {
		fh = os.Stdout
	} else {
		fh, err = os.Create(out)
		if err != nil {
			return tracerr.Wrap(err)
		}
		defer fh.Close()
	}

	// extract selected tables
	outData := &DataDef{
		Fixed:   data.Fixed,
		Schemas: []Schema{},
	}
	// function to check whether the pattern is valid
	var isEligibleTable = func(name string) bool {
		if pattern == "" {
			return true
		}
		parts := strings.Split(pattern, ",")
		// fmt.Printf("*** %v\n", parts)
		for i := 0; i < len(parts); i++ {
			if wildCardMatch(strings.ToLower(strings.TrimSpace(parts[i])), strings.ToLower(name)) {
				return true
			}
		}
		return false
	}

	for _, schema := range data.Schemas {
		tables := []Table{}
		// filter the desired tables
		for _, table := range schema.Tables {
			if isEligibleTable(table.Name) {
				tables = append(tables, table)
			}
		}
		if len(tables) > 0 {
			schema.Tables = tables
			outData.Schemas = append(outData.Schemas, schema)
		}
	}

	// merge
	if err := view.Execute(fh, nil, *outData); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}
