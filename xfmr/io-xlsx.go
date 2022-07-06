package xfmr

import (
	"fmt"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/xuri/excelize/v2"
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

func (s *Xfmr) LoadXlsx(infile string) error {
	excel, err := excelize.OpenFile(infile)
	if err != nil {
		return eris.Wrapf(err, "failed to open %s", infile)
	}

	var data InDB

	sheets := excel.GetSheetList()
	for _, sheet := range sheets {
		rows, err := excel.GetRows(sheet)
		if err != nil {
			eris.Wrapf(err, "failed to get the rows from the sheet %s", sheet)
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
					return eris.New("failed to get the table name")
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
				table.OutColumns = append(table.OutColumns, OutColumn{Values: incol})
			}
		}
		// last table
		schema.Tables = append(schema.Tables, *table)
		data.Schemas = append(data.Schemas, *schema)
	}
	s.Data = &data
	return nil
}

func (s *Xfmr) SaveToXlsx(outfile string) error {
	excel := excelize.NewFile()

	style, err := s.definedExcelStyle(excel)
	if err != nil {
		return eris.Wrap(err, "failed to create a bold style")
	}

	for _, schema := range s.Data.Schemas {
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
			for i, column := range s.Data.Fixed {
				index := i + rowctnr
				setColValue(index, column)
				excel.SetCellStyle(sheet, fmt.Sprintf("B%d", index), fmt.Sprintf("%c%d", 65+len(headings), index), (*style)["fixcol"])
			}

			rowctnr += len(s.Data.Fixed)
		}
	}
	excel.DeleteSheet("Sheet1")
	if err := excel.SaveAs(outfile); err != nil {
		return eris.Wrapf(err, "failed to save to file %s", outfile)
	}
	return nil

}
func (s *Xfmr) definedExcelStyle(excel *excelize.File) (*map[string]int, error) {
	style := make(map[string]int, 0)

	// style for the header cell
	header, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#b4c7dc"}, Pattern: 1},
	})
	if err != nil {
		return nil, eris.Wrap(err, "failed to create a header style")
	}
	style["header"] = header

	// style for the table name cell
	table, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#dee6ef"}, Pattern: 1},
	})
	if err != nil {
		return nil, eris.Wrap(err, "failed to create a table style")
	}
	style["table"] = table

	// style for the fix column cell
	fixColStyle, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "#205375"},
	})
	if err != nil {
		return nil, eris.Wrap(err, "failed to create a fixed column style")
	}
	style["fixcol"] = fixColStyle

	return &style, nil
}
