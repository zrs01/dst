package edit

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/service"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

var (
	UIApp     *tview.Application
	Container *tview.Pages
	Schema    *model.Schema
)

// type UpdateCallback func(key tcell.Key) error

func Launch() error {
	if err := readData(); err != nil {
		return tracerr.Wrap(err)
	}
	tview.Styles.PrimitiveBackgroundColor = tcell.Color16 // 0x000000

	UIApp = tview.NewApplication()
	Container = tview.NewPages()

	baseWidget := NewTableWidget()

	Container.AddPage("root", baseWidget.View(), true, true)
	return UIApp.SetRoot(Container, true).SetFocus(Container).Run()
}

func readData() error {
	builder := service.NewLoadBuilder()
	logrus.Tracef("Reading data from %s", config.MergeSetting.Sdf)
	schema, err := builder.LoadFromFile(config.MergeSetting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	Schema = schema
	return nil
}
