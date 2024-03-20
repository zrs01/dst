package dbm

type MetaData interface {
	Tables(schemaName string) any
	Columns(tableName string) any
}
