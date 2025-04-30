package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

type ViewManager struct {
	TargetTableSchema string
	TargetTableName   string
}

// View represents a view in the database
type View struct {
	TableSchema         string
	TableName           string
	ViewDefinition      string
	CheckOption         string
	IsUpdatable         string
	Definer             string
	SecurityType        string
	CharacterSetClient  string
	CollationConnection string
	Algorithm           string
}

func NewViewManager() *ViewManager {
	return &ViewManager{}
}

func (m *ViewManager) WithTableSchema(value string) *ViewManager {
	m.TargetTableSchema = value
	return m
}

func (m *ViewManager) WithTableName(value string) *ViewManager {
	m.TargetTableName = value
	return m
}

func (m *ViewManager) GetDef(db *sql.DB) ([]*View, error) {
	query := `
SELECT
	TABLE_SCHEMA,
	TABLE_NAME,
	VIEW_DEFINITION,
	CHECK_OPTION,
	IS_UPDATABLE,
	DEFINER,
	SECURITY_TYPE,
	CHARACTER_SET_CLIENT,
	COLLATION_CONNECTION,
	ALGORITHM
FROM
	INFORMATION_SCHEMA.VIEWS`

	args := []any{}
	cond := util.NewCondition()
	if m.TargetTableSchema != "" {
		cond.AndWithParam("TABLE_SCHEMA = ?", m.TargetTableSchema)
	}
	if m.TargetTableName != "" {
		cond.AndWithParam("TABLE_NAME = ?", m.TargetTableName)
	}
	if cond.HasCondition() {
		condStmt, condArgs := cond.Build()
		query = fmt.Sprintf("%s WHERE %s", query, condStmt)
		args = append(args, condArgs...)
	}
	query = fmt.Sprintf("%s ORDER BY TABLE_NAME", query)

	views := []*View{}
	if err := util.Query(db, query, args, func(rows *sql.Rows) error {
		var view View
		if err := rows.Scan(
			&view.TableSchema,
			&view.TableName,
			&view.ViewDefinition,
			&view.CheckOption,
			&view.IsUpdatable,
			&view.Definer,
			&view.SecurityType,
			&view.CharacterSetClient,
			&view.CollationConnection,
			&view.Algorithm); err != nil {
			return tracerr.Wrap(err)
		}

		views = append(views, &view)
		return nil
	}); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return views, nil
}
