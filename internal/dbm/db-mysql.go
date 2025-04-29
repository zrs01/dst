package dbm

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-sqlx/sqlx"
	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

type MysqlService struct {
	dataSourceName string
}

func NewMysqlService(dataSourceName string) DataService {
	return &MysqlService{
		dataSourceName: dataSourceName,
	}
}

func (s *MysqlService) DataSourceName() string {
	return s.dataSourceName
}

func (s *MysqlService) Read(schemaName string) (*model.DataDef, error) {
	dataDef := model.DataDef{}

	fmt.Println("Connecting to MySQL ...")
	db, err := sqlx.Connect("mysql", s.dataSourceName)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer db.Close()

	dSchemas, err := s.getSchemas(db, schemaName)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	dTables, err := s.getTables(db, dSchemas)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	dColumns, err := s.getColumns(db, dSchemas, dTables)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	dKeyColumnUsages, err := s.getKeyColumnUsages(db, dSchemas, dTables)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	mSchemas, err := s.buildSchemas(dSchemas, dTables, dColumns, dKeyColumnUsages)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	dataDef.Schemas = *mSchemas
	return &dataDef, nil
}

func (s *MysqlService) buildSchemas(dSchemas *[]MySqlSchema, dTables *[]MySqlTable, dColumns *[]MySqlColumn, dKeyColumnUsages *[]MySqlKeyColumnUsage) (*[]*model.Schema, error) {
	mSchemas := make([]*model.Schema, len(*dSchemas))
	for index, dSchema := range *dSchemas {
		schemaTables := lo.Filter(*dTables, func(dTable MySqlTable, _ int) bool {
			return *dTable.TableSchema == *dSchema.SchemaName
		})
		mTables, err := s.buildTable(dSchema, &schemaTables, dColumns, dKeyColumnUsages)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		mSchemas[index] = &model.Schema{
			Name:   *dSchema.SchemaName,
			Tables: *mTables,
		}
		if *dSchema.SchemaComment != "" {
			mSchemas[index].Desc = *dSchema.SchemaComment
		}
	}
	return &mSchemas, nil
}

func (s *MysqlService) buildTable(dSchema MySqlSchema, dTables *[]MySqlTable, dColumns *[]MySqlColumn, dKeyColumnUsages *[]MySqlKeyColumnUsage) (*[]model.Table, error) {
	mTables := make([]model.Table, len(*dTables))
	for index, dTable := range *dTables {
		tableColumns := lo.Filter(*dColumns, func(dColumn MySqlColumn, _ int) bool {
			return *dColumn.TableSchema == *dSchema.SchemaName && *dColumn.TableName == *dTable.TableName
		})
		mColumns, err := s.buildColumn(dTable, &tableColumns, dKeyColumnUsages)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		mTables[index] = model.Table{
			Name:    *dTable.TableName,
			Columns: *mColumns,
		}
	}
	return &mTables, nil
}

