package edit

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TableWidget struct {
	tableWidget *tview.Table
	view        *tview.Flex
}

func NewTableWidget() *TableWidget {
	w := &TableWidget{}
	table := tview.NewTable()
	table.SetTitle("Tables")
	table.Clear()
	table.SetBorder(true)
	table.SetEvaluateAllRows(true)
	table.SetSelectable(true, false)
	table.SetInputCapture(w.keyHandler)
	table.SetSelectedFunc(w.selectedHandler)

	for row := 0; row < len(Schema.Tables); row++ {
		tb := Schema.Tables[row]
		tableNameCell := tview.NewTableCell(tb.Name)
		tableTitleCell := tview.NewTableCell(tb.Title)
		tableDescCell := tview.NewTableCell(tb.Desc)
		table.SetCell(row, 0, tableNameCell)
		table.SetCell(row, 1, tableTitleCell)
		table.SetCell(row, 2, tableDescCell)
	}
	w.tableWidget = table

	// Horizontal flex: 3 columns
	hFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 3, 0, false).  // left spacer
		AddItem(table, 0, 1, true). // centre widget
		AddItem(nil, 3, 0, false)   // right spacer

	// Vertical flex: 3 rows
	vFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 2, 0, false).  // top spacer
		AddItem(hFlex, 0, 1, true). // middle row (with centre widget)
		AddItem(nil, 2, 0, false)   // bottom spacer

	w.view = vFlex
	return w
}

func (w *TableWidget) View() *tview.Flex {
	return w.view
}

func (w *TableWidget) selectedHandler(row, column int) {
	fields := map[string]string{
		"Name":  "John",
		"Title": "John",
		"Desc":  "John",
	}
	InputDialog("Edit Table", fields, func(fields map[string]string) {
		// fmt.Printf("%s\n%s\n%s", fields["Name"], fields["Title"], fields["Desc"])
	})
}

func (w *TableWidget) keyHandler(event *tcell.EventKey) *tcell.EventKey {
	return event
}
