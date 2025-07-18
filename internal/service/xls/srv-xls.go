package xls

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
	"github.com/zrs01/dst/config"
	"github.com/zrs01/dst/internal/service"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

type StyleType string

const (
	StyleHeader StyleType = "header"
	StyleTable  StyleType = "table"
	StyleKey    StyleType = "keydata"
)

func Generate() error {
	builder := service.NewLoadBuilder(
		service.WithTablePattern(config.Setting.TableName))
	schema, err := builder.LoadFromFile(config.Setting.Sdf)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// filter the data if table name is provided
	if err := schema.Filter(config.Setting.TableName, config.Setting.ColumnName); err != nil {
		return tracerr.Wrap(err)
	}
	// append common columns
	if err := service.AppendCommonColumns(schema); err != nil {
		return tracerr.Wrap(err)
	}

	if !strings.HasSuffix(config.Setting.Output, ".xlsx") {
		config.Setting.Output = config.Setting.Output + ".xlsx"
	}
	if err := WriteXlsx(schema, config.Setting.Output); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func WriteXlsx(schema *model.Schema, out string) error {
	excel := excelize.NewFile()
	if err := createColumnSheet(excel, schema); err != nil {
		return tracerr.Wrap(err)
	}
	if err := createTableSheet(excel, schema); err != nil {
		return tracerr.Wrap(err)
	}
	excel.DeleteSheet("Sheet1")
	if err := excel.SaveAs(out); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func createColumnSheet(excel *excelize.File, schema *model.Schema) error {
	sheet := fmt.Sprintf("%s (columns)", schema.Name)
	excel.NewSheet(sheet)

	style, err := definedExcelStyle(excel)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// header
	header := []string{"#", "Table Name", "Column Name", "Data Type", "Identity", "Not Null", "Default", "Foreign Key", "Description"}
	for i, item := range header {
		cell := fmt.Sprintf("%c1", 65+i)
		excel.SetCellValue(sheet, cell, item)
	}
	// header style
	excel.SetCellStyle(sheet, "A1", fmt.Sprintf("%c1", 65+len(header)-1), (*style)[StyleHeader])
	// column width
	widths := []float64{4, 26, 20, 15, 8, 8, 10, 25, 50}
	for i, width := range widths {
		col := fmt.Sprintf("%c", 65+i)
		excel.SetColWidth(sheet, col, col, width)
	}

	// content starts under the header
	rowCounter := 2

	for _, table := range schema.Tables {
		// table columns
		{
			for i, column := range table.Columns {
				index := i + rowCounter
				excel.SetCellValue(sheet, fmt.Sprintf("A%d", index), i+1)
				excel.SetCellValue(sheet, fmt.Sprintf("B%d", index), table.Name)
				excel.SetCellValue(sheet, fmt.Sprintf("C%d", index), column.Name)
				excel.SetCellValue(sheet, fmt.Sprintf("D%d", index), column.DataType)
				excel.SetCellValue(sheet, fmt.Sprintf("E%d", index), column.Identity)
				excel.SetCellValue(sheet, fmt.Sprintf("F%d", index), column.NotNull)
				excel.SetCellValue(sheet, fmt.Sprintf("G%d", index), column.Value)
				excel.SetCellValue(sheet, fmt.Sprintf("H%d", index), column.ForeignKey)
				excel.SetCellValue(sheet, fmt.Sprintf("I%d", index), column.Desc)
			}
			rowCounter += len(table.Columns)
		}
	}
	return nil
}

func createTableSheet(excel *excelize.File, schema *model.Schema) error {
	sheet := fmt.Sprintf("%s (tables)", schema.Name)
	excel.NewSheet(sheet)

	style, err := definedExcelStyle(excel)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// header
	header := []string{"#", "Title Name", "Table Description", "Key Data Item"}
	for i, item := range header {
		cell := fmt.Sprintf("%c1", 65+i)
		excel.SetCellValue(sheet, cell, item)
	}
	// header style
	excel.SetCellStyle(sheet, "A1", fmt.Sprintf("%c1", 65+len(header)-1), (*style)[StyleHeader])
	// column width
	widths := []float64{4, 26, 60, 50}
	for i, width := range widths {
		col := fmt.Sprintf("%c", 65+i)
		excel.SetColWidth(sheet, col, col, width)
	}

	rowCounter := 2

	for i, table := range schema.Tables {
		excel.SetCellValue(sheet, fmt.Sprintf("A%d", rowCounter), i+1)
		excel.SetCellValue(sheet, fmt.Sprintf("B%d", rowCounter), table.Name)
		excel.SetCellValue(sheet, fmt.Sprintf("C%d", rowCounter), table.Desc)

		// find PK and all FK of the table
		keyData := ""
		keyDataRow := 0
		for k, column := range table.Columns {
			if strings.ToUpper(column.Identity) == "Y" || lo.IsNotEmpty(column.ForeignKey) {
				if k != 0 {
					keyData += "\n"
				}
				if strings.ToUpper(column.Identity) == "Y" {
					keyData += fmt.Sprintf("%s (%s)", column.Name, "PK")
					keyDataRow += 1
				} else if lo.IsNotEmpty(column.ForeignKey) {
					keyData += fmt.Sprintf("%s (%s)", column.Name, "FK")
					keyDataRow += 1
				}
			}
		}
		excel.SetCellValue(sheet, fmt.Sprintf("D%d", rowCounter), keyData)
		excel.SetCellStyle(sheet, fmt.Sprintf("D%d", rowCounter), fmt.Sprintf("D%d", rowCounter), (*style)[StyleKey])

		h, _ := excel.GetRowHeight(sheet, 1)
		excel.SetRowHeight(sheet, rowCounter, h*float64(keyDataRow))
		rowCounter += 1
	}
	return nil
}

func definedExcelStyle(excel *excelize.File) (*map[StyleType]int, error) {
	style := make(map[StyleType]int, 0)

	/* ------------------------ style for the header cell ----------------------- */
	header, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#b4c7dc"}, Pattern: 1},
	})
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	style[StyleHeader] = header

	/* ---------------------- style for the table name cell --------------------- */
	table, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#dee6ef"}, Pattern: 1},
	})
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	style[StyleTable] = table

	/* ----------------------- style for the key data cell ---------------------- */
	keyStyle, err := excel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: true},
	})
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	style[StyleKey] = keyStyle

	return &style, nil
}
