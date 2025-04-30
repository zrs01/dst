package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

type TableManager struct {
	TargetTableSchema string
	TargetTableName   string
	TargetTableType   string
}

// Table represents a table in the database
type Table struct {
	TableSchema string
	TableName   string
	TableType   string // BASE TABLE | VIEW | SYSTEM VIEW | VERSIONED | SEQUENCE | TEMPORARY
	Engine      string // Storage engine (e.g., InnoDB, MyISAM)
	Comment     string
}

func NewTableManager() *TableManager {
	return &TableManager{}
}

func (m *TableManager) WithTableSchema(value string) *TableManager {
	m.TargetTableSchema = value
	return m
}

func (m *TableManager) WithTableName(value string) *TableManager {
	m.TargetTableName = value
	return m
}

func (m *TableManager) WithTableType(value string) *TableManager {
	m.TargetTableType = value
	return m
}

func (m *TableManager) GetDef(db *sql.DB) ([]*Table, error) {
	query := `
SELECT
	TABLE_SCHEMA,
	TABLE_NAME,
	TABLE_TYPE,
	IFNULL(ENGINE, '') AS ENGINE,
	IFNULL(TABLE_COMMENT, '') AS TABLE_COMMENT
FROM
	INFORMATION_SCHEMA.TABLES`

	args := []any{}
	cond := util.NewCondition()
	if m.TargetTableSchema != "" {
		cond.AndWithParam("TABLE_SCHEMA = ?", m.TargetTableSchema)
	}
	if m.TargetTableName != "" {
		cond.AndWithParam("TABLE_NAME = ?", m.TargetTableName)
	}
	if m.TargetTableType != "" {
		cond.AndWithParam("TABLE_TYPE = ?", m.TargetTableType)
	}
	if cond.HasCondition() {
		condStmt, condArgs := cond.Build()
		query = fmt.Sprintf("%s WHERE %s", query, condStmt)
		args = append(args, condArgs...)
	}
	query = fmt.Sprintf("%s ORDER BY TABLE_NAME", query)

	tables := []*Table{}
	if err := util.Query(db, query, args, func(rows *sql.Rows) error {
		var table Table
		if err := rows.Scan(
			&table.TableSchema,
			&table.TableName,
			&table.TableType,
			&table.Engine,
			&table.Comment); err != nil {
			return tracerr.Wrap(err)
		}
		tables = append(tables, &table)
		return nil
	}); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return tables, nil
}
