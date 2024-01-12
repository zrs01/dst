package transform

import (
	"fmt"
	"strings"

	"github.com/shomali11/util/xstrings"
	"github.com/xuri/excelize/v2"
	"github.com/ztrue/tracerr"
)

const (
	CEmpty      = 0
	CName       = 1
	CTitle      = 2
	CDataType   = 3
	CIdentity   = 4
	CNotNull    = 5
	CValue      = 6
	CForeignKey = 7
	CDesc       = 8
)

func ReadXlsx(infile string) (*DataDef, error) {
	excel, err := excelize.OpenFile(infile)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var data DataDef

	sheets := excel.GetSheetList()
	for _, sheet := range sheets {
		rows, err := excel.GetRows(sheet)
		if err != nil {
			tracerr.Wrap(err)
		}

		schema := &Schema{Name: sheet}

		var table *Table
		for rowIndex, row := range rows {
			if rowIndex == 0 {
				// skip the heading row
				continue
			}

			if len(row) > 0 && row[0] != "" {
				// table title (first column contains table name and descrition only)

				if table != nil {
					// append last table instance
					schema.Tables = append(schema.Tables, *table)
				}
				table = &Table{}

				tableText := row[0]
				parts := strings.Split(tableText, " - ")
				switch len(parts) {
				case 0:
					return nil, tracerr.New("failed to get the table name")
				case 1:
					table.Name = strings.TrimSpace(tableText)
				case 2:
					table.Name = strings.TrimSpace(parts[0])
					table.Desc = strings.TrimSpace(parts[1])
				default:
					table.Name = strings.TrimSpace(parts[0])
					table.Title = strings.TrimSpace(parts[1])
					table.Desc = strings.TrimSpace(parts[2])
				}
			} else {
				incol := Column{}
				for idx, cell := range row {
					cell = strings.TrimSpace(cell)
					if idx == CName {
						incol.Name = cell
					}
					if idx == CTitle {
						incol.Title = cell
					}
					if idx == CDataType {
						incol.DataType = cell
					}
					if idx == CIdentity {
						incol.Identity = cell
					}
					if idx == CNotNull {
						incol.NotNull = cell
					}
					if idx == CValue {
						incol.Value = cell
					}
					if idx == CForeignKey {
						incol.ForeignKey = cell
					}
					if idx == CDesc {
						incol.Desc = cell
					}
				}
				table.OutColumns = append(table.OutColumns, OutColumn{Value: incol})
			}
		}
		// last table
		schema.Tables = append(schema.Tables, *table)
		data.Schemas = append(data.Schemas, *schema)
	}
	return &data, err
}

func WriteXlsx(data *DataDef, out string, simple bool) error {
	if simple {
		return writeSimpleDataDict(data, out)
	}
	return writeDataDict(data, out)
}

