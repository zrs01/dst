package config

import (
	"github.com/jinzhu/configor"
)

type SettingDef struct {
	Input string
	Text  struct {
		Output       string
		Template     string
		SchemaFilter string
		TableFilter  string
		ColumnFilter string
	}
	Erd struct {
		Output       string `default:"output.png"`
		Template     string
		SchemaFilter string
		TableFilter  string
		UmlLib       string
	}
}

var (
	Debug    bool        = false
	Filename string      = "config.yml"
	Setting  *SettingDef = &SettingDef{}
)

func init() {
	// load the configuration by default
	config := configor.New(&configor.Config{
		Silent: true, // suppress "file not found" error message
	})
	config.Load(Setting, Filename)
}
