package config

import (
	mp "github.com/geraldo-labs/merge-struct"
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

type ConfigDef struct {
	Dsn     string // default database source name
	Sdf     string `yaml:"schemaDefinitionFile"` // default input file
	Ccf     string `yaml:"commonColumnFile"`     // default common column file
	DataDef SettingDef
	Erd     SettingDef
	Xls     SettingDef
	Export  SettingDef
	Text    SettingDef
}

type SettingDef struct {
	Dsn        string // database source name
	Sdf        string `yaml:"schemaDefinitionFile"` // schema definition file
	Output     string // output file
	Template   string // template file
	TableName  string // table name filter
	ColumnName string // column name filter
	Plantuml   string // plantuml library path
	Ccf        string `yaml:"commonColumnFile"` // common column file
	IsAlter    bool   // used int table creation
}

type ConfigType int

const (
	DataDefConf = iota
	ERDConf
	ExcelConf
	ExportConf
	TextConf
)

var (
	Version  string     = "development"
	Debug    bool       = false
	Filename string     = "config.yml"
	Setting  SettingDef = SettingDef{}
)

// InitSetting initializes the global Setting variable by merging configurations from multiple sources:
// 1. Loads base config from config.yml file
// 2. Merges specific settings based on configType (DataDef, Export, ERD, or Text)
// 3. Finally merges with provided source parameter to override any settings
func InitSetting(configType ConfigType, source any) {
	configDef := &ConfigDef{}
	// load the configuration by default
	config := configor.New(&configor.Config{
		Silent: true, // suppress "file not found" error message
	})
	config.Load(configDef, Filename)

	Setting.Dsn = configDef.Dsn
	Setting.Sdf = configDef.Sdf
	Setting.Ccf = configDef.Ccf

	switch configType {
	case DataDefConf:
		if _, err := mp.Struct(&Setting, configDef.DataDef); err != nil {
			logrus.Fatal(err)
		}
	case ExcelConf:
		if _, err := mp.Struct(&Setting, configDef.Xls); err != nil {
			logrus.Fatal(err)
		}
	case ERDConf:
		if _, err := mp.Struct(&Setting, configDef.Erd); err != nil {
			logrus.Fatal(err)
		}
	case ExportConf:
		if _, err := mp.Struct(&Setting, configDef.Export); err != nil {
			logrus.Fatal(err)
		}
	case TextConf:
		if _, err := mp.Struct(&Setting, configDef.Text); err != nil {
			logrus.Fatal(err)
		}
	}
	if _, err := mp.Struct(&Setting, source); err != nil {
		logrus.Fatal(err)
	}
}
