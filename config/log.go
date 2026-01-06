package config

import (
	"os"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

func SetLogFormatter() {
	logrus.SetFormatter(&nested.Formatter{
		NoColors:        false,
		NoFieldsColors:  false,
		HideKeys:        true,
		TimestampFormat: time.RFC3339,
	})
}

func SetLogLevel() {
	lvlstr := os.Getenv("LOG_LEVEL")
	if lvlstr == "" {
		lvlstr = "info"
	}
	lvl, err := logrus.ParseLevel(lvlstr)
	if err != nil {
		logrus.Error(err)
	}
	logrus.SetLevel(lvl)

	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		logrus.SetReportCaller(true)
	}
}
