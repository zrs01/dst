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
