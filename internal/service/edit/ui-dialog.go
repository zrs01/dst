package edit

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"github.com/ztrue/tracerr"
)

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
		Container.RemovePage(pid)
	})
	Container.AddPage(pid, modalView, false, true)
}

func NotifyDialog(text string, doneFunc func()) {
	pid := uuid.New().String()
	modalView := tview.NewModal()
	modalView.SetText(text)
	modalView.AddButtons([]string{"OK"})
	modalView.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		Container.RemovePage(pid)
		if doneFunc != nil {
			doneFunc()
		}
	})
	Container.AddPage(pid, modalView, false, true)
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
		Container.RemovePage(pid)
		if doneFunc != nil {
			doneFunc()
		}
	})
	Container.AddPage(pid, modalView, false, true)
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
			Container.RemovePage(pid)
			if doneFunc != nil {
				doneFunc()
			}
			return nil
		}
		return event
	})
	Container.AddPage(pid, container, true, true)
}

func DebugDialog(text string) {
	pid := uuid.New().String()
	textView := tview.NewTextView()
	textView.SetText(text)
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			Container.RemovePage(pid)
		}
		return event
	})
	Container.AddPage(pid, textView, true, true)
}

// InputDialog displays a dialog with input fields for the user to enter data.
// Example:
//
//	InputDialog("Enter your name", map[string]string{"Name": ""}, func(fields map[string]string) {
//	    fmt.Println("Name:", fields["Name"])
//	})
func InputDialog(title string, fields map[string]string, callback func(map[string]string)) {
	pid := uuid.New().String()

	c1 := tview.NewForm()
	c1.SetBorderPadding(1, 1, 2, 2)
	c1.SetFieldBackgroundColor(tcell.ColorGray)

	for key, value := range fields {
		c1.AddInputField(key, value, 0, nil, func(text string) {
			fields[key] = text
		})
	}

	c1.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			callback(fields)
			Container.RemovePage(pid)
		case tcell.KeyEscape:
			Container.RemovePage(pid)
		}
		return event
	})

	// c2 := tview.NewTextView()
	// c2.SetDynamicColors(true)
	// c2.SetBackgroundColor(tcell.Color(0x16))
	// c2.SetText(`  [blue]F2:Confirm  ESC:Exit`)

	frame := tview.NewFlex()
	frame.SetDirection(tview.FlexRow)
	frame.SetBorder(true)
	frame.SetTitle(fmt.Sprintf(" %s ", title))
	frame.AddItem(c1, 0, 1, true)
	// frame.AddItem(c2, 1, 1, true)

	view := Center(80, len(fields)*2+3, frame)
	Container.AddPage(pid, view, true, true)
}
