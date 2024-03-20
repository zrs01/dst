package yml

import (
	"strings"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/zrs01/dst/internal/dbm"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

func readFromDb(in string) (*model.DataDef, error) {
	if strings.HasPrefix(in, "sqlserver") {
		return dbm.ReadMSSQL(in)
	}
	return nil, tracerr.New("unknown database type")
}
