package ddwriter

import (
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zrs01/dst/model"
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

func ReadXlsx(infile string) (*model.Root, error) {
	excel, err := excelize.OpenFile(infile)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var data model.Root

	sheets := excel.GetSheetList()
	for _, sheet := range sheets {
		rows, err := excel.GetRows(sheet)
		if err != nil {
			tracerr.Wrap(err)
		}

		schema := &model.Schema{Name: sheet}

		var table *model.Table
		for rowIndex, row := range rows {
			if rowIndex == 0 {
				// skip the heading row
				continue
			}

			if len(row) > 0 && row[0] != "" {
				// table title (first column contains table name and descrition only)

				if table != nil {
					// append last table instance
					schema.Tables = append(schema.Tables, table)
				}
				table = &model.Table{}

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
				incol := model.Column{}
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
				// table.OutColumns = append(table.OutColumns, model.OutColumn{Value: incol})
			}
		}
		// last table
		schema.Tables = append(schema.Tables, table)
		data.Schemas = append(data.Schemas, schema)
	}
	return &data, err
}
