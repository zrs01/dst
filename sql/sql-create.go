package sql

import (
	"embed"
	"fmt"

	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

//go:embed template/*.jet
var fs embed.FS

func WriteCreateTable(data *model.DataDef, out string) error {
	b, err := fs.ReadFile("template/create.jet")
	if err != nil {
		return tracerr.Wrap(err)
	}
	fmt.Println(string(b))
	return nil
}
