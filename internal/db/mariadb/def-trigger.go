package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

type TriggerManager struct {
	TargetTriggerSchema string
	TargetTriggerName   string
}

// Trigger represents a trigger in the database
type Trigger struct {
	TriggerSchema     string // Database name in which the trigger occurs
	TriggerName       string // Name of the trigger
	EventManipulation string // INSERT, UPDATE, DELETE
	EventObjectSchema string // Database name on which the trigger acts
	EventObjectTable  string // Table name on which the trigger acts
	ActionOrder       int    // Order of the trigger
	ActionStatement   string // SQL statement to execute when the trigger is fired
	ActionTiming      string // BEFORE or AFTER
	Definer           string // User account that created the trigger
}

func NewTriggerManager() *TriggerManager {
	return &TriggerManager{}
}

func (m *TriggerManager) WithTriggerSchema(schema string) *TriggerManager {
	m.TargetTriggerSchema = schema
	return m
}

func (m *TriggerManager) WithTriggerName(name string) *TriggerManager {
	m.TargetTriggerName = name
	return m
}

func (m *TriggerManager) GetDef(db *sql.DB) ([]*Trigger, error) {
	query := `
SELECT
	TRIGGER_SCHEMA,
	TRIGGER_NAME,
	EVENT_MANIPULATION,
	EVENT_OBJECT_SCHEMA,
	EVENT_OBJECT_TABLE,
	ACTION_ORDER,
	ACTION_STATEMENT,
	ACTION_TIMING,
	DEFINER
FROM
  INFORMATION_SCHEMA.TRIGGERS`

	args := []any{}
	cond := util.NewCondition()
	if m.TargetTriggerSchema != "" {
		cond.AndWithParam("TRIGGER_SCHEMA = ?", m.TargetTriggerSchema)
	}
	if m.TargetTriggerName != "" {
		cond.AndWithParam("TRIGGER_NAME = ?", m.TargetTriggerName)
	}
	if cond.HasCondition() {
		condStmt, condArgs := cond.Build()
		query = fmt.Sprintf("%s WHERE %s", query, condStmt)
		args = append(args, condArgs...)
	}

	var triggers []*Trigger
	if err := util.Query(db, query, args, func(rows *sql.Rows) error {
		var trigger Trigger
		if err := rows.Scan(
			&trigger.TriggerSchema,
			&trigger.TriggerName,
			&trigger.EventManipulation,
			&trigger.EventObjectSchema,
			&trigger.EventObjectTable,
			&trigger.ActionOrder,
			&trigger.ActionStatement,
			&trigger.Definer); err != nil {
			return tracerr.Wrap(err)
		}
		triggers = append(triggers, &trigger)
		return nil
	}); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return triggers, nil
}
