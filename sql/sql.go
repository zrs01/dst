package sql

import (
	"embed"
	"fmt"

	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/transform"
	"github.com/ztrue/tracerr"
)

//go:embed template/*.jet
var fs embed.FS

func CreateTable(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("template/%s-create-table.jet", db), out)
}

func DropTable(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("template/%s-drop-table.jet", db), out)
}

func AddColumn(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("template/%s-add-column.jet", db), out)
}

func DropColumn(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("template/%s-drop-column.jet", db), out)
}

func RenameColumn(data *model.DataDef, db string, out string, newcol string) error {
	data.CustomData.NewColumnName = newcol
	return writeDDL(data, fmt.Sprintf("template/%s-rename-column.jet", db), out)
}

func writeDDL(data *model.DataDef, template string, out string) error {
	b, err := fs.ReadFile(template)
	if err != nil {
		return tracerr.Wrap(err)
	}
	if err := transform.WriteMemoryTpl(data, template, string(b), out); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}
