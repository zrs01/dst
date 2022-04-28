package main

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/rotisserie/eris"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v2"
)

type Database struct {
	Schemas []struct {
		Name   string `yaml:"name" default:"Schema"`
		Desc   string `yaml:"description"`
		Tables []struct {
			Name    string          `yaml:"name"`
			Desc    string          `yaml:"description"`
			Columns [][]interface{} `yaml:"columns"`
		} `yaml:"tables"`
	} `yaml:"schemas"`
}

type Column struct {
	Name        string
	DataType    string
	IsPK        string
	IsUnique    string
	Nullable    string
	ForeignHint string
	Desc        string
}

func NewDatabase() *Database {
	return &Database{}
}

func (s *Database) Load(f string) error {
	yamlFile, err := ioutil.ReadFile(f)
	if err != nil {
		return eris.Wrapf(err, "failed to read the file %s", f)
	}
	if err := yaml.Unmarshal(yamlFile, s); err != nil {
		return eris.Wrapf(err, "failed to unmarshal the file %s", f)
	}
	return nil
}

func (s *Database) convertToColumn(icols [][]interface{}) []Column {
	var trueFalse = func(value string) string {
		if value == "true" {
			return "Y"
		}
		return "N"
	}
	var createColumn = func(value []interface{}) Column {
		return Column{
			Name:        fmt.Sprintf("%v", value[0]),
			DataType:    fmt.Sprintf("%v", value[1]),
			IsPK:        trueFalse(fmt.Sprintf("%v", value[2])),
			IsUnique:    trueFalse(fmt.Sprintf("%v", value[3])),
			Nullable:    trueFalse(fmt.Sprintf("%v", value[4])),
			ForeignHint: fmt.Sprintf("%v", value[5]),
			Desc:        fmt.Sprintf("%v", value[6]),
		}
	}

	var columns []Column
	for _, icol := range icols {
		// due to merge array in yaml will cause increase array one level, so need check is array of array
		if reflect.TypeOf(icol[0]).Kind() == reflect.Slice {
			for _, sub := range icol {
				columns = append(columns, createColumn(sub.([]interface{})))
			}
		} else {
			columns = append(columns, createColumn(icol))
		}
	}
	return columns
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

	for _, schema := range s.Schemas {
		sheet := schema.Name
		excel.NewSheet(sheet)

		// heading
		headings := []string{"Column Name", "Data Type", "Primary Key", "Unique", "Nullable", "Foreign Key", "Description"}
		for i := 0; i < len(headings); i++ {
			cell := fmt.Sprintf("%c1", 66+i) // start from column 2
			excel.SetCellValue(sheet, cell, headings[i])
			excel.SetCellStyle(sheet, cell, cell, boldStyle)
		}
		// column width
		widths := []float64{2, 20, 15, 12, 12, 12, 20, 50}
		for i := 0; i < len(widths); i++ {
			col := fmt.Sprintf("%c", 65+i)
			excel.SetColWidth(sheet, col, col, widths[i])
		}

		offset := 3

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
			for i := 0; i < len(columns); i++ {
				index := i + offset
				column := columns[i]
				excel.SetCellValue(sheet, fmt.Sprintf("B%d", index), column.Name)
				excel.SetCellValue(sheet, fmt.Sprintf("C%d", index), column.DataType)
				excel.SetCellValue(sheet, fmt.Sprintf("D%d", index), column.IsPK)
				excel.SetCellValue(sheet, fmt.Sprintf("E%d", index), column.IsUnique)
				excel.SetCellValue(sheet, fmt.Sprintf("F%d", index), column.Nullable)
				excel.SetCellValue(sheet, fmt.Sprintf("G%d", index), column.ForeignHint)
				excel.SetCellValue(sheet, fmt.Sprintf("H%d", index), column.Desc)
			}
			offset = offset + len(columns) + 1
		}
	}
	excel.DeleteSheet("Sheet1")
	if err := excel.SaveAs(f); err != nil {
		return eris.Wrapf(err, "failed to save to file %s", f)
	}
	return nil
}
