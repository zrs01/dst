package edit

import "github.com/rivo/tview"

func TableWidget() *tview.Flex {
	table := tview.NewTable()
	table.SetBorder(true)

	for row := 0; row < len(Schema.Tables); row++ {
		tb := Schema.Tables[row]
		tableNameCell := tview.NewTableCell(tb.Name)
		table.SetCell(row, 0, tableNameCell)
	}

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
	return vFlex
}
