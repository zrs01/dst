package edit

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/service"
	"github.com/zrs01/dst/internal/service/edit/app"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

// type UpdateCallback func(key tcell.Key) error

func Launch() error {
	schema, err := readData()
	if err != nil {
		return tracerr.Wrap(err)
	}
	tview.Styles.ContrastBackgroundColor = tcell.ColorGray
	tview.Styles.InverseTextColor = tcell.ColorGray
	tview.Styles.PrimitiveBackgroundColor = tcell.Color16 // 0x000000

	app.UIApp = tview.NewApplication()
	app.Container = tview.NewPages()

	baseWidget := NewTableWidget(schema)

	app.Container.AddPage("root", baseWidget.View(), true, true)
	return app.UIApp.SetRoot(app.Container, true).SetFocus(app.Container).Run()
}

func readData() (*model.Schema, error) {
	builder := service.NewLoadBuilder()
	logrus.Tracef("Reading data from %s", config.MergeSetting.Sdf)
	schema, err := builder.LoadFromFile(config.MergeSetting.Sdf)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return schema, nil
}
