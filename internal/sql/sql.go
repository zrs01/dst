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
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-table-create.jet", db, db), out)
}

func DropTable(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-table-drop.jet", db, db), out)
}

func AddColumn(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-column-add.jet", db, db), out)
}

func DropColumn(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-column-drop.jet", db, db), out)
}

func RenameColumn(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-column-rename.jet", db, db), out)
}

func ModifyColumn(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-column-modify.jet", db, db), out)
}

func CreateIndex(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-index-create.jet", db, db), out)
}

func DropIndex(data *model.DataDef, db string, out string) error {
	return writeDDL(data, fmt.Sprintf("templates/%s/%s-index-drop.jet", db, db), out)
}

func writeDDL(data *model.DataDef, template string, out string) error {
	if err := tpl.WriteEmbedFSTpl(fs, data, template, out); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}
