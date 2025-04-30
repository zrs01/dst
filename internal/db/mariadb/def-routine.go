package mariadb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

type RoutineManager struct {
	TargetSpecificName  string
	TargetRoutineSchema string
}

type Routine struct {
	RoutineName       string
	RoutineSchema     string // Database name associated with routine
	RoutineType       string // Type of routine (FUNCTION or PROCEDURE)
	DataType          string // The return value's data type (for stored functions)
	RoutineDefinition string // The routine's SQL code
	Definer           string // The user who created the routine
	Parameters        []*RoutineParameter
}

type RoutineParameter struct {
	OrdinalPosition int
	ParameterMode   string
	ParameterName   string
	DataType        string
	DTDIentifier    string
}

func NewRoutineManager() *RoutineManager {
	return &RoutineManager{}
}

func (m *RoutineManager) WithSpecificName(value string) *RoutineManager {
	m.TargetSpecificName = value
	return m
}

func (m *RoutineManager) WithRoutineSchema(value string) *RoutineManager {
	m.TargetRoutineSchema = value
	return m
}

func (m *RoutineManager) GetDef(db *sql.DB) ([]*Routine, error) {
	query := `
SELECT
	ROUTINE_NAME,
	ROUTINE_SCHEMA,
	ROUTINE_TYPE,
	DATA_TYPE,
	ROUTINE_DEFINITION,
	DEFINER
FROM
  INFORMATION_SCHEMA.ROUTINES`

	args := []any{}
	cond := util.NewCondition()
	if m.TargetSpecificName != "" {
		cond.AndWithParam("SPECIFIC_NAME = ?", m.TargetSpecificName)
	}
	if m.TargetRoutineSchema != "" {
		cond.AndWithParam("ROUTINE_SCHEMA = ?", m.TargetRoutineSchema)
	}
	if cond.HasCondition() {
		condStmt, condArgs := cond.Build()
		query = fmt.Sprintf("%s WHERE %s", query, condStmt)
		args = append(args, condArgs...)
	}

	var routines []*Routine
	if err := util.Query(db, query, args, func(rows *sql.Rows) error {
		var routine Routine
		if err := rows.Scan(
			&routine.RoutineName,
			&routine.RoutineSchema,
			&routine.RoutineType,
			&routine.DataType,
			&routine.RoutineDefinition,
			&routine.Definer); err != nil {
			return tracerr.Wrap(err)
		}
		routine.RoutineDefinition = strings.ReplaceAll(routine.RoutineDefinition, "\r\n", "\n")
		routines = append(routines, &routine)
		return nil
	}); err != nil {
		return nil, tracerr.Wrap(err)
	}

	paramQuery := `
	SELECT
		ORDINAL_POSITION,
		IFNULL(PARAMETER_MODE, '') AS PARAMETER_MODE,
		IFNULL(PARAMETER_NAME, '') AS PARAMETER_NAME,
		DATA_TYPE,
		DTD_IDENTIFIER
	FROM
		INFORMATION_SCHEMA.PARAMETERS
	WHERE
		SPECIFIC_SCHEMA = ?	AND SPECIFIC_NAME = ?
	ORDER BY ORDINAL_POSITION`

	for _, routine := range routines {
		var params []*RoutineParameter
		if err := util.Query(db, paramQuery, []any{routine.RoutineSchema, routine.RoutineName}, func(rows *sql.Rows) error {
			var param RoutineParameter
			if err := rows.Scan(
				&param.OrdinalPosition,
				&param.ParameterMode,
				&param.ParameterName,
				&param.DataType,
				&param.DTDIentifier); err != nil {
				return tracerr.Wrap(err)
			}
			params = append(params, &param)
			return nil
		}); err != nil {
			return nil, tracerr.Wrap(err)
		}
		routine.Parameters = params
	}
	return routines, nil
}
