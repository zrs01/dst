package db

import (
	"fmt"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/xuri/excelize/v2"
)

func (s *Database) loadExcel(infile string) (*InDB, error) {
	excel, err := excelize.OpenFile(infile)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to open the excel file %s", infile)
	}
	defer func() {
		if err := excel.Close(); err != nil {
			fmt.Println(eris.ToString(err, true))
		}
	}()

	var data InDB

	sheets := excel.GetSheetList()
	for _, sheet := range sheets {
		rows, err := excel.GetRows(sheet)
		if err != nil {
			return nil, eris.Wrapf(err, "faled to get the rows from the sheet %s", sheet)
		}

		schema := &InSchema{Name: sheet}

		var table *InTable
		for rowIndex, row := range rows {
			if rowIndex == 0 {
				// skip the title row
				continue
			}

			if len(row) > 0 && row[0] != "" {
				// new table (first column contains table name and descrition only)

				if table != nil {
					schema.Tables = append(schema.Tables, *table)
				}
				table = &InTable{}

				tableInfo := row[0]
				descIndex := strings.Index(tableInfo, " - ")
				if descIndex != -1 {
					table.Name = tableInfo[0:descIndex]
					table.Desc = tableInfo[descIndex+3:]
				} else {
					table.Name = tableInfo
				}
			} else {
				incol := InColumn{}
				for idx, cell := range row {
					if idx == 1 {
						incol.Name = cell
					}
					if idx == 2 {
						incol.DataType = cell
					}
					if idx == 3 {
						incol.Identity = cell
					}
					if idx == 4 {
						incol.NotNull = cell
					}
					if idx == 5 {
						incol.Value = cell
					}
					if idx == 6 {
						incol.ForeignKeyHint = cell
					}
					if idx == 7 {
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
	return &data, nil
}

func (s *Database) saveExcel(data *InDB, outfile string) error {
	excel := excelize.NewFile()

	style, err := s.definedExcelStyle(excel)
	if err != nil {
		return eris.Wrap(err, "failed to create a bold style")
	}

	for _, schema := range data.Schemas {
		sheet := schema.Name
		excel.NewSheet(sheet)

		// heading
		headings := []string{"Column Name", "Data Type", "Identity", "Nullable", "Default", "Foreign Key", "Comment"}
		for i, heading := range headings {
			cell := fmt.Sprintf("%c1", 66+i) // start from column 2
			excel.SetCellValue(sheet, cell, heading)
		}
		// styling
		excel.SetCellStyle(sheet, "A1", fmt.Sprintf("%c1", 65+len(headings)), (*style)["header"])
		// column width
		widths := []float64{2, 20, 15, 8, 8, 15, 20, 50}
		for i, width := range widths {
			col := fmt.Sprintf("%c", 65+i)
			excel.SetColWidth(sheet, col, col, width)
		}

		var setColValue = func(rowIndex int, column InColumn) {
			excel.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), column.Name)
			excel.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), column.DataType)
			excel.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), column.Identity)
			excel.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), column.NotNull)
			excel.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), column.Value)
			excel.SetCellValue(sheet, fmt.Sprintf("G%d", rowIndex), column.ForeignKeyHint)
			excel.SetCellValue(sheet, fmt.Sprintf("H%d", rowIndex), column.Desc)
		}

		offset := 2

		for _, table := range schema.Tables {
			// table infomation
			tableText := table.Name
			if table.Desc != "" {
				tableText = tableText + " - " + table.Desc
			}
			tableCell := fmt.Sprintf("A%d", offset)
			excel.SetCellValue(sheet, tableCell, tableText)
			excel.SetCellStyle(sheet, tableCell, fmt.Sprintf("%c%d", 65+len(headings), offset), (*style)["table"])
			offset = offset + 1

			// table columns
			for i, column := range table.Columns {
				index := i + offset
				setColValue(index, column)
			}
			offset = offset + len(table.Columns)

			// fixed columns
			for i, column := range data.Fixed {
				index := i + offset
				setColValue(index, column)
				excel.SetCellStyle(sheet, fmt.Sprintf("B%d", index), fmt.Sprintf("%c%d", 65+len(headings), index), (*style)["fixcol"])
			}

			offset = offset + len(data.Fixed)
		}
	}
	excel.DeleteSheet("Sheet1")
	if err := excel.SaveAs(outfile); err != nil {
		return eris.Wrapf(err, "failed to save to file %s", outfile)
	}
	return nil

}
func (s *Database) definedExcelStyle(excel *excelize.File) (*map[string]int, error) {
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
