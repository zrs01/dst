package yml

import (
	"strings"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

func readFromDb(in string) (*model.DataDef, error) {
	if strings.HasPrefix(in, "sqlserver") {
		return readMSSQL(in)
	}
	return nil, tracerr.New("unknown database type")
}

func readMSSQL(in string) (*model.DataDef, error) {
	return nil, nil
}
