package config

import (
	"github.com/jinzhu/configor"
)

type SettingDef struct {
	Text struct {
		Input    string
		Output   string
		Template string
		Schema   string
		Table    string
	}
	Erd struct {
		Input    string
		Output   string `default:"output.png"`
		Template string
		Schema   string
		Table    string
		UmlLib   string
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