func writeSimpleDataDict(data *DataDef, out string) error {
	excel := excelize.NewFile()
	sheet := "Tables Description"
	excel.NewSheet(sheet)

	style, err := definedExcelStyle(excel)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// heading
	headings := []string{"Table ID", "Title Name", "Table Description", "Key Data Item"}
	for i, heading := range headings {
		cell := fmt.Sprintf("%c1", 65+i)
		excel.SetCellValue(sheet, cell, heading)
	}
	// styling
	excel.SetCellStyle(sheet, "A1", fmt.Sprintf("%c1", 65+len(headings)-1), (*style)["header"])
	// column width
	widths := []float64{8, 15, 60, 30}
	for i, width := range widths {
		col := fmt.Sprintf("%c", 65+i)
		excel.SetColWidth(sheet, col, col, width)
	}

	rowctnr := 2 // row counter
	for i, schema := range data.Schemas {
		// schema name
		excel.SetCellValue(sheet, fmt.Sprintf("A%d", rowctnr), schema.Name)
		for j := 0; j < len(headings); j++ {
			excel.SetCellStyle(sheet, fmt.Sprintf("A%d", rowctnr), fmt.Sprintf("%c%d", 65+len(headings)-1, rowctnr), (*style)["table"])
		}

		for j, table := range schema.Tables {
			rowctnr += 1
			excel.SetCellValue(sheet, fmt.Sprintf("A%d", rowctnr), fmt.Sprintf("T%d", (i+1)*100+j))
			excel.SetCellValue(sheet, fmt.Sprintf("B%d", rowctnr), table.Name)
			excel.SetCellValue(sheet, fmt.Sprintf("C%d", rowctnr), table.Desc)

			// find PK and all FK of the table
			keyData := ""
			keyDataRow := 0
			for k, column := range table.Columns {
				if strings.ToUpper(column.Identity) == "Y" || xstrings.IsNotBlank(column.ForeignKey) {
					if k != 0 {
						keyData += "\n"
					}
					if strings.ToUpper(column.Identity) == "Y" {
						keyData += fmt.Sprintf("%s (%s)", column.Name, "PK")
						keyDataRow += 1
					} else if xstrings.IsNotBlank(column.ForeignKey) {
						keyData += fmt.Sprintf("%s (%s)", column.Name, "FK")
						keyDataRow += 1
					}
				}
			}
			excel.SetCellValue(sheet, fmt.Sprintf("D%d", rowctnr), keyData)

			h, _ := excel.GetRowHeight(sheet, 1)
			excel.SetRowHeight(sheet, rowctnr, h*float64(keyDataRow))
		}
		rowctnr += 1
	}
	excel.DeleteSheet("Sheet1")
	if err := excel.SaveAs(out); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func writeDataDict(data *DataDef, out string) error {
	excel := excelize.NewFile()

	style, err := definedExcelStyle(excel)
	if err != nil {
		return tracerr.Wrap(err)
	}

	for _, schema := range data.Schemas {
		sheet := schema.Name
		excel.NewSheet(sheet)

		// heading
		headings := []string{"Column Name", "Title", "Data Type", "Identity", "Not Null", "Default", "Foreign Key", "Description"}
		for i, heading := range headings {
			cell := fmt.Sprintf("%c1", 66+i) // start from column 2
			excel.SetCellValue(sheet, cell, heading)
		}
		// styling
		excel.SetCellStyle(sheet, "A1", fmt.Sprintf("%c1", 65+len(headings)), (*style)["header"])
		// column width
		widths := []float64{2, 20, 20, 15, 8, 8, 10, 25, 50}
		for i, width := range widths {
			col := fmt.Sprintf("%c", 65+i)
			excel.SetColWidth(sheet, col, col, width)
		}

		var setColValue = func(rowIndex int, column Column) {
			excel.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), column.Name)
			excel.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), column.Title)
			excel.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), column.DataType)
			excel.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), column.Identity)
			excel.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), column.NotNull)
			excel.SetCellValue(sheet, fmt.Sprintf("G%d", rowIndex), column.Value)
			excel.SetCellValue(sheet, fmt.Sprintf("H%d", rowIndex), column.ForeignKey)
			excel.SetCellValue(sheet, fmt.Sprintf("I%d", rowIndex), column.Desc)
		}

		rowctnr := 2 // row counter

		for _, table := range schema.Tables {
			{
				// table infomation
				text := table.Name
				if table.Title != "" {
					text += " - " + table.Title
				}
				if table.Desc != "" {
					text += " - " + table.Desc
				}
				cell := fmt.Sprintf("A%d", rowctnr)
				excel.SetCellValue(sheet, cell, text)
				excel.SetCellStyle(sheet, cell, fmt.Sprintf("%c%d", 65+len(headings), rowctnr), (*style)["table"])
				rowctnr += 1
			}

			// table columns
			for i, column := range table.Columns {
				index := i + rowctnr
				setColValue(index, column)
			}
			rowctnr += len(table.Columns)

			// fixed columns
			for i, column := range data.Fixed {
				index := i + rowctnr
				setColValue(index, column)
				excel.SetCellStyle(sheet, fmt.Sprintf("B%d", index), fmt.Sprintf("%c%d", 65+len(headings), index), (*style)["fixcol"])
			}

			rowctnr += len(data.Fixed)
		}
	}
	excel.DeleteSheet("Sheet1")
	if err := excel.SaveAs(out); err != nil {
		return tracerr.Wrap(err)
	}
	return nil

}
func definedExcelStyle(excel *excelize.File) (*map[string]int, error) {
	style := make(map[string]int, 0)

	// style for the header cell
	header, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#b4c7dc"}, Pattern: 1},
	})
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	style["header"] = header

	// style for the table name cell
	table, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#dee6ef"}, Pattern: 1},
	})
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	style["table"] = table

	// style for the fix column cell
	fixColStyle, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "#205375"},
	})
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	style["fixcol"] = fixColStyle

	return &style, nil
}
