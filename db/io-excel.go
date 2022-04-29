package db

import (
	"fmt"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/xuri/excelize/v2"
)

func (s *Database) loadExcel(infile string) (*OutDB, error) {
	excel, err := excelize.OpenFile(infile)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to open the excel file %s", infile)
	}
	defer func() {
		if err := excel.Close(); err != nil {
			fmt.Println(eris.ToString(err, true))
		}
	}()

	var data OutDB

	sheets := excel.GetSheetList()
	for _, sheet := range sheets {
		rows, err := excel.GetRows(sheet)
		if err != nil {
			return nil, eris.Wrapf(err, "faled to get the rows from the sheet %s", sheet)
		}

		schema := &OutSchema{Name: sheet}

		var table *OutTable
		for rowIndex, row := range rows {
			if rowIndex == 0 {
				// skip the title row
				continue
			}

			if row[0] != "" {
				// new table (first column contains table name and descrition only)

				if table != nil {
					schema.Tables = append(schema.Tables, *table)
				}
				table = &OutTable{}

				tableInfo := row[0]
				descIndex := strings.Index(tableInfo, " - ")
				if descIndex != -1 {
					table.Name = tableInfo[0:descIndex]
					table.Desc = tableInfo[descIndex+3:]
				} else {
					table.Name = tableInfo
				}
			} else {
				// database columns start from column B
				table.Columns = append(table.Columns, OutColumn{Values: row[1:]})
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
		headings := []string{"Column Name", "Data Type", "Primary Key", "Unique", "Nullable", "Foreign Key", "Comment"}
		for i, heading := range headings {
			cell := fmt.Sprintf("%c1", 66+i) // start from column 2
			excel.SetCellValue(sheet, cell, heading)
		}
		// styling
		excel.SetCellStyle(sheet, "A1", fmt.Sprintf("%c1", 65+len(headings)), (*style)["header"])
		// column width
		widths := []float64{2, 20, 15, 12, 12, 12, 20, 50}
		for i, width := range widths {
			col := fmt.Sprintf("%c", 65+i)
			excel.SetColWidth(sheet, col, col, width)
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

			columns := s.convertToColumn(table.Columns)
			for i, column := range columns {
				index := i + offset
				excel.SetCellValue(sheet, fmt.Sprintf("B%d", index), column.Name)
				excel.SetCellValue(sheet, fmt.Sprintf("C%d", index), column.DataType)
				excel.SetCellValue(sheet, fmt.Sprintf("D%d", index), column.IsPK)
				excel.SetCellValue(sheet, fmt.Sprintf("E%d", index), column.IsUnique)
				excel.SetCellValue(sheet, fmt.Sprintf("F%d", index), column.Nullable)
				excel.SetCellValue(sheet, fmt.Sprintf("G%d", index), column.ForeignHint)
				excel.SetCellValue(sheet, fmt.Sprintf("H%d", index), column.Desc)
			}
			offset = offset + len(columns)
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

	header, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#b4c7dc"}, Pattern: 1},
	})
	if err != nil {
		return nil, eris.Wrap(err, "failed to create a header style")
	}
	style["header"] = header

	table, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#dee6ef"}, Pattern: 1},
	})
	if err != nil {
		return nil, eris.Wrap(err, "failed to create a table style")
	}
	style["table"] = table

	return &style, nil
}
