package mariadb

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/zrs01/dst/internal/db/dbcm"
	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/util"
)

// DDLBuilder is a builder for generating DDL statements
type DDLBuilder struct{}

func NewDDBuilder() dbcm.DDL {
	return &DDLBuilder{}
}

// CreateDatabase generates the DDL statement for creating the database
func (m *DDLBuilder) CreateDatabase(schemaModel *model.Schema) (string, error) {
	var stmt strings.Builder
	var args []any
	stmt.WriteString("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
	args = append(args, schemaModel.Name)
	if schemaModel.Desc != "" {
		stmt.WriteString(" COMMENT '%s'")
		args = append(args, schemaModel.Desc)
	}
	stmt.WriteString(";")
	return fmt.Sprintf(stmt.String(), args...), nil
}

func (m *DDLBuilder) CreateIndex(indexName, tableName string, fields []string, isUnique bool) string {
	unique := lo.If(isUnique, "UNIQUE ").Else("")
	return fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS `%s` ON `%s` (`%s`);", unique, indexName, tableName, strings.Join(fields, "`,`"))
}

// CreateTable implements db.DDL.
func (m *DDLBuilder) CreateTableWithDirect(table *model.Table) string {
	return m.buildCreateTableStmt(table.Name, table.Desc, table.Columns)
}

func (m *DDLBuilder) CreateTableWithAlter(table *model.Table) string {
	var stmts []string
	columns := lo.Filter(table.Columns, func(column *model.Column, index int) bool {
		return util.IsYes(column.Identity)
	})
	stmts = []string{m.buildCreateTableStmt(table.Name, table.Desc, columns)}

	for i, mColumn := range table.Columns {
		if !util.IsYes(mColumn.Identity) {
			stmts = append(stmts, m.AddColumn(table.Name, mColumn, table.Columns[i-1].Name))
		}
	}
	stmts = append(stmts, "")
	return strings.Join(stmts, "\n")
}

func (m *DDLBuilder) AddColumn(tableName string, column *model.Column, previousColumnName string) string {
	col := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN IF NOT EXISTS `%s` %s", tableName, column.Name, column.DataType)
	col += m.buildColumnAttribues(column)
	if previousColumnName != "" {
		col += fmt.Sprintf(" AFTER `%s`", previousColumnName)
	}
	col += ";"
	return col
}

func (m *DDLBuilder) AddConstraint(tableName string, column *model.Column) string {
	ref := strings.Split(column.ForeignKey, ".")
	return fmt.Sprintf("ALTER TABLE `%s` ADD CONSTRAINT `fk_%s_%s` FOREIGN KEY IF NOT EXISTS (`%s`) REFERENCES `%s` (`%s`);",
		tableName, tableName, column.Name, column.Name, ref[0], ref[1])
}

func (m *DDLBuilder) AlterTableModifyColumn(tableName string, column *model.Column) string {
	return fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN IF EXISTS `%s` %s%s;",
		tableName, column.Name, column.DataType, m.buildColumnAttribues(column))
}

// DropIndex implements db.DDL.
func (m *DDLBuilder) DropIndex(tableName, columnName string) string {
	return fmt.Sprintf("DROP INDEX IF EXISTS `idx_%s_%s` on `%s`;", tableName, columnName, tableName)
}

func (m *DDLBuilder) DropConstraint(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE IF EXISTS `%s` DROP CONSTRAINT IF EXISTS `fk_%s_%s`;", tableName, tableName, columnName)
}

