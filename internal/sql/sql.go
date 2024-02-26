package sql

import (
	"embed"
	"fmt"

	"github.com/zrs01/dst/internal/tpl"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

//go:embed templates/*/*.jet
var fs embed.FS

func CreateTable(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-create-table.jet", db, db), out)
}

func DropTable(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-drop-table.jet", db, db), out)
}

func AddColumn(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-add-column.jet", db, db), out)
}

func DropColumn(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-drop-column.jet", db, db), out)
}

func RenameColumn(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-rename-column.jet", db, db), out)
}

func ModifyColumn(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-modify-column.jet", db, db), out)
}

func CreateIndex(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-create-index.jet", db, db), out)
}

func DropIndex(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-drop-index.jet", db, db), out)
}

func writeDDL(data *model.DataDef, template string, out string) error {
	b, err := fs.ReadFile(template)
	if err != nil {
		return tracerr.Wrap(err)
	}
	if err := tpl.WriteMemoryTpl(data, template, string(b), out); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}
