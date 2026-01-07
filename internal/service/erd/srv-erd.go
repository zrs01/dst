package erd

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/CloudyKit/jet/v6"
	"github.com/CloudyKit/jet/v6/loaders/embedfs"
	"github.com/codeskyblue/go-sh"
	"github.com/sirupsen/logrus"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/service"
	"github.com/zrs01/dst/internal/service/txt"
	"github.com/ztrue/tracerr"
)

//go:embed templates/default.jet
var templateFS embed.FS

func Generate() error {
	logrus.Info("Generating ER Diagram...")
	builder := service.NewLoadBuilder(
		service.WithTablePattern(config.MergeSetting.TableName))
	schema, err := builder.LoadFromFile(config.MergeSetting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// filter the data if table name provided
	if err := schema.Filter(config.MergeSetting.TableName, config.MergeSetting.ColumnName); err != nil {
		return tracerr.Wrap(err)
	}

	var loader jet.Loader
	// Use default template if not specified
	if config.MergeSetting.Template == "" {
		logrus.Info("Template not provided, using default template")
		config.MergeSetting.Template = "templates/default.jet"
		loader = embedfs.NewLoader(filepath.Dir(config.MergeSetting.Template), templateFS)
	} else {
		logrus.Infof("Using template: %s", config.MergeSetting.Template)
		loader = jet.NewOSFileSystemLoader(filepath.Dir(config.MergeSetting.Template))
	}

	ext := filepath.Ext(config.MergeSetting.Output)
	switch ext {
	case ".puml":
		if err := txt.WriteWithLoader(loader, schema, config.MergeSetting.Template, config.MergeSetting.Output); err != nil {
			return tracerr.Wrap(err)
		}
		logrus.Infof("Generated file: %s", config.MergeSetting.Output)
	case ".png":
		if config.MergeSetting.Plantuml == "" {
			return tracerr.New("plantuml library path is not specified")
		}
		if _, err := os.Stat(config.MergeSetting.Plantuml); os.IsNotExist(err) {
			return tracerr.Wrap(err)
		}
		if err := txt.WriteWithLoader(loader, schema, config.MergeSetting.Template, "output.puml"); err != nil {
			return tracerr.Wrap(err)
		}
		defer os.Remove("output.puml")
		if err := sh.Command("java", "-jar", config.MergeSetting.Plantuml, "-o", filepath.Dir(config.MergeSetting.Output), "output.puml").Run(); err != nil {
			return tracerr.Wrap(err)
		}
		logrus.Infof("Generated file: %s", config.MergeSetting.Output)
	default:
		return tracerr.New(fmt.Sprintf("Output file extension '%s' is not supported", config.MergeSetting.Output))
	}
	return nil
}