// DropTable implements db.DDL.
func (m *DDLBuilder) DropTable(table *model.Table) []string {
	var stmts []string
	if table.Version {
		stmts = append(stmts, fmt.Sprintf("ALTER TABLE IF EXISTS `%s` DROP SYSTEM VERSIONING;", table.Name))
	}
	stmts = append(stmts, fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", table.Name))
	return stmts
}

func (m *DDLBuilder) buildCreateTableStmt(tableName, tableDesc string, columns []*model.Column) string {
	var stmts []string
	var pkColumn *model.Column
	for _, column := range columns {
		col := fmt.Sprintf("  `%s` %s", column.Name, column.DataType)
		col += m.buildColumnAttribues(column)
		stmts = append(stmts, col)

		if util.IsYes(column.Identity) {
			pkColumn = column
		}
	}
	// primary index
	if pkColumn != nil {
		stmts = append(stmts, fmt.Sprintf("  PRIMARY KEY (`%s`)", pkColumn.Name))
	}
	joinedRows := strings.Join(stmts, ",\n")
	return fmt.Sprintf("CREATE TABLE `%s` IF NOT EXISTS (\n%s\n)%s;\n",
		tableName, joinedRows, lo.If(tableDesc != "", fmt.Sprintf(" COMMENT `%s`", tableDesc)).Else(""))
}

func (m *DDLBuilder) buildColumnAttribues(col *model.Column) string {
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

/* -------------------------------- Separator ------------------------------- */
// type DDLService struct {
// 	tableName string
// }

// func NewDDLService() *DDLService {
// 	return &DDLService{}
// }

// func (m *DDLService) WithTableName(tableName string) *DDLService {
// 	m.tableName = tableName
// 	return m
// }

// func (m *DDLService) GenerateStatement(db *model.Schema, stmtType string) (*string, error) {
// 	switch stmtType {
// 	case "create":
// 		return m.createStmt(db)
// 	case "alter":
// 		return m.alterStmt(db)
// 	default:
// 		return nil, tracerr.Wrap(fmt.Errorf("unsupported statement type: %s", stmtType))
// 	}
// }

// func (m *DDLService) createStmt(mdb *model.Schema) (*string, error) {
// 	mTables := mdb.Tables
// 	if m.tableName != "" {
// 		mTables = lo.Filter(mTables, func(mTable *model.Table, _ int) bool {
// 			return mTable.Name == m.tableName
// 		})
// 	}

// 	var stmts []string
// 	for _, mTable := range mTables {
// 		stmt, err := m.createTableStmt(mTable)
// 		if err != nil {
// 			return nil, tracerr.Wrap(err)
// 		}
// 		stmts = append(stmts, stmt...)
// 	}

// 	output := strings.Join(stmts, "\n")
// 	return &output, nil
// }

// func (m *DDLService) createTableStmt(mTable *model.Table) ([]string, error) {
// 	var stmts, colStmts []string
// 	for _, mColumn := range mTable.Columns {
// 		col := fmt.Sprintf("\t`%s` %s", mColumn.Name, mColumn.DataType)
// 		col += m.buildColumnAttribues(mColumn)
// 		colStmts = append(colStmts, col)
// 	}

// 	// primary index
// 	for _, column := range mTable.Columns {
// 		if util.IsYes(column.Identity) {
// 			colStmts = append(colStmts, fmt.Sprintf("\tPRIMARY KEY (`%s`)", column.Name))
// 		}
// 	}

// 	joinedItems := strings.Join(colStmts, ",\n")
// 	if mTable.Desc == "" {
// 		mTable.Desc = mTable.Name
// 	}
// 	stmts = append(stmts, fmt.Sprintf("CREATE TABLE `%s` IF NOT EXISTS (\n%s\n) COMMENT `%s`;\n", mTable.Name, joinedItems, mTable.Desc))

// 	// column index
// 	for _, column := range mTable.Columns {
// 		index := m.buildColumnIndex(mTable.Name, column)
// 		if index != "" {
// 			stmts = append(stmts, index)
// 		}
// 	}

// 	// additional index
// 	for i, index := range mTable.Indexes {
// 		keyName := index.Name
// 		if keyName == "" {
// 			keyName = fmt.Sprintf("idx_%s_%s", mTable.Name, fmt.Sprintf("%03d", (i+1)))
// 		}
// 		stmts = append(stmts, m.createIndex(keyName, mTable.Name, index.Columns, util.IsYes(index.Unique)))
// 	}

// 	// foreign key
// 	fkColumns := lo.Filter(mTable.Columns, func(column *model.Column, index int) bool {
// 		return column.ForeignKey != ""
// 	})
// 	for _, fkc := range fkColumns {
// 		ref := strings.Split(fkc.ForeignKey, ".")
// 		stmts = append(stmts, fmt.Sprintf("ALTER TABLE `%s` ADD CONSTRAINT `fk_%s_%s` FOREIGN KEY (`%s`) REFERENCES `%s` (`%s`);", mTable.Name, mTable.Name, fkc.Name, fkc.Name, ref[0], ref[1]))
// 	}

// 	return stmts, nil
// }

// func (m *DDLService) alterStmt(mdb *model.Schema) (*string, error) {
// 	mTables := mdb.Tables
// 	if m.tableName != "" {
// 		mTables = lo.Filter(mTables, func(mTable *model.Table, _ int) bool {
// 			return mTable.Name == m.tableName
// 		})
// 	}

// 	var stmts []string
// 	for _, mTable := range mTables {
// 		stmt, err := m.alterTableStmt(mdb.Name, mTable)
// 		if err != nil {
// 			return nil, tracerr.Wrap(err)
// 		}
// 		stmts = append(stmts, stmt...)
// 	}

// 	output := strings.Join(stmts, "\n")
// 	return &output, nil
// }

// func (m *DDLService) alterTableStmt(srcDBName string, srcTB *model.Table) ([]string, error) {
// 	// map current database schema to model
// 	expDB, err := NewExportBuilder().WithDatabaseName(srcDBName).ToSchemaModel()
// 	if err != nil {
// 		return nil, tracerr.Wrap(err)
// 	}

// 	stmts := []string{}

// 	expTB, ok := lo.Find(expDB.Tables, func(t *model.Table) bool {
// 		return t.Name == srcTB.Name
// 	})
// 	if !ok {
// 		/* ------------------------------ create table ------------------------------ */
// 		stmt, err := m.createTableStmt(srcTB)
// 		if err != nil {
// 			return nil, tracerr.Wrap(err)
// 		}
// 		stmts = append(stmts, stmt...)
// 	} else {
// 		/* ------------------------------- alter table ------------------------------ */
// 		// column ID support function
// 		ciRe := regexp.MustCompile(`\[(\d+)\]`)
// 		findCI := func(v string) string {
// 			match := ciRe.FindStringSubmatch(v)
// 			if len(match) > 1 { // column ID found
// 				return match[1]
// 			}
// 			return ""
// 		}

// 		eofStmts := []string{}
// 		for _, srcCol := range srcTB.Columns {
// 			srcCI := findCI(srcCol.Desc)
// 			if srcCI == "" { // column ID does not found
// 				continue
// 			}
// 			expCol, ok := lo.Find(expTB.Columns, func(c *model.Column) bool {
// 				return findCI(c.Desc) == srcCI
// 			})
// 			if !ok { // column ID not found
// 				continue
// 			}

// 			if srcCol.DataType != expCol.DataType || srcCol.NotNull != expCol.NotNull || srcCol.Value != expCol.Value || srcCol.Desc != expCol.Desc || srcCol.Unique != expCol.Unique {
// 				stmts = append(stmts, fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN `%s` %s%s;", srcTB.Name, expCol.Name, srcCol.DataType, m.buildColumnAttribues(srcCol)))
// 			}
// 			if srcCol.Name != expCol.Name {
// 				stmts = append(stmts, fmt.Sprintf("ALTER TABLE `%s` RENAME COLUMN `%s` TO `%s`;", srcTB.Name, expCol.Name, srcCol.Name))
// 			}

// 			if srcCol.Index != expCol.Index || srcCol.Unique != expCol.Unique {
// 				eofStmts = append(eofStmts, fmt.Sprintf("DROP INDEX IF EXISTS `%s` ON `%s`;", m.formatIndexName(srcTB.Name, expCol.Name), srcTB.Name))
// 				eofStmts = append(eofStmts, m.buildColumnIndex(srcTB.Name, srcCol))
// 			}
// 		}
// 		stmts = append(stmts, eofStmts...)
// 	}
// 	return stmts, nil
// }

// func (m *DDLService) buildColumnAttribues(col *model.Column) string {
// 	attr := ""
// 	if col.Compute != "" {
// 		attr += fmt.Sprintf(" GENERATED ALWAYS AS (%s)", col.Compute)
// 	}
// 	if util.IsYes(col.NotNull) {
// 		attr += " NOT NULL"
// 	}
// 	if col.Value != "" {
// 		attr += fmt.Sprintf(" DEFAULT %s", col.Value)
// 	}
// 	// if util.IsYes(col.AutoIncrement) {
// 	// 	attr += " AUTO_INCREMENT"
// 	// }
// 	if col.Desc != "" {
// 		attr += fmt.Sprintf(" COMMENT `%s`", col.Desc)
// 	}
// 	return attr
// }

// func (m *DDLService) buildColumnIndex(tbn string, col *model.Column) string {
// 	if util.IsYes(col.Index) {
// 		name := m.formatIndexName(tbn, col.Name)
// 		return m.createIndex(name, tbn, []string{col.Name}, util.IsYes(col.Unique))
// 	}
// 	return ""
// }

// func (m *DDLService) createIndex(name, tbn string, fields []string, isUnique bool) string {
// 	unique := lo.If(isUnique, "UNIQUE ").Else("")
// 	return fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS `%s` ON `%s` (`%s`);", unique, name, tbn, strings.Join(fields, "`,`"))
// }

// func (m *DDLService) formatIndexName(tb, col string) string {
// 	return fmt.Sprintf("idx_%s_%s", tb, col)
// }
