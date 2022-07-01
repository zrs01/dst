package xfmr

import (
	"fmt"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/xuri/excelize/v2"
)

func (s *Xfmr) LoadExcel(infile string) {
	excel, err := excelize.OpenFile(infile)
	if err != nil {
		panic(fmt.Sprintf("failed to open the excel file %s", infile))
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
			panic(fmt.Sprintf("faled to get the rows from the sheet %s", sheet))
		}

		schema := &Schema{Name: sheet}

		var table *Table
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
				table = &Table{}

				tableInfo := row[0]
				descIndex := strings.Index(tableInfo, " - ")
				if descIndex != -1 {
					table.Name = tableInfo[0:descIndex]
					table.Desc = tableInfo[descIndex+3:]
				} else {
					table.Name = tableInfo
				}
			} else {
				incol := Column{}
				for idx, cell := range row {
					if idx == CName {
						incol.Name = cell
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
					if idx == CUnique {
						incol.Unique = cell
					}
					if idx == CValue {
						incol.Value = cell
					}
					if idx == CForeignKeyHint {
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
}

func (s *Xfmr) SaveToExcel(outfile string) error {
	excel := excelize.NewFile()

	style, err := s.definedExcelStyle(excel)
	if err != nil {
		return eris.Wrap(err, "failed to create a bold style")
	}

	for _, schema := range s.Data.Schemas {
		sheet := schema.Name
		excel.NewSheet(sheet)

		// heading
		headings := []string{"Column Name", "Data Type", "Identity", "Not Null", "Default", "Foreign Key", "Comment"}
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

		var setColValue = func(rowIndex int, column Column) {
			excel.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), column.Name)
			excel.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), column.DataType)
			excel.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), column.Identity)
			excel.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), column.NotNull)
			excel.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), column.Value)
			excel.SetCellValue(sheet, fmt.Sprintf("G%d", rowIndex), column.ForeignKey)
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
			for i, column := range s.Data.Fixed {
				index := i + offset
				setColValue(index, column)
				excel.SetCellStyle(sheet, fmt.Sprintf("B%d", index), fmt.Sprintf("%c%d", 65+len(headings), index), (*style)["fixcol"])
			}

			offset = offset + len(s.Data.Fixed)
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
