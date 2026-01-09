package edit

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zrs01/dst/internal/service/edit/app"
	"github.com/zrs01/dst/internal/service/edit/dialog"
	"github.com/zrs01/dst/model"
)

type TableWidget struct {
	tableView   *tview.Table
	view        *tview.Flex
	schemaModel *model.Schema
}

func NewTableWidget(schemaModel *model.Schema) *TableWidget {
	w := &TableWidget{
		schemaModel: schemaModel,
	}

	table := tview.NewTable()
	table.SetTitle("Tables")
	table.Clear()
	table.SetBorder(true)
	table.SetEvaluateAllRows(true)
	table.SetSelectable(true, false)
	table.SetInputCapture(w.keyHandler)
	table.SetSelectedFunc(w.selectedHandler)

	for row := 0; row < len(schemaModel.Tables); row++ {
		tb := schemaModel.Tables[row]
		tableNameCell := tview.NewTableCell(tb.Name)
		tableTitleCell := tview.NewTableCell(tb.Title)
		tableDescCell := tview.NewTableCell(tb.Desc)
		table.SetCell(row, 0, tableNameCell)
		table.SetCell(row, 1, tableTitleCell)
		table.SetCell(row, 2, tableDescCell)
	}
	w.tableView = table

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
	table := w.schemaModel.Tables[row]
	columnWidget := NewColumnWidget(table)
	app.Container.AddPage(table.Name, columnWidget.View(), true, true)
}

func (w *TableWidget) keyHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyF2:
		row, _ := w.tableView.GetSelection()

		fields := []app.Field{
			{Label: "Name", Value: w.schemaModel.Tables[row].Name},
			{Label: "Title", Value: w.schemaModel.Tables[row].Title},
			{Label: "Description", Value: w.schemaModel.Tables[row].Desc},
		}
		dialog.NewDialog().WithTitle("Edit Table").Input(fields, func(fields []app.Field) {
			w.schemaModel.Tables[row].Name = fields[0].Value
			w.schemaModel.Tables[row].Title = fields[1].Value
			w.schemaModel.Tables[row].Desc = fields[2].Value

			w.tableView.SetCell(row, 0, tview.NewTableCell(fields[0].Value))
			w.tableView.SetCell(row, 1, tview.NewTableCell(fields[1].Value))
			w.tableView.SetCell(row, 2, tview.NewTableCell(fields[2].Value))
		})
	}
	return event
}