func (s *MysqlService) buildColumn(dTable MySqlTable, dColumns *[]MySqlColumn, dKeyColumnUsages *[]MySqlKeyColumnUsage) (*[]*model.Column, error) {
	mColumns := make([]*model.Column, len(*dColumns))
	for index, dColumn := range *dColumns {
		mColumn := &model.Column{
			Name:     *dColumn.ColumnName,
			DataType: *dColumn.DataType,
		}
		// compu field
		if strings.Contains(*dColumn.Extra, "GENERATED") {
			mColumn.Compute = *dColumn.GenerationExpression
		} else {

			if *dColumn.IsNullable != "YES" {
				mColumn.NotNull = "Y"
			}

			if _, ok := lo.Find(*dKeyColumnUsages, func(dKeyColumnUsage MySqlKeyColumnUsage) bool {
				return *dKeyColumnUsage.ConstraintName == "PRIMARY" && *dKeyColumnUsage.ColumnName == *dColumn.ColumnName
			}); ok {
				mColumn.Identity = "Y"
			}

			if usage, ok := lo.Find(*dKeyColumnUsages, func(dKeyColumnUsage MySqlKeyColumnUsage) bool {
				return *dKeyColumnUsage.ConstraintName != "PRIMARY" && *dKeyColumnUsage.ColumnName == *dColumn.ColumnName
			}); ok {
				mColumn.ForeignKey = fmt.Sprintf("%s.%s", *usage.ReferencedTableName, *usage.ReferencedColumnName)
			}

			if strings.Contains(strings.ToLower(*dColumn.DataType), "char") {
				mColumn.DataType = fmt.Sprintf("%s(%d)", *dColumn.DataType, *dColumn.CharacterMaximumLength)
			} else if strings.Contains(strings.ToLower(*dColumn.DataType), "decimal") {
				if dColumn.NumericScale != nil && *dColumn.NumericScale > 0 {
					mColumn.DataType = fmt.Sprintf("%s(%d,%d)", *dColumn.DataType, *dColumn.NumericPrecision, *dColumn.NumericScale)
				}
			}
		}
		if *dColumn.ColumnComment != "" {
			mColumn.Desc = *dColumn.ColumnComment
		}
		mColumns[index] = mColumn
	}
	return &mColumns, nil
}

func (s *MysqlService) getSchemas(db *sqlx.DB, schemaName string) (*[]MySqlSchema, error) {
	dSchemas := []MySqlSchema{}
	if err := db.Select(&dSchemas, "select * from INFORMATION_SCHEMA.SCHEMATA where SCHEMA_NAME = ?", schemaName); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &dSchemas, nil
}

func (s *MysqlService) getTables(db *sqlx.DB, schemas *[]MySqlSchema) (*[]MySqlTable, error) {
	dTables := []MySqlTable{}
	dSchemaNames := lo.Map(*schemas, func(dSchema MySqlSchema, _ int) string { return *dSchema.SchemaName })
	query, args, err := sqlx.In(`
		select * from INFORMATION_SCHEMA.TABLES where TABLE_TYPE = 'BASE TABLE' and TABLE_SCHEMA in (?) order by TABLE_NAME`,
		dSchemaNames)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	query = db.Rebind(query)
	if err := db.Select(&dTables, query, args...); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &dTables, nil
}

func (s *MysqlService) getColumns(db *sqlx.DB, schemas *[]MySqlSchema, tables *[]MySqlTable) (*[]MySqlColumn, error) {
	dColumns := []MySqlColumn{}
	dSchemaNames := lo.Map(*schemas, func(dSchema MySqlSchema, _ int) string { return *dSchema.SchemaName })
	dTableNames := lo.Map(*tables, func(dTable MySqlTable, _ int) string { return *dTable.TableName })
	query, args, err := sqlx.In(
		`select * from INFORMATION_SCHEMA.COLUMNS where TABLE_SCHEMA in (?) and TABLE_NAME in (?) order by TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION`,
		dSchemaNames, dTableNames)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	query = db.Rebind(query)
	if err := db.Select(&dColumns, query, args...); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &dColumns, nil
}

func (s *MysqlService) getKeyColumnUsages(db *sqlx.DB, schemas *[]MySqlSchema, tables *[]MySqlTable) (*[]MySqlKeyColumnUsage, error) {
	dColumnUsages := []MySqlKeyColumnUsage{}
	dSchemaNames := lo.Map(*schemas, func(dSchema MySqlSchema, _ int) string { return *dSchema.SchemaName })
	dTableNames := lo.Map(*tables, func(dTable MySqlTable, _ int) string { return *dTable.TableName })
	query, args, err := sqlx.In(`
		select * from INFORMATION_SCHEMA.KEY_COLUMN_USAGE where TABLE_SCHEMA in (?) and TABLE_NAME in (?) order by TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION `,
		dSchemaNames, dTableNames)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	query = db.Rebind(query)
	if err := db.Select(&dColumnUsages, query, args...); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &dColumnUsages, nil
}
