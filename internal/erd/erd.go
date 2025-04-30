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
	"github.com/zrs01/dst/internal/ddloader"
	"github.com/zrs01/dst/internal/tpl"
	"github.com/ztrue/tracerr"
)

//go:embed templates/default.jet
var templateFS embed.FS

func Generate() error {
	logrus.Info("Generating ER Diagram...")
	data, err := ddloader.Load(config.Setting.Input,
		ddloader.WithSchemaPattern(config.Setting.Erd.Schema),
		ddloader.WithTablePattern(config.Setting.Erd.Table))
	if err != nil {
		return tracerr.Wrap(err)
	}

	var loader jet.Loader
	// Use default template if not specified
	if config.Setting.Erd.Template == "" {
		logrus.Info("Template not provided, using default template")
		config.Setting.Erd.Template = "templates/default.jet"
		loader = embedfs.NewLoader(filepath.Dir(config.Setting.Erd.Template), templateFS)
	} else {
		logrus.Infof("Using template: %s", config.Setting.Erd.Template)
		loader = jet.NewOSFileSystemLoader(filepath.Dir(config.Setting.Erd.Template))
	}

	ext := filepath.Ext(config.Setting.Erd.Output)
	switch ext {
	case ".puml":
		if err := tpl.WriteWithLoader(loader, data, config.Setting.Erd.Template, config.Setting.Erd.Output); err != nil {
			return tracerr.Wrap(err)
		}
		logrus.Infof("Generated file: %s", config.Setting.Erd.Output)
	case ".png":
		if config.Setting.Erd.UmlLib == "" {
			return tracerr.New("plantuml library path is not specified")
		}
		if _, err := os.Stat(config.Setting.Erd.UmlLib); os.IsNotExist(err) {
			return tracerr.Wrap(err)
		}
		if err := tpl.WriteWithLoader(loader, data, config.Setting.Erd.Template, "output.puml"); err != nil {
			return tracerr.Wrap(err)
		}
		defer os.Remove("output.puml")
		if err := sh.Command("java", "-jar", config.Setting.Erd.UmlLib, "-o", filepath.Dir(config.Setting.Erd.Output), "output.puml").Run(); err != nil {
			return tracerr.Wrap(err)
		}
		logrus.Infof("Generated file: %s", config.Setting.Erd.Output)
	default:
		return tracerr.New(fmt.Sprintf("Output file extension '%s' is not supported", config.Setting.Erd.Output))
	}
	return nil
}
