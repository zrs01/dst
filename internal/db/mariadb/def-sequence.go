package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

type SequenceManager struct {
	TargetSequenceSchema string
	TargetSequenceName   string
}

// Sequence represents a sequence in the database
type Sequence struct {
	SequenceSchema        string // Database name
	SequenceName          string // Sequence name
	DataType              string // Data type
	NumbericPrecision     int    // Numeric precision
	NumericPrecisionRadix int    // Numeric precision radix
	NumbericScale         int    // Numeric scale
	StartValue            int64  // Start value
	MinimumValue          int64  // Minimum value
	MaximumValue          int64  // Maximum value
	Increment             int64  // Increment
	CycleOption           int    // Cycle option
}

func NewSequenceManager() *SequenceManager {
	return &SequenceManager{}
}

func (m *SequenceManager) WithSequenceSchema(value string) *SequenceManager {
	m.TargetSequenceSchema = value
	return m
}

func (m *SequenceManager) WithSequenceName(value string) *SequenceManager {
	m.TargetSequenceName = value
	return m
}

func (m *SequenceManager) GetDef(db *sql.DB) ([]*Sequence, error) {
	query := `
SELECT
	SEQUNCE_SCHEMA,
	SEQUENCE_NAME,
	DATA_TYPE,
	NUMERIC_PRECISION,
	NUMERIC_PRECISION_RADIX,
	NUMERIC_SCALE,
	START_VALUE,
	MINIMUM_VALUE,
	MAXIMUM_VALUE,
	INCREMENT,
	CYCLE_OPTION
FROM
  INFORMATION_SCHEMA.SEQUENCES`

	args := []any{}
	cond := util.NewCondition()
	if m.TargetSequenceSchema != "" {
		cond.AndWithParam("SEQUENCE_SCHEMA = ?", m.TargetSequenceSchema)
	}
	if m.TargetSequenceName != "" {
		cond.AndWithParam("SEQUENCE_NAME = ?", m.TargetSequenceName)
	}
	if cond.HasCondition() {
		condStmt, condArgs := cond.Build()
		query = fmt.Sprintf("%s WHERE %s", query, condStmt)
		args = append(args, condArgs...)
	}

	var sequences []*Sequence
	if err := util.Query(db, query, args, func(rows *sql.Rows) error {
		var sequence Sequence
		if err := rows.Scan(
			&sequence.SequenceSchema,
			&sequence.SequenceName,
			&sequence.DataType,
			&sequence.NumbericPrecision,
			&sequence.NumericPrecisionRadix,
			&sequence.NumbericScale,
			&sequence.StartValue,
			&sequence.MinimumValue,
			&sequence.MaximumValue,
			&sequence.Increment,
			&sequence.CycleOption); err != nil {
			return tracerr.Wrap(err)
		}
		sequences = append(sequences, &sequence)
		return nil
	}); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return sequences, nil

}
