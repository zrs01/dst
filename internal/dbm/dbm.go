package dbm

import "github.com/zrs01/dst/model"

type DataSource interface {
	Read(dataSourceName string, schemaName string) (*model.DataDef, error)
}
