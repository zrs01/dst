package edit

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zrs01/dst/internal/service/edit/app"
	"github.com/zrs01/dst/internal/service/edit/dialog"
	"github.com/zrs01/dst/model"
)

type ColumnWidget struct {
	tableView  *tview.Table
	view       *tview.Flex
	tableModel *model.Table
}

func NewColumnWidget(tableModel *model.Table) *ColumnWidget {
	w := &ColumnWidget{
		tableModel: tableModel,
	}
	table := tview.NewTable()
	table.SetTitle("Columns")
	table.Clear()
	table.SetBorder(true)
	table.SetEvaluateAllRows(true)
	table.SetSelectable(true, false)
	table.SetInputCapture(w.keyHandler)
	// table.SetSelectedFunc(w.selectedHandler)

	for row := 0; row < len(tableModel.Columns); row++ {
		col := tableModel.Columns[row]
		tableNameCell := tview.NewTableCell(col.Name)
		tableDataTypeCell := tview.NewTableCell(col.DataType)
		tableIdentityCell := tview.NewTableCell(col.Identity)
		tableNotNullCell := tview.NewTableCell(col.NotNull)
		tableUniqueCell := tview.NewTableCell(col.Unique)
		tableValueCell := tview.NewTableCell(col.Value)
		tableForeignKeyCell := tview.NewTableCell(col.ForeignKey)
		tableCardinalityCell := tview.NewTableCell(col.Cardinality)
		tableTitleCell := tview.NewTableCell(col.Title)
		tableIndexCell := tview.NewTableCell(col.Index)
		tableDescCell := tview.NewTableCell(col.Desc)
		tableComputeCell := tview.NewTableCell(col.Compute)
		tableComputeTypeCell := tview.NewTableCell(col.ComputeType)
		table.SetCell(row, 0, tableNameCell)
		table.SetCell(row, 1, tableDataTypeCell)
		table.SetCell(row, 2, tableIdentityCell)
		table.SetCell(row, 3, tableNotNullCell)
		table.SetCell(row, 4, tableUniqueCell)
		table.SetCell(row, 5, tableValueCell)
		table.SetCell(row, 6, tableForeignKeyCell)
		table.SetCell(row, 7, tableCardinalityCell)
		table.SetCell(row, 8, tableTitleCell)
		table.SetCell(row, 9, tableIndexCell)
		table.SetCell(row, 10, tableDescCell)
		table.SetCell(row, 11, tableComputeCell)
		table.SetCell(row, 12, tableComputeTypeCell)
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

func (w *ColumnWidget) View() *tview.Flex {
	return w.view
}

func (w *ColumnWidget) keyHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEsc:
		id, _ := app.Container.GetFrontPage()
		app.Container.RemovePage(id)
	case tcell.KeyF2:
		row, _ := w.tableView.GetSelection()

		fields := []app.Field{
			{Label: "Name", Value: w.tableModel.Columns[row].Name},
			{Label: "Data Type", Value: w.tableModel.Columns[row].DataType},
			{Label: "Identity", Value: w.tableModel.Columns[row].Identity},
			{Label: "Not Null", Value: w.tableModel.Columns[row].NotNull},
			{Label: "Unique", Value: w.tableModel.Columns[row].Unique},
			{Label: "Value", Value: w.tableModel.Columns[row].Value},
			{Label: "Foreign Key", Value: w.tableModel.Columns[row].ForeignKey},
			{Label: "Cardinality", Value: w.tableModel.Columns[row].Cardinality},
			{Label: "Title", Value: w.tableModel.Columns[row].Title},
			{Label: "Index", Value: w.tableModel.Columns[row].Index},
			{Label: "Description", Value: w.tableModel.Columns[row].Desc},
			{Label: "Compute", Value: w.tableModel.Columns[row].Compute},
			{Label: "Compute Type", Value: w.tableModel.Columns[row].ComputeType},
		}
		dialog.NewDialog().WithTitle("Edit Column").WithWidth(100).Input(fields, func(fields []app.Field) {
			w.tableModel.Columns[row].Name = fields[0].Value
			w.tableModel.Columns[row].DataType = fields[1].Value
			w.tableModel.Columns[row].Identity = fields[2].Value
			w.tableModel.Columns[row].NotNull = fields[3].Value
			w.tableModel.Columns[row].Unique = fields[4].Value
			w.tableModel.Columns[row].Value = fields[5].Value
			w.tableModel.Columns[row].ForeignKey = fields[6].Value
			w.tableModel.Columns[row].Cardinality = fields[7].Value
			w.tableModel.Columns[row].Title = fields[8].Value
			w.tableModel.Columns[row].Index = fields[9].Value
			w.tableModel.Columns[row].Desc = fields[10].Value
			w.tableModel.Columns[row].Compute = fields[11].Value
			w.tableModel.Columns[row].ComputeType = fields[12].Value

			w.tableView.SetCell(row, 0, tview.NewTableCell(fields[0].Value))
			w.tableView.SetCell(row, 1, tview.NewTableCell(fields[1].Value))
			w.tableView.SetCell(row, 2, tview.NewTableCell(fields[2].Value))
			w.tableView.SetCell(row, 3, tview.NewTableCell(fields[3].Value))
			w.tableView.SetCell(row, 4, tview.NewTableCell(fields[4].Value))
			w.tableView.SetCell(row, 5, tview.NewTableCell(fields[5].Value))
			w.tableView.SetCell(row, 6, tview.NewTableCell(fields[6].Value))
			w.tableView.SetCell(row, 7, tview.NewTableCell(fields[7].Value))
			w.tableView.SetCell(row, 8, tview.NewTableCell(fields[8].Value))
			w.tableView.SetCell(row, 9, tview.NewTableCell(fields[9].Value))
			w.tableView.SetCell(row, 10, tview.NewTableCell(fields[10].Value))
			w.tableView.SetCell(row, 11, tview.NewTableCell(fields[11].Value))
			w.tableView.SetCell(row, 12, tview.NewTableCell(fields[12].Value))
		})
	}
	return event
}
