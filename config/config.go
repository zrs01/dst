package config

import (
	mp "github.com/geraldo-labs/merge-struct"
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

type ConfigDef struct {
	Dsn     string // Global Data Source Name (ex: mariadb://root:DockerMySQL@tcp(10.24.64.157:10076)/mdis_dc_new_sit)
	Sdf     string // Global Schema Definition File
	Ccf     string // Global Common Column File
	DataDef SettingDef
	Erd     SettingDef
	Xls     SettingDef
	Export  SettingDef
	Text    SettingDef
}

type SettingDef struct {
	Dsn        string // Specific Data Source Name
	Sdf        string // Specific Schema Definition File
	Output     string // output file
	Template   string // template file
	TableName  string // table name filter
	ColumnName string // column name filter
	Plantuml   string // plantuml library path
	Ccf        string // common column file
	IsAlter    bool   // used int table creation
}

type ConfigType int

const (
	DataDefConf ConfigType = iota
	ERDConf
	ExcelConf
	ExportConf
	TextConf
)

var (
	Version  string     = "development"
	Filename string     = "config.yml"
	Default  SettingDef = SettingDef{}
)

// LoadConfig initializes the global Setting variable by merging configurations from multiple sources:
// 1. Loads base config from config.yml file
// 2. Merges specific settings based on configType (DataDef, Export, ERD, or Text)
// 3. Finally merges with provided source parameter to override any settings
func LoadConfig(configType ConfigType, source any) {
	configDef := &ConfigDef{}
	// load the configuration by default
	config := configor.New(&configor.Config{
		Silent: true, // suppress "file not found" error message
	})
	config.Load(configDef, Filename)

	Default.Dsn = configDef.Dsn
	Default.Sdf = configDef.Sdf
	Default.Ccf = configDef.Ccf

	switch configType {
	case DataDefConf:
		if _, err := mp.Struct(&Default, configDef.DataDef); err != nil {
			logrus.Fatal(err)
		}
	case ExcelConf:
		if _, err := mp.Struct(&Default, configDef.Xls); err != nil {
			logrus.Fatal(err)
		}
	case ERDConf:
		if _, err := mp.Struct(&Default, configDef.Erd); err != nil {
			logrus.Fatal(err)
		}
	case ExportConf:
		if _, err := mp.Struct(&Default, configDef.Export); err != nil {
			logrus.Fatal(err)
		}
	case TextConf:
		if _, err := mp.Struct(&Default, configDef.Text); err != nil {
			logrus.Fatal(err)
		}
	}
	if _, err := mp.Struct(&Default, source); err != nil {
		logrus.Fatal(err)
	}
}
