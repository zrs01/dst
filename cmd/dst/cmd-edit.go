package dst

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/flagbuilder"
	"github.com/zrs01/dst/internal/service/edit"
	"github.com/ztrue/tracerr"
)

func RegisterEditCmd() *cli.Command {
	var options config.SettingDef

	return &cli.Command{
		Name:  "edit",
		Usage: "edit the schema definition file",
		Flags: []cli.Flag{
			flagbuilder.SdfFlag().WithDestination(&options.Sdf).Build(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := config.LoadConfig(config.EmptyConf, options); err != nil {
				return tracerr.Wrap(err)
			}
			return edit.Launch()
		},
	}
}
