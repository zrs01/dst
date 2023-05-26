package xfmr

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/rotisserie/eris"
	"github.com/shomali11/util/xstrings"
)

// //go:embed templates
// var tpl embed.FS

// func (s *Xfmr) SaveToText(outfile, outtpl string) error {
func (s *Xfmr) SaveToText(args TextArgs) error {

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
	loader := jet.NewOSFileSystemLoader(filepath.Dir(args.TemplateFile))

	views := jet.NewSet(loader, jet.InDevelopmentMode())
	view, err := views.GetTemplate(filepath.Base(args.TemplateFile))
	if err != nil {
		return eris.Wrapf(err, "failed to get the template %s", args.TemplateFile)
	}

	// output
	var fh *os.File
	if args.OutFile == "stdout" {
		fh = os.Stdout
	} else {
		fh, err = os.Create(args.OutFile)
		if err != nil {
			return eris.Wrapf(err, "failed to create the file %s", args.OutFile)
		}
		defer fh.Close()
	}

	// extract selected tables
	if xstrings.IsNotBlank(args.Pattern) {
		var isValidTable = func(name string) bool {
			parts := strings.Split(args.Pattern, ",")
			// fmt.Printf("*** %v\n", parts)
			for i := 0; i < len(parts); i++ {
				if wildCardMatch(strings.ToLower(strings.TrimSpace(parts[i])), strings.ToLower(name)) {
					return true
				}
			}
			return false
		}

		schemas := []Schema{}
		for _, schema := range s.Data.Schemas {
			tables := []Table{}
			// filter the desired tables
			for _, table := range schema.Tables {
				if isValidTable(table.Name) {
					tables = append(tables, table)
				}
			}
			if len(tables) > 0 {
				schema.Tables = tables
				schemas = append(schemas, schema)
			}
		}
		s.Data.Schemas = schemas
	}

	// merge
	if err := view.Execute(fh, nil, *s.Data); err != nil {
		return eris.Wrapf(err, "failed to merge the template %s", args.OutFile)
	}
	return nil
}
