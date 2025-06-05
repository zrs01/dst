package text

import (
	"embed"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/CloudyKit/jet/v6/loaders/embedfs"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/ddwriter"
	"github.com/zrs01/dst/internal/service"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

func Generate() error {
	schema, err := service.NewLoadBuilder(
		service.WithSchemaPattern(config.Setting.SchemaFilter),
		service.WithTablePattern(config.Setting.TableFilter),
		service.WithColumnPattern(config.Setting.ColumnFilter)).
		LoadFromFile(config.Setting.Input) // the data does not filter by schema and table
	if err != nil {
		return tracerr.Wrap(err)
	}
	if config.Setting.Template == "" {
		return ddwriter.OutputYml(schema, config.Setting.Output)
	}
	return WriteWithFileLoader(schema, config.Setting.Template, config.Setting.Output)
}

func WriteWithFileLoader(schema *model.Schema, tplf string, out string) error {
	return WriteWithLoader(jet.NewOSFileSystemLoader(filepath.Dir(tplf)), schema, tplf, out)
}

func WriteWithEmbedFSLoader(fs embed.FS, schema *model.Schema, tplf string, out string) error {
	return WriteWithLoader(embedfs.NewLoader(filepath.Dir(tplf), fs), schema, tplf, out)
}

func WriteWithInMemoryLoader(schema *model.Schema, tplf, tplc string, out string) error {
	loader := jet.NewInMemLoader()
	loader.Set(filepath.Base(tplf), tplc)
	return WriteWithLoader(loader, schema, tplf, out)
}

func WriteWithLoader(loader jet.Loader, schema *model.Schema, tplf string, out string) error {
	views := jet.NewSet(loader)
	setJetFunc(views)
	view, err := views.GetTemplate(filepath.Base(tplf))
	if err != nil {
		return tracerr.Wrap(err)
	}

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

	// merge
	if err := view.Execute(fh, nil, *schema); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func setJetFunc(views *jet.Set) {
	jetToLowerCamel(views)
	jetToCamel(views)
	jetPlural(views)
	jetSingular(views)
	jetJavaType(views)
	jetTypescriptType(views)
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
	view.AddGlobalFunc("toPlural", func(args jet.Arguments) reflect.Value {
		pluralize := pluralize.NewClient()
		return reflect.ValueOf(pluralize.Plural(args.Get(0).Interface().(string)))
	})
}

func jetSingular(view *jet.Set) {
	view.AddGlobalFunc("toSingular", func(args jet.Arguments) reflect.Value {
		pluralize := pluralize.NewClient()
		return reflect.ValueOf(pluralize.Singular(args.Get(0).Interface().(string)))
	})
}

func jetJavaType(view *jet.Set) {
	view.AddGlobalFunc("toJavaType", func(args jet.Arguments) reflect.Value {
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

func jetTypescriptType(view *jet.Set) {
	view.AddGlobalFunc("toTypescriptType", func(args jet.Arguments) reflect.Value {
		srcType := strings.ToLower(args.Get(0).Interface().(string))
		switch {
		case strings.Contains(srcType, "bigint"):
			return reflect.ValueOf("number")
		case strings.Contains(srcType, "bit"):
			return reflect.ValueOf("boolean")
		case strings.Contains(srcType, "date"):
			return reflect.ValueOf("Date")
		case strings.Contains(srcType, "datetime"):
			return reflect.ValueOf("Date")
		case strings.Contains(srcType, "decimal"):
			return reflect.ValueOf("number")
		case strings.Contains(srcType, "float"):
			return reflect.ValueOf("number")
		case strings.Contains(srcType, "int"):
			return reflect.ValueOf("number")
		case strings.Contains(srcType, "longtext"):
			return reflect.ValueOf("string")
		case strings.Contains(srcType, "mediumint"):
			return reflect.ValueOf("number")
		case strings.Contains(srcType, "mediumtext"):
			return reflect.ValueOf("string")
		case strings.Contains(srcType, "smallint"):
			return reflect.ValueOf("string")
		case strings.Contains(srcType, "text"):
			return reflect.ValueOf("string")
		case strings.Contains(srcType, "time"):
			return reflect.ValueOf("date")
		case strings.Contains(srcType, "timestamp"):
			return reflect.ValueOf("date")
		case strings.Contains(srcType, "tinyint"):
			return reflect.ValueOf("byte")
		case strings.Contains(srcType, "tinytext"):
			return reflect.ValueOf("string")
		case strings.Contains(srcType, "varbinary"):
			return reflect.ValueOf("byte[]")
		case strings.Contains(srcType, "varchar"):
			return reflect.ValueOf("string")
		case strings.Contains(srcType, "char"):
			return reflect.ValueOf("string")
		default:
			return reflect.ValueOf("Unknown Type: " + srcType)
		}
	})
}
