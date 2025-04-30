package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

type IndexManager struct {
	TargetTableSchema string
	TargetTableName   string
}

// Index represents an index in a table
type Index struct {
	TableSchema string
	TableName   string
	IndexName   string
	Columns     []string
	Unique      bool
}

func NewIndexManager() *IndexManager {
	return &IndexManager{}
}

func (m *IndexManager) WithTableSchema(value string) *IndexManager {
	m.TargetTableSchema = value
	return m
}

func (m *IndexManager) WithTableName(value string) *IndexManager {
	m.TargetTableName = value
	return m
}

func (m *IndexManager) GetDef(db *sql.DB) ([]*Index, error) {
	query := `
SELECT
	TABLE_SCHEMA,
	TABLE_NAME,
	INDEX_NAME,
	NON_UNIQUE,
	COLUMN_NAME
FROM
	INFORMATION_SCHEMA.STATISTICS`

	args := []any{}
	cond := util.NewCondition()
	if m.TargetTableSchema != "" {
		cond = cond.AndWithParam("TABLE_SCHEMA = ?", m.TargetTableSchema)
	}
	if m.TargetTableName != "" {
		cond = cond.AndWithParam("TABLE_NAME = ?", m.TargetTableName)
	}
	if cond.HasCondition() {
		condStmt, condArgs := cond.Build()
		query = fmt.Sprintf("%s WHERE %s", query, condStmt)
		args = append(args, condArgs...)
	}
	query = fmt.Sprintf("%s ORDER BY INDEX_SCHEMA, INDEX_NAME, SEQ_IN_INDEX", query)

	indexes := []*Index{}
	indexMap := make(map[string]*Index)
	if err := util.Query(db, query, args, func(rows *sql.Rows) error {
		var tableSchema string
		var tableName string
		var indexName string
		var columnName string
		var nonUnique int
		if err := rows.Scan(
			&tableSchema,
			&tableName,
			&indexName,
			&nonUnique,
			&columnName); err != nil {
			return tracerr.Wrap(err)
		}
		unique := nonUnique == 0
		if idx, ok := indexMap[indexName]; ok {
			idx.Columns = append(idx.Columns, columnName)
		} else {
			indexMap[indexName] = &Index{
				TableSchema: tableSchema,
				TableName:   tableName,
				IndexName:   indexName,
				Columns:     []string{columnName},
				Unique:      unique,
			}
		}
		return nil
	}); err != nil {
		return nil, tracerr.Wrap(err)
	}
	for _, idx := range indexMap {
		indexes = append(indexes, idx)
	}
	return indexes, nil
}
