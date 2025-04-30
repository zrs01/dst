package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

type DatabaseManager struct {
	targetSchemaName string
	targetTableName  string
}

// Database represents a database in the database
type Database struct {
	Name         string
	CharacterSet string
	Collate      string
	Comment      string
}

func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{}
}

func (m *DatabaseManager) WithSchemaName(value string) *DatabaseManager {
	m.targetSchemaName = value
	return m
}

func (m *DatabaseManager) WithTableName(value string) *DatabaseManager {
	m.targetTableName = value
	return m
}

func (m *DatabaseManager) GetDef(db *sql.DB) ([]*Database, error) {
	query := `
SELECT
	SCHEMA_NAME,
	DEFAULT_CHARACTER_SET_NAME,
	DEFAULT_COLLATION_NAME,
	SCHEMA_COMMENT
FROM
	INFORMATION_SCHEMA.SCHEMATA`

	args := []any{}
	cond := util.NewCondition()
	if m.targetSchemaName != "" {
		cond.AndWithParam("SCHEMA_NAME = ?", m.targetSchemaName)
	}
	if cond.HasCondition() {
		condStmt, condArgs := cond.Build()
		query = fmt.Sprintf("%s WHERE %s", query, condStmt)
		args = append(args, condArgs...)
	}

	var databases []*Database
	if err := util.Query(db, query, args, func(rows *sql.Rows) error {
		var database Database
		if err := rows.Scan(
			&database.Name,
			&database.CharacterSet,
			&database.Collate,
			&database.Comment); err != nil {
			return tracerr.Wrap(err)
		}
		databases = append(databases, &database)
		return nil
	}); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return databases, nil
}
