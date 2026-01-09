package dialog

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"github.com/zrs01/dst/internal/service/edit/app"
	"github.com/ztrue/tracerr"
)

type Dialog struct {
	title string
	width int
}

func (d *Dialog) WithTitle(t string) *Dialog {
	d.title = t
	return d
}

func (d *Dialog) WithWidth(w int) *Dialog {
	d.width = w
	return d
}

func NewDialog() *Dialog {
	return &Dialog{
		width: 80,
	}
}

func ConfirmDialog(text string, doneFunc func(ok bool) error) {
	pid := uuid.New().String()
	modalView := tview.NewModal()
	modalView.SetText(text)
	modalView.AddButtons([]string{"Yes", "No"})
	modalView.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		var err error
		if buttonLabel == "Yes" {
			err = doneFunc(true)
		}
		if buttonLabel == "No" {
			err = doneFunc(false)
		}
		if err != nil {
			ErrorDialog(err, nil)
			return
		}
		app.Container.RemovePage(pid)
	})
	app.Container.AddPage(pid, modalView, false, true)
}

func NotifyDialog(text string, doneFunc func()) {
	pid := uuid.New().String()
	modalView := tview.NewModal()
	modalView.SetText(text)
	modalView.AddButtons([]string{"OK"})
	modalView.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.Container.RemovePage(pid)
		if doneFunc != nil {
			doneFunc()
		}
	})
	app.Container.AddPage(pid, modalView, false, true)
}

func ErrorDialog(err error, doneFunc func()) {
	if _, ok := err.(tracerr.Error); ok {
		traceErr := err.(tracerr.Error)
		if len(traceErr.StackTrace()) > 0 {
			stackTrace(tracerr.SprintSource(err, 0), doneFunc)
			return
		}
	}

	pid := uuid.New().String()
	modalView := tview.NewModal()
	modalView.SetText(err.Error())
	modalView.SetBackgroundColor(tcell.ColorRed)
	modalView.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorRed))
	modalView.AddButtons([]string{"OK"})
	modalView.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.Container.RemovePage(pid)
		if doneFunc != nil {
			doneFunc()
		}
	})
	app.Container.AddPage(pid, modalView, false, true)
}

func stackTrace(text string, doneFunc func()) {
	pid := uuid.New().String()

	textView := tview.NewTextView()
	textView.SetDynamicColors(true)
	textView.SetText(text)
	textView.SetDisabled(true)
	textView.SetBackgroundColor(tcell.ColorRed)
	textView.SetTextColor(tcell.ColorWhite)

	container := tview.NewFlex()
	container.SetDirection(tview.FlexRow)
	container.SetTitle(" Error ")
	container.SetBorder(true)
	container.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorRed))
	container.AddItem(textView, 0, 1, false)
	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.Container.RemovePage(pid)
			if doneFunc != nil {
				doneFunc()
			}
			return nil
		}
		return event
	})
	app.Container.AddPage(pid, container, true, true)
}

func DebugDialog(text string) {
	pid := uuid.New().String()
	textView := tview.NewTextView()
	textView.SetText(text)
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.Container.RemovePage(pid)
		}
		return event
	})
	app.Container.AddPage(pid, textView, true, true)
}

// InputDialog displays a dialog with input fields for the user to enter data.
func (d *Dialog) Input(fields []app.Field, callback func([]app.Field)) {
	pid := uuid.New().String()

	c1 := tview.NewForm()
	c1.SetBorderPadding(1, 1, 2, 2)

	for key, value := range fields {
		c1.AddInputField(value.Label, value.Value, 0, nil, func(text string) {
			fields[key] = app.Field{Label: value.Label, Value: text}
		})
	}

	c1.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			callback(fields)
			app.Container.RemovePage(pid)
		case tcell.KeyEscape:
			app.Container.RemovePage(pid)
		}
		return event
	})

	frame := tview.NewFlex()
	frame.SetDirection(tview.FlexRow)
	frame.SetBorder(true)
	frame.SetTitle(fmt.Sprintf(" %s ", d.title))
	frame.AddItem(c1, 0, 1, true)

	view := app.Center(d.width, len(fields)*2+3, frame)
	app.Container.AddPage(pid, view, true, true)
}
