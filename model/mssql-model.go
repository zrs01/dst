package model

type MSSQLSchema struct {
	CatalogName                string `db:"CATALOG_NAME"`
	SchemaName                 string `db:"SCHEMA_NAME"`
	SchemaOwner                string `db:"SCHEMA_OWNER"`
	DefaultCharacterSetCatalog string `db:"DEFAULT_CHARACTER_SET_CATALOG"`
	DefaultCharacterSetSchema  string `db:"DEFAULT_CHARACTER_SET_SCHEMA"`
	DefaultCharacterSetName    string `db:"DEFAULT_CHARACTER_SET_NAME"`
}

type MSSQLTable struct {
	TableCatalog string `db:"TABLE_CATALOG"`
	TableSchema  string `db:"TABLE_SCHEMA"`
	TableName    string `db:"TABLE_NAME"`
	TableType    string `db:"TABLE_TYPE"`
}
