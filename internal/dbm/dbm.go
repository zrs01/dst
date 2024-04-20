package dbm

import "github.com/zrs01/dst/model"

type DatabaseService interface {
	Read(dataSourceName string, schemaName string) (*model.DataDef, error)
}
