package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/cmd/dst"
	"github.com/zrs01/dst/config"
	"github.com/ztrue/tracerr"
)

var version = "development"

func main() {
	config.Version = version
	config.SetLogFormatter()
	config.SetLogLevel()

	cmd := &cli.Command{
		Name:    "dst",
		Usage:   "database schema tool",
		Version: config.Version,
	}

	cmd.Commands = append(cmd.Commands, dst.RegisterDDLCmd())    // SQL DDL statement
	cmd.Commands = append(cmd.Commands, dst.RegisterERDCmd())    // ER diagram
	cmd.Commands = append(cmd.Commands, dst.RegisterXlsCmd())    // Excel
	cmd.Commands = append(cmd.Commands, dst.RegisterExportCmd()) // Export
	cmd.Commands = append(cmd.Commands, dst.RegisterTextCmd())   // Template
	cmd.Commands = append(cmd.Commands, dst.RegisterEditCmd())   // Edit

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		if logrus.GetLevel() >= logrus.DebugLevel {
			logrus.Error(tracerr.SprintSourceColor(err, 0))
		} else {
			logrus.Error(tracerr.Sprint(err))
		}
		os.Exit(1)
	}
}
