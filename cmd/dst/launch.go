package dst

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/ztrue/tracerr"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/zrs01/dst/config"
)

func Launch() {
	config.SetLogFormatter()
	config.SetLogLevel()

	cmd := &cli.Command{
		Name:    "dst",
		Usage:   "database schema tool",
		Version: config.Version,
	}

	cmd.Commands = append(cmd.Commands, RegisterDDLCmd())    // SQL DDL statement
	cmd.Commands = append(cmd.Commands, RegisterERDCmd())    // ER diagram
	cmd.Commands = append(cmd.Commands, RegisterXlsCmd())    // Excel
	cmd.Commands = append(cmd.Commands, RegisterExportCmd()) // Export
	cmd.Commands = append(cmd.Commands, RegisterTextCmd())   // Template

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		if logrus.GetLevel() >= logrus.DebugLevel {
			logrus.Error(tracerr.SprintSourceColor(err, 0))
		} else {
			logrus.Error(tracerr.Sprint(err))
		}
		os.Exit(1)
	}
}
