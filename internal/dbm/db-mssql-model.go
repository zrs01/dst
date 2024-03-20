package dbm

type MsSqlSchema struct {
	CatalogName                string  `db:"CATALOG_NAME"`
	SchemaName                 string  `db:"SCHEMA_NAME"`
	SchemaOwner                string  `db:"SCHEMA_OWNER"`
	DefaultCharacterSetCatalog *string `db:"DEFAULT_CHARACTER_SET_CATALOG"`
	DefaultCharacterSetSchema  *string `db:"DEFAULT_CHARACTER_SET_SCHEMA"`
	DefaultCharacterSetName    *string `db:"DEFAULT_CHARACTER_SET_NAME"`
}

type MsSqlTable struct {
	TableCatalog string `db:"TABLE_CATALOG"`
	TableSchema  string `db:"TABLE_SCHEMA"`
	TableName    string `db:"TABLE_NAME"`
	TableType    string `db:"TABLE_TYPE"`
}

type MsSqlColumn struct {
	TableCatalog           string  `db:"TABLE_CATALOG"`
	TableSchema            string  `db:"TABLE_SCHEMA"`
	TableName              string  `db:"TABLE_NAME"`
	ColumnName             string  `db:"COLUMN_NAME"`
	OrdinalPosition        int     `db:"ORDINAL_POSITION"`
	ColumnDefault          *string `db:"COLUMN_DEFAULT"`
	IsNullable             string  `db:"IS_NULLABLE"`
	DataType               string  `db:"DATA_TYPE"`
	CharacterMaximumLength *int    `db:"CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   *int    `db:"CHARACTER_OCTET_LENGTH"`
	NumericPrecision       *int    `db:"NUMERIC_PRECISION"`
	NumericPrecisionRadix  *int    `db:"NUMERIC_PRECISION_RADIX"`
	NumericScale           *int    `db:"NUMERIC_SCALE"`
	DateTimePrecision      *int    `db:"DATETIME_PRECISION"`
	CharacterSetCatalog    *string `db:"CHARACTER_SET_CATALOG"`
	CharacterSetSchema     *string `db:"CHARACTER_SET_SCHEMA"`
	CharacterSetName       *string `db:"CHARACTER_SET_NAME"`
	CollationCatalog       *string `db:"COLLATION_CATALOG"`
	CollationSchema        *string `db:"COLLATION_SCHEMA"`
	CollationName          *string `db:"COLLATION_NAME"`
	DomainCatalog          *string `db:"DOMAIN_CATALOG"`
	DomainSchema           *string `db:"DOMAIN_SCHEMA"`
	DomainName             *string `db:"DOMAIN_NAME"`
}
