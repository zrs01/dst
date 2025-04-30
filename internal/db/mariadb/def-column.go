package mariadb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

type ColumnManager struct {
	targetTableSchema string
	targetTableName   string
	targetColumnName  string
}

// Column represents a column in a table
type Column struct {
	TableSchema          string
	TableName            string
	ColumnName           string
	OrdinalPosition      int
	ColumnDefault        string // Default value
	IsNullable           string
	DataType             string
	Length               int
	Precision            int
	ColumnType           string // e.g. VARCHAR(255)
	ColumnKey            string // PRI|MUL|empty
	Extra                string
	AutoIncrement        bool
	Comment              string
	IsGenerated          string
	GenerationExpression string
}

func NewColumnManager() *ColumnManager {
	return &ColumnManager{}
}

func (m *ColumnManager) WithTableSchema(value string) *ColumnManager {
	m.targetTableSchema = value
	return m
}

func (m *ColumnManager) WithTableName(value string) *ColumnManager {
	m.targetTableName = value
	return m
}

func (m *ColumnManager) WithColumnName(value string) *ColumnManager {
	m.targetColumnName = value
	return m
}

func (m *ColumnManager) GetDef(db *sql.DB) ([]*Column, error) {
	query := `
	SELECT
		TABLE_SCHEMA,
		TABLE_NAME,
		COLUMN_NAME,
		ORDINAL_POSITION,
		IFNULL(COLUMN_DEFAULT, '') AS COLUMN_DEFAULT,
		IS_NULLABLE,
		DATA_TYPE,
		IFNULL(CHARACTER_MAXIMUM_LENGTH, -1) AS CHARACTER_MAXIMUM_LENGTH,
		IFNULL(NUMERIC_PRECISION, -1) AS NUMERIC_PRECISION,
		COLUMN_TYPE,
		COLUMN_KEY,
		EXTRA,
		COLUMN_COMMENT,
		IS_GENERATED,
		IFNULL(GENERATION_EXPRESSION, '') AS GENERATION_EXPRESSION
	FROM
		INFORMATION_SCHEMA.COLUMNS`

	args := []any{}
	cond := util.NewCondition()
	if m.targetTableSchema != "" {
		cond.AndWithParam("TABLE_SCHEMA = ?", m.targetTableSchema)
	}
	if m.targetTableName != "" {
		cond.AndWithParam("TABLE_NAME = ?", m.targetTableName)
	}
	if m.targetColumnName != "" {
		cond.AndWithParam("COLUMN_NAME = ?", m.targetColumnName)
	}
	if cond.HasCondition() {
		condStmt, condArgs := cond.Build()
		query = fmt.Sprintf("%s WHERE %s", query, condStmt)
		args = append(args, condArgs...)
	}
	query = fmt.Sprintf("%s ORDER BY TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION", query)

	columns := []*Column{}
	if err := util.Query(db, query, args, func(rows *sql.Rows) error {
		var column Column
		if err := rows.Scan(
			&column.TableSchema,
			&column.TableName,
			&column.ColumnName,
			&column.OrdinalPosition,
			&column.ColumnDefault,
			&column.IsNullable,
			&column.DataType,
			&column.Length,
			&column.Precision,
			&column.ColumnType,
			&column.ColumnKey,
			&column.Extra,
			&column.Comment,
			&column.IsGenerated,
			&column.GenerationExpression,
		); err != nil {
			return tracerr.Wrap(err)
		}
		column.AutoIncrement = strings.Contains(column.Extra, "auto_increment")
		columns = append(columns, &column)
		return nil
	}); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return columns, nil
}
