package db

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v2"
)

type InData struct {
	Schemas []InSchema `yaml:"schemas"`
}

type InSchema struct {
	Name   string    `yaml:"name" default:"Schema"`
	Desc   string    `yaml:"description"`
	Tables []InTable `yaml:"tables"`
}

type InTable struct {
	Name    string     `yaml:"name"`
	Desc    string     `yaml:"description"`
	Columns []InColumn `yaml:"columns"`
}

type InColumn struct {
	Values []string `yaml:"columnValues,flow"`
}

func (s *Database) ExportToExcel(f string) error {
	excel := excelize.NewFile()

	// define global style
	boldStyle, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return eris.Wrap(err, "failed to create a bold style")
	}

	for _, schema := range s.Data.Schemas {
		sheet := schema.Name
		excel.NewSheet(sheet)

		// heading
		headings := []string{"Column Name", "Data Type", "Primary Key", "Unique", "Nullable", "Foreign Key", "Description"}
		for i, heading := range headings {
			cell := fmt.Sprintf("%c1", 66+i) // start from column 2
			excel.SetCellValue(sheet, cell, heading)
			excel.SetCellStyle(sheet, cell, cell, boldStyle)
		}
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
			excel.SetCellStyle(sheet, tableCell, tableCell, boldStyle)
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
	if err := excel.SaveAs(f); err != nil {
		return eris.Wrapf(err, "failed to save to file %s", f)
	}
	return nil
}

func (s *Database) ImportFromExcel(f string, outfile string) error {
	excel, err := excelize.OpenFile(f)
	if err != nil {
		return eris.Wrapf(err, "failed to open the excel file %s", f)
	}
	defer func() {
		if err := excel.Close(); err != nil {
			fmt.Println(eris.ToString(err, true))
		}
	}()

	var data InData

	sheets := excel.GetSheetList()
	for _, sheet := range sheets {
		rows, err := excel.GetRows(sheet)
		if err != nil {
			return eris.Wrapf(err, "faled to get the rows from the sheet %s", sheet)
		}

		schema := &InSchema{Name: sheet}

		var table *InTable
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
				// database columns start from column B
				table.Columns = append(table.Columns, InColumn{Values: row[1:]})
			}
		}
		// last table
		schema.Tables = append(schema.Tables, *table)
		data.Schemas = append(data.Schemas, *schema)
	}

	yaml.FutureLineWrap()
	bytes, err := yaml.Marshal(&data)
	if err != nil {
		return eris.Wrap(err, "failed to transform to yaml")
	}
	output := string(bytes)
	// remove useless key
	output = strings.ReplaceAll(output, "columnValues: ", "")
	ioutil.WriteFile(outfile, []byte(output), 0744)
	return nil
}
