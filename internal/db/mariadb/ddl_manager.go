package mariadb

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

type DDLManager struct {
	targetTableName string
}

func NewDDLManager() *DDLManager {
	return &DDLManager{}
}

func (m *DDLManager) WithTableName(tableName string) *DDLManager {
	m.targetTableName = tableName
	return m
}

func (m *DDLManager) GenerateStatement(db *model.Schema, stmtType string) (*string, error) {
	switch stmtType {
	case "auto", "create":
		return m.createStmt(db)
	case "alter":
		return m.alterStmt(db)
	default:
		return nil, tracerr.Wrap(fmt.Errorf("unsupported statement type: %s", stmtType))
	}
}

func (m *DDLManager) createStmt(mdb *model.Schema) (*string, error) {
	mTables := mdb.Tables
	if m.targetTableName != "" {
		mTables = lo.Filter(mTables, func(mTable *model.Table, _ int) bool {
			return mTable.Name == m.targetTableName
		})
	}

	var stmts []string
	for _, mTable := range mTables {
		stmt, err := m.createTableStmt(mTable)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		stmts = append(stmts, stmt...)
	}

	output := strings.Join(stmts, "\n")
	return &output, nil
}

func (m *DDLManager) createTableStmt(mTable *model.Table) ([]string, error) {
	var stmts, colStmts []string
	for _, mColumn := range mTable.Columns {
		col := fmt.Sprintf("\t`%s` %s", mColumn.Name, mColumn.DataType)
		col += m.buildColumnAttribues(mColumn)
		colStmts = append(colStmts, col)
	}

	// primary index
	for _, column := range mTable.Columns {
		if util.IsYes(column.Identity) {
			colStmts = append(colStmts, fmt.Sprintf("\tPRIMARY KEY (`%s`)", column.Name))
		}
	}

	joinedItems := strings.Join(colStmts, ",\n")
	if mTable.Desc == "" {
		mTable.Desc = mTable.Name
	}
	stmts = append(stmts, fmt.Sprintf("CREATE TABLE `%s` IF NOT EXISTS (\n%s\n) COMMENT `%s`;\n", mTable.Name, joinedItems, mTable.Desc))

	// column index
	for _, column := range mTable.Columns {
		index := m.buildColumnIndex(mTable.Name, column)
		if index != "" {
			stmts = append(stmts, index)
		}
	}

	// additional index
	for i, index := range mTable.Indexes {
		keyName := index.Name
		if keyName == "" {
			keyName = fmt.Sprintf("idx_%s_%s", mTable.Name, fmt.Sprintf("%03d", (i+1)))
		}
		stmts = append(stmts, m.createIndex(keyName, mTable.Name, index.Columns, util.IsYes(index.Unique)))
	}

	// foreign key
	fkColumns := lo.Filter(mTable.Columns, func(column *model.Column, index int) bool {
		return column.ForeignKey != ""
	})
	for _, fkc := range fkColumns {
		ref := strings.Split(fkc.ForeignKey, ".")
		stmts = append(stmts, fmt.Sprintf("ALTER TABLE `%s` ADD CONSTRAINT `fk_%s_%s` FOREIGN KEY (`%s`) REFERENCES `%s` (`%s`);", mTable.Name, mTable.Name, fkc.Name, fkc.Name, ref[0], ref[1]))
	}

	return stmts, nil
}

func (m *DDLManager) alterStmt(mdb *model.Schema) (*string, error) {
	mTables := mdb.Tables
	if m.targetTableName != "" {
		mTables = lo.Filter(mTables, func(mTable *model.Table, _ int) bool {
			return mTable.Name == m.targetTableName
		})
	}

	var stmts []string
	for _, mTable := range mTables {
		stmt, err := m.alterTableStmt(mdb.Name, mTable)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		stmts = append(stmts, stmt...)
	}

	output := strings.Join(stmts, "\n")
	return &output, nil
}

func (m *DDLManager) alterTableStmt(srcDBName string, srcTB *model.Table) ([]string, error) {
	// map current database schema to model
	expDB, err := NewExportManager().WithDatabaseName(srcDBName).ToModel()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	stmts := []string{}

	expTB, ok := lo.Find(expDB.Tables, func(t *model.Table) bool {
		return t.Name == srcTB.Name
	})
	if !ok {
		/* ------------------------------ create table ------------------------------ */
		stmt, err := m.createTableStmt(srcTB)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		stmts = append(stmts, stmt...)
	} else {
		/* ------------------------------- alter table ------------------------------ */
		// column ID support function
		ciRe := regexp.MustCompile(`\[(\d+)\]`)
		findCI := func(v string) string {
			match := ciRe.FindStringSubmatch(v)
			if len(match) > 1 { // column ID found
				return match[1]
			}
			return ""
		}

		eofStmts := []string{}
		for _, srcCol := range srcTB.Columns {
			srcCI := findCI(srcCol.Desc)
			if srcCI == "" { // column ID does not found
				continue
			}
			expCol, ok := lo.Find(expTB.Columns, func(c *model.Column) bool {
				return findCI(c.Desc) == srcCI
			})
			if !ok { // column ID not found
				continue
			}

			if srcCol.DataType != expCol.DataType || srcCol.NotNull != expCol.NotNull || srcCol.Value != expCol.Value || srcCol.Desc != expCol.Desc || srcCol.Unique != expCol.Unique {
				stmts = append(stmts, fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN `%s` %s%s;", srcTB.Name, expCol.Name, srcCol.DataType, m.buildColumnAttribues(srcCol)))
			}
			if srcCol.Name != expCol.Name {
				stmts = append(stmts, fmt.Sprintf("ALTER TABLE `%s` RENAME COLUMN `%s` TO `%s`;", srcTB.Name, expCol.Name, srcCol.Name))
			}

			if srcCol.Index != expCol.Index || srcCol.Unique != expCol.Unique {
				eofStmts = append(eofStmts, fmt.Sprintf("DROP INDEX IF EXISTS `%s` ON `%s`;", m.formatIndexName(srcTB.Name, expCol.Name), srcTB.Name))
				eofStmts = append(eofStmts, m.buildColumnIndex(srcTB.Name, srcCol))
			}
		}
		stmts = append(stmts, eofStmts...)
	}
	return stmts, nil
}

func (m *DDLManager) buildColumnAttribues(col *model.Column) string {
	attr := ""
	if col.Compute != "" {
		attr += fmt.Sprintf(" GENERATED ALWAYS AS (%s)", col.Compute)
	}
	if util.IsYes(col.NotNull) {
		attr += " NOT NULL"
	}
	if col.Value != "" {
		attr += fmt.Sprintf(" DEFAULT %s", col.Value)
	}
	// if util.IsYes(col.AutoIncrement) {
	// 	attr += " AUTO_INCREMENT"
	// }
	if col.Desc != "" {
		attr += fmt.Sprintf(" COMMENT `%s`", col.Desc)
	}
	return attr
}

func (m *DDLManager) buildColumnIndex(tbn string, col *model.Column) string {
	if util.IsYes(col.Index) {
		name := m.formatIndexName(tbn, col.Name)
		return m.createIndex(name, tbn, []string{col.Name}, util.IsYes(col.Unique))
	}
	return ""
}

func (m *DDLManager) createIndex(name, tbn string, fields []string, isUnique bool) string {
	unique := lo.If(isUnique, "UNIQUE ").Else("")
	return fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS `%s` ON `%s` (`%s`);", unique, name, tbn, strings.Join(fields, "`,`"))
}

func (m *DDLManager) formatIndexName(tb, col string) string {
	return fmt.Sprintf("idx_%s_%s", tb, col)
}
