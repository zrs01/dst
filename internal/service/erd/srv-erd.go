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
		service.WithTablePattern(config.Setting.TableName))
	schema, err := builder.LoadFromFile(config.Setting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// filter the data if table name provided
	if err := schema.Filter(config.Setting.TableName, config.Setting.ColumnName); err != nil {
		return tracerr.Wrap(err)
	}

	var loader jet.Loader
	// Use default template if not specified
	if config.Setting.Template == "" {
		logrus.Info("Template not provided, using default template")
		config.Setting.Template = "templates/default.jet"
		loader = embedfs.NewLoader(filepath.Dir(config.Setting.Template), templateFS)
	} else {
		logrus.Infof("Using template: %s", config.Setting.Template)
		loader = jet.NewOSFileSystemLoader(filepath.Dir(config.Setting.Template))
	}

	ext := filepath.Ext(config.Setting.Output)
	switch ext {
	case ".puml":
		if err := txt.WriteWithLoader(loader, schema, config.Setting.Template, config.Setting.Output); err != nil {
			return tracerr.Wrap(err)
		}
		logrus.Infof("Generated file: %s", config.Setting.Output)
	case ".png":
		if config.Setting.Plantuml == "" {
			return tracerr.New("plantuml library path is not specified")
		}
		if _, err := os.Stat(config.Setting.Plantuml); os.IsNotExist(err) {
			return tracerr.Wrap(err)
		}
		if err := txt.WriteWithLoader(loader, schema, config.Setting.Template, "output.puml"); err != nil {
			return tracerr.Wrap(err)
		}
		defer os.Remove("output.puml")
		if err := sh.Command("java", "-jar", config.Setting.Plantuml, "-o", filepath.Dir(config.Setting.Output), "output.puml").Run(); err != nil {
			return tracerr.Wrap(err)
		}
		logrus.Infof("Generated file: %s", config.Setting.Output)
	default:
		return tracerr.New(fmt.Sprintf("Output file extension '%s' is not supported", config.Setting.Output))
	}
	return nil
}
