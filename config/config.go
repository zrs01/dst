package config

import (
	"github.com/jinzhu/configor"
)

type SettingDef struct {
	Export struct {
		Dsn              string
		CommonColumnFile string `yaml:"ccf"`
		Output           string
		TableFilter      string
	}
	Text struct {
		Input        string
		Output       string
		Template     string
		SchemaFilter string
		TableFilter  string
		ColumnFilter string
	}
	Erd struct {
		Input        string
		Output       string `default:"output.png"`
		Template     string
		SchemaFilter string
		TableFilter  string
		PlantumlLib  string `yaml:"plantuml"`
	}
}

var (
	Version  string      = "development"
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
