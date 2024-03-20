package tpl

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/CloudyKit/jet/v6"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
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
}

func WriteMemoryTpl(data *model.DataDef, tplf, tplc string, out string) error {
	return writeTpl(data, memoryLoader(data, tplf, tplc), tplf, out)
}

func writeTpl(data *model.DataDef, loader jet.Loader, tplf string, out string) error {
	views := jet.NewSet(loader)
	setJetFunc(views)
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

func setJetFunc(views *jet.Set) {
	jetToLowerCamel(views)
	jetToCamel(views)
	jetPlural(views)
	jetSingular(views)
	jetJavaType(views)
}

func jetToLowerCamel(views *jet.Set) {
	views.AddGlobalFunc("toLowerCamel", func(args jet.Arguments) reflect.Value {
		return reflect.ValueOf(strcase.ToLowerCamel(args.Get(0).Interface().(string)))
	})
}

func jetToCamel(views *jet.Set) {
	views.AddGlobalFunc("toCamel", func(args jet.Arguments) reflect.Value {
		return reflect.ValueOf(strcase.ToCamel(args.Get(0).Interface().(string)))
	})
}

func jetPlural(view *jet.Set) {
	view.AddGlobalFunc("plural", func(args jet.Arguments) reflect.Value {
		pluralize := pluralize.NewClient()
		return reflect.ValueOf(pluralize.Plural(args.Get(0).Interface().(string)))
	})
}

func jetSingular(view *jet.Set) {
	view.AddGlobalFunc("singular", func(args jet.Arguments) reflect.Value {
		pluralize := pluralize.NewClient()
		return reflect.ValueOf(pluralize.Singular(args.Get(0).Interface().(string)))
	})
}

func jetJavaType(view *jet.Set) {
	view.AddGlobalFunc("javaType", func(args jet.Arguments) reflect.Value {
		srcType := strings.ToLower(args.Get(0).Interface().(string))
		switch {
		case strings.Contains(srcType, "bigint"):
			return reflect.ValueOf("Long")
		case strings.Contains(srcType, "bit"):
			return reflect.ValueOf("Boolean")
		case strings.Contains(srcType, "date"):
			return reflect.ValueOf("LocalDate")
		case strings.Contains(srcType, "datetime"):
			return reflect.ValueOf("LocalDateTime")
		case strings.Contains(srcType, "decimal"):
			return reflect.ValueOf("BigDecimal")
		case strings.Contains(srcType, "float"):
			return reflect.ValueOf("Double")
		case strings.Contains(srcType, "int"):
			return reflect.ValueOf("Integer")
		case strings.Contains(srcType, "longtext"):
			return reflect.ValueOf("String")
		case strings.Contains(srcType, "mediumint"):
			return reflect.ValueOf("Integer")
		case strings.Contains(srcType, "mediumtext"):
			return reflect.ValueOf("String")
		case strings.Contains(srcType, "smallint"):
			return reflect.ValueOf("Short")
		case strings.Contains(srcType, "text"):
			return reflect.ValueOf("String")
		case strings.Contains(srcType, "time"):
			return reflect.ValueOf("LocalTime")
		case strings.Contains(srcType, "timestamp"):
			return reflect.ValueOf("Timestamp")
		case strings.Contains(srcType, "tinyint"):
			return reflect.ValueOf("Byte")
		case strings.Contains(srcType, "tinytext"):
			return reflect.ValueOf("String")
		case strings.Contains(srcType, "varbinary"):
			return reflect.ValueOf("byte[]")
		case strings.Contains(srcType, "varchar"):
			return reflect.ValueOf("String")
		case strings.Contains(srcType, "char"):
			return reflect.ValueOf("String")
		default:
			return reflect.ValueOf("Unknown Type: " + srcType)
		}
	})
}
