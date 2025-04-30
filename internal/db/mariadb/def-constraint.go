package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

type ConstraintManager struct {
	TargetConstraintName        string
	TargetTableSchema           string
	TargetTableName             string
	TargetReferencedTableSchema string
	TargetReferencedTableName   string
}

// Constraint represents a constraint in a table
type Constraint struct {
	ConstraintName        string
	TableSchema           string
	TableName             string
	ColumnName            string
	ReferencedTableSchema string
	ReferencedTableName   string
	ReferencedColumnName  string
}

func NewConstraintManager() *ConstraintManager {
	return &ConstraintManager{}
}

func (m *ConstraintManager) WithConstraintName(value string) *ConstraintManager {
	m.TargetConstraintName = value
	return m
}

func (m *ConstraintManager) WithTableSchema(value string) *ConstraintManager {
	m.TargetTableSchema = value
	return m
}

func (m *ConstraintManager) WithTableName(value string) *ConstraintManager {
	m.TargetTableName = value
	return m
}

func (m *ConstraintManager) WithReferencedTableSchema(value string) *ConstraintManager {
	m.TargetReferencedTableSchema = value
	return m
}

func (m *ConstraintManager) WithReferencedTableName(value string) *ConstraintManager {
	m.TargetReferencedTableName = value
	return m
}

func (m *ConstraintManager) GetDef(db *sql.DB) ([]*Constraint, error) {
	query := `
SELECT
	CONSTRAINT_NAME,
	TABLE_SCHEMA,
	TABLE_NAME,
	COLUMN_NAME,
	IFNULL(REFERENCED_TABLE_SCHEMA, '') AS REFERENCED_TABLE_SCHEMA,
	IFNULL(REFERENCED_TABLE_NAME, '') AS REFERENCED_TABLE_NAME,
	IFNULL(REFERENCED_COLUMN_NAME, '') AS REFERENCED_COLUMN_NAME
FROM
	INFORMATION_SCHEMA.KEY_COLUMN_USAGE`

	args := []any{}
	cond := util.NewCondition()
	if m.TargetConstraintName != "" {
		cond.AndWithParam("CONSTRAINT_NAME = ?", m.TargetConstraintName)
	}
	if m.TargetTableSchema != "" {
		cond.AndWithParam("TABLE_SCHEMA = ?", m.TargetTableSchema)
	}
	if m.TargetTableName != "" {
		cond.AndWithParam("TABLE_NAME = ?", m.TargetTableName)
	}
	if m.TargetReferencedTableSchema != "" {
		cond.AndWithParam("REFERENCED_TABLE_SCHEMA = ?", m.TargetReferencedTableSchema)
	}
	if m.TargetReferencedTableName != "" {
		cond.AndWithParam("REFERENCED_TABLE_NAME = ?", m.TargetReferencedTableName)
	}
	if cond.HasCondition() {
		condStmt, condArgs := cond.Build()
		query = fmt.Sprintf("%s WHERE %s", query, condStmt)
		args = append(args, condArgs...)
	}

	constraints := []*Constraint{}
	if err := util.Query(db, query, args, func(rows *sql.Rows) error {
		var constraint Constraint
		if err := rows.Scan(
			&constraint.ConstraintName,
			&constraint.TableSchema,
			&constraint.TableName,
			&constraint.ColumnName,
			&constraint.ReferencedTableSchema,
			&constraint.ReferencedTableName,
			&constraint.ReferencedColumnName); err != nil {
			return tracerr.Wrap(err)
		}
		constraints = append(constraints, &constraint)
		return nil
	}); err != nil {
		return nil, err
	}

	return constraints, nil
}
