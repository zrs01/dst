package sql

import (
	"embed"
	"fmt"

	"github.com/zrs01/dst/internal/service/txt"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

//go:embed templates/*/*.jet
var fs embed.FS

func CreateTable(schema *model.Schema, db string, out string) error {
	return writeDDL(schema, fmt.Sprintf("templates/%s/%s-create-table.jet", db, db), out)
}

func DropTable(schema *model.Schema, db string, out string) error {
	return writeDDL(schema, fmt.Sprintf("templates/%s/%s-drop-table.jet", db, db), out)
}

func AddColumn(schema *model.Schema, db string, out string) error {
	return writeDDL(schema, fmt.Sprintf("templates/%s/%s-add-column.jet", db, db), out)
}

func DropColumn(schema *model.Schema, db string, out string) error {
	return writeDDL(schema, fmt.Sprintf("templates/%s/%s-drop-column.jet", db, db), out)
}

func RenameColumn(schema *model.Schema, db string, out string) error {
	return writeDDL(schema, fmt.Sprintf("templates/%s/%s-rename-column.jet", db, db), out)
}

func ModifyColumn(schema *model.Schema, db string, out string) error {
	return writeDDL(schema, fmt.Sprintf("templates/%s/%s-modify-column.jet", db, db), out)
}

func CreateIndex(schema *model.Schema, db string, out string) error {
	return writeDDL(schema, fmt.Sprintf("templates/%s/%s-create-index.jet", db, db), out)
}

func DropIndex(schema *model.Schema, db string, out string) error {
	return writeDDL(schema, fmt.Sprintf("templates/%s/%s-drop-index.jet", db, db), out)
}

func writeDDL(schema *model.Schema, template string, out string) error {
	if err := txt.WriteWithEmbedFSLoader(fs, schema, template, out); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}
