package dst

import (
	"context"
	"os"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/ztrue/tracerr"

	nested "github.com/antonfisher/nested-logrus-formatter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/zrs01/dst/config"
)

func Launch() {
	initLogrus()

	cmd := &cli.Command{
		Name:    "dst",
		Usage:   "database schema tool",
		Version: config.Version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "debug mode",
				Required:    false,
				Destination: &config.Debug,
			},
		},
	}

	cmd.Commands = append(cmd.Commands, RegisterDLLCmd())    // SQL DDL statement
	cmd.Commands = append(cmd.Commands, RegisterERDCmd())    // ER diagram
	cmd.Commands = append(cmd.Commands, RegisterExcelCmd())  // Excel
	cmd.Commands = append(cmd.Commands, RegisterExportCmd()) // Export
	cmd.Commands = append(cmd.Commands, RegisterTextCmd())   // Template

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		if config.Debug {
			logrus.Error(tracerr.SprintSourceColor(err, 0))
		} else {
			logrus.Errorf("%s", err)
		}
	}
}

func initLogrus() {
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: time.RFC3339,
	})
}
