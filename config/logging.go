package config

import (
	"io"
	"os"
	"path/filepath"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type RollingFileDef struct {
	File       string `yaml:"file" default:"system.log"`
	MaxBackups int    `yaml:"maxBackups" default:"10"`
	MaxSize    int    `yaml:"maxSize" default:"5"` // MB
	MaxAge     int    `yaml:"maxAge" default:"30"` // days
	Compress   bool   `yaml:"compress" default:"true"`
}

func init() {
	SetLogFormatter()
	SetLogLevel()
}

func SetLogFormatter() {
	formatter := &nested.Formatter{
		NoColors:        false,
		NoFieldsColors:  false,
		HideKeys:        true,
		TimestampFormat: time.RFC3339,
	}
	if logrus.GetLevel() == logrus.InfoLevel {
		formatter.NoColors = true
		formatter.NoFieldsColors = true
	}
	logrus.SetFormatter(formatter)
}

func SetLogLevel() {
	level := logrus.InfoLevel
	l := os.Getenv("LOG_LEVEL")
	if l != "" {
		newLevel, err := logrus.ParseLevel(l)
		if err != nil {
			logrus.Error(err)
		}
		level = newLevel
	}
	logrus.SetLevel(level)
	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		logrus.SetReportCaller(true)
	}
}

func SetLogOutput(logfile string, config RollingFileDef) {
	var multiWriter io.Writer

	if logfile != "" {
		logfile = filepath.ToSlash(logfile)
		lumberjackLogger := &lumberjack.Logger{
			// Log file abbsolute path, os agnostic
			Filename:   logfile,
			MaxSize:    config.MaxSize, // MB
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,   // days
			Compress:   config.Compress, // disabled by default
		}
		logrus.Debugf("log file set to %s", logfile)

		// Fork writing into two outputs
		multiWriter = io.MultiWriter(os.Stderr, lumberjackLogger)
		logrus.SetOutput(multiWriter)
	}
}
