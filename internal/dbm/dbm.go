package dbm

import "github.com/zrs01/dst/model"

type DataService interface {
	DataSourceName() string
	Read(schemaName string) (*model.SchemaDef, error)
}
