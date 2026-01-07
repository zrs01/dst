package config

import (
	mp "github.com/geraldo-labs/merge-struct"
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
	"github.com/ztrue/tracerr"
)

type ConfigDef struct {
	Logging struct {
		Output  string         `yaml:"output"`
		Rolling RollingFileDef `yaml:"rolling"`
	} `yaml:"logging"`
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
	EmptyConf
)

var (
	Version      string = "development"
	Filename     string = "config.yml"
	Default      ConfigDef
	MergeSetting SettingDef
)

// LoadConfig initializes the global Setting variable by merging configurations from multiple sources:
// 1. Loads base config from config.yml file
// 2. Merges specific settings based on configType (DataDef, Export, ERD, or Text)
// 3. Finally merges with provided source parameter to override any settings
func LoadConfig(configType ConfigType, source any) error {
	// load the configuration by default
	config := configor.New(&configor.Config{
		Silent: true, // suppress "file not found" error message
	})
	logrus.Tracef("Loading config file: %s", Filename)
	if err := config.Load(&Default, Filename); err != nil {
		return tracerr.Wrap(err)
	}
	SetLogOutput(Default.Logging.Output, Default.Logging.Rolling)

	MergeSetting.Dsn = Default.Dsn
	MergeSetting.Sdf = Default.Sdf
	MergeSetting.Ccf = Default.Ccf

	switch configType {
	case DataDefConf:
		if _, err := mp.Struct(&MergeSetting, Default.DataDef); err != nil {
			return tracerr.Wrap(err)
		}
	case ExcelConf:
		if _, err := mp.Struct(&MergeSetting, Default.Xls); err != nil {
			return tracerr.Wrap(err)
		}
	case ERDConf:
		if _, err := mp.Struct(&MergeSetting, Default.Erd); err != nil {
			return tracerr.Wrap(err)
		}
	case ExportConf:
		if _, err := mp.Struct(&MergeSetting, Default.Export); err != nil {
			return tracerr.Wrap(err)
		}
	case TextConf:
		if _, err := mp.Struct(&MergeSetting, Default.Text); err != nil {
			return tracerr.Wrap(err)
		}
	}
	if _, err := mp.Struct(&MergeSetting, source); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}
