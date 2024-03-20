package dbm

import "time"

// select * from information_schema.SCHEMATA;
type MySqlSchema struct {
	CatalogName             string `db:"CATALOG_NAME"`
	SchemaName              string `db:"SCHEMA_NAME"`
	DefaultCharacterSetName string `db:"DEFAULT_CHARACTER_SET_NAME"`
	DefaultCollationName    string `db:"DEFAULT_COLLATION_NAME"`
	SQLPath                 string `db:"SQL_PATH"`
	SchemaComment           string `db:"SCHEMA_COMMENT"`
}

// select * from INFORMATION_SCHEMA.TABLES;
type MySqlTable struct {
	TableCatalog   string    `db:"TABLE_CATALOG"`
	TableSchema    string    `db:"TABLE_SCHEMA"`
	TableName      string    `db:"TABLE_NAME"`
	TableType      string    `db:"TABLE_TYPE"`
	Engine         string    `db:"ENGINE"`
	Version        uint64    `db:"VERSION"`
	RowFormat      string    `db:"ROW_FORMAT"`
	TableRows      uint64    `db:"TABLE_ROWS"`
	AvgRowLength   uint64    `db:"AVG_ROW_LENGTH"`
	DataLength     uint64    `db:"DATA_LENGTH"`
	MaxDataLength  uint64    `db:"MAX_DATA_LENGTH"`
	IndexLength    uint64    `db:"INDEX_LENGTH"`
	DataFree       uint64    `db:"DATA_FREE"`
	AutoIncrement  uint64    `db:"AUTO_INCREMENT"`
	CreateTime     time.Time `db:"CREATE_TIME"`
	UpdateTime     time.Time `db:"UPDATE_TIME"`
	CheckTime      time.Time `db:"CHECK_TIME"`
	TableCollation string    `db:"TABLE_COLLATION"`
	Checksum       uint64    `db:"CHECKSUM"`
	CreateOptions  string    `db:"CREATE_OPTIONS"`
	TableComment   string    `db:"TABLE_COMMENT"`
	MaxIndexLength uint64    `db:"MAX_INDEX_LENGTH"`
	Temporary      string    `db:"TEMPORARY"`
}

// select * from INFORMATION_SCHEMA.columns;
type MySqlColumn struct {
	TableCatalog           string `db:"TABLE_CATALOG"`
	TableSchema            string `db:"TABLE_SCHEMA"`
	TableName              string `db:"TABLE_NAME"`
	ColumnName             string `db:"COLUMN_NAME"`
	OrdinalPosition        uint64 `db:"ORDINAL_POSITION"`
	ColumnDefault          string `db:"COLUMN_DEFAULT"`
	IsNullable             string `db:"IS_NULLABLE"`
	DataType               string `db:"DATA_TYPE"`
	CharacterMaximumLength uint64 `db:"CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   uint64 `db:"CHARACTER_OCTET_LENGTH"`
	NumericPrecision       uint64 `db:"NUMERIC_PRECISION"`
	NumericScale           uint64 `db:"NUMERIC_SCALE"`
	DatetimePrecision      uint64 `db:"DATETIME_PRECISION"`
	CharacterSetName       string `db:"CHARACTER_SET_NAME"`
	CollationName          string `db:"COLLATION_NAME"`
	ColumnType             string `db:"COLUMN_TYPE"`
	ColumnKey              string `db:"COLUMN_KEY"`
	Extra                  string `db:"EXTRA"`
	Privileges             string `db:"PRIVILEGES"`
	ColumnComment          string `db:"COLUMN_COMMENT"`
	IsGenerated            string `db:"IS_GENERATED"`
	GenerationExpression   string `db:"GENERATION_EXPRESSION"`
}

// select * from information_schema.key_column_usage;
type MySqlKeyColumnUsage struct {
	ConstraintCatalog          string `db:"CONSTRAINT_CATALOG"`
	ConstraintSchema           string `db:"CONSTRAINT_SCHEMA"`
	ConstraintName             string `db:"CONSTRAINT_NAME"`
	TableCatalog               string `db:"TABLE_CATALOG"`
	TableSchema                string `db:"TABLE_SCHEMA"`
	TableName                  string `db:"TABLE_NAME"`
	ColumnName                 string `db:"COLUMN_NAME"`
	OrdinalPosition            int64  `db:"ORDINAL_POSITION"`
	PositionInUniqueConstraint int64  `db:"POSITION_IN_UNIQUE_CONSTRAINT"`
	ReferencedTableSchema      string `db:"REFERENCED_TABLE_SCHEMA"`
	ReferencedTableName        string `db:"REFERENCED_TABLE_NAME"`
	ReferencedColumnName       string `db:"REFERENCED_COLUMN_NAME"`
}
