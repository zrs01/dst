package edit

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/service"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

var (
	UIApp  *tview.Application
	pages  *tview.Pages
	Schema *model.Schema
)

func Launch() error {
	if err := readData(); err != nil {
		return tracerr.Wrap(err)
	}
	tview.Styles.PrimitiveBackgroundColor = tcell.Color16 // 0x000000

	UIApp = tview.NewApplication()
	pages = tview.NewPages()
	pages.AddPage("root", TableWidget(), true, true)
	return UIApp.SetRoot(pages, true).SetFocus(pages).Run()
}

func readData() error {
	builder := service.NewLoadBuilder()
	schema, err := builder.LoadFromFile(config.MergeSetting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	Schema = schema
	return nil
}
