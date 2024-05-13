package dbm

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	// "github.com/go-sqlx/sqlx"
	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

type MssqlService struct{}

func NewMssqlService() DataSource {
	return &MssqlService{}
}

func (s *MssqlService) Read(dataSourceName string, dbName string) (*model.DataDef, error) {
	dataDef := model.DataDef{}

	fmt.Println("Connecting to SqlServer ...")
	db, err := sqlx.Connect("sqlserver", dataSourceName)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	dSchemas, err := s.getSchemas(db, "dbo")
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	dTables, err := s.getTables(db, dSchemas, dbName)
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

func (s *MssqlService) buildSchemas(dSchemas *[]MsSqlSchema, dTables *[]MsSqlTable, dColumns *[]MsSqlColumn, dKeyColumnUsages *[]MsSqlKeyColumnUsage) (*[]model.Schema, error) {
	mSchemas := make([]model.Schema, len(*dSchemas))
	for index, dSchema := range *dSchemas {
		schemaTables := lo.Filter(*dTables, func(dTable MsSqlTable, _ int) bool {
			return *dTable.TableSchema == *dSchema.SchemaName
		})
		mTables, err := s.buildTable(dSchema, &schemaTables, dColumns, dKeyColumnUsages)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		mSchemas[index] = model.Schema{
			Name:   *dSchema.SchemaName,
			Tables: *mTables,
		}
	}
	return &mSchemas, nil
}

func (s *MssqlService) buildTable(dSchema MsSqlSchema, dTables *[]MsSqlTable, dColumns *[]MsSqlColumn, dKeyColumnUsages *[]MsSqlKeyColumnUsage) (*[]model.Table, error) {
	mTables := make([]model.Table, len(*dTables))
	for index, dTable := range *dTables {
		tableColumns := lo.Filter(*dColumns, func(dColumn MsSqlColumn, _ int) bool {
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

func (s *MssqlService) buildColumn(dTable MsSqlTable, dColumns *[]MsSqlColumn, dKeyColumnUsages *[]MsSqlKeyColumnUsage) (*[]model.Column, error) {
	mColumns := make([]model.Column, len(*dColumns))
	for index, dColumn := range *dColumns {
		mColumn := model.Column{
			Name:     *dColumn.ColumnName,
			DataType: *dColumn.DataType,
		}

		if *dColumn.IsNullable != "YES" {
			mColumn.NotNull = "Y"
		}

		if _, ok := lo.Find(*dKeyColumnUsages, func(dKeyColumnUsage MsSqlKeyColumnUsage) bool {
			return *dKeyColumnUsage.ConstraintName == "PRIMARY" && *dKeyColumnUsage.ColumnName == *dColumn.ColumnName
		}); ok {
			mColumn.Identity = "Y"
		}

		// if usage, ok := lo.Find(*dKeyColumnUsages, func(dKeyColumnUsage MsSqlKeyColumnUsage) bool {
		// 	return *dKeyColumnUsage.ConstraintName != "PRIMARY" && *dKeyColumnUsage.ColumnName == *dColumn.ColumnName
		// }); ok {
		// 	mColumn.ForeignKey = fmt.Sprintf("%s.%s", *usage.ReferencedTableName, *usage.ReferencedColumnName)
		// }

		if strings.Contains(strings.ToLower(*dColumn.DataType), "char") {
			mColumn.DataType = fmt.Sprintf("%s(%d)", *dColumn.DataType, *dColumn.CharacterMaximumLength)
		} else if strings.Contains(strings.ToLower(*dColumn.DataType), "decimal") {
			if dColumn.NumericScale != nil && *dColumn.NumericScale > 0 {
				mColumn.DataType = fmt.Sprintf("%s(%d,%d)", *dColumn.DataType, *dColumn.NumericPrecision, *dColumn.NumericScale)
			}
		}
		mColumns[index] = mColumn
	}
	return &mColumns, nil
}

func (s *MssqlService) getSchemas(db *sqlx.DB, schemaName string) (*[]MsSqlSchema, error) {
	dSchemas := []MsSqlSchema{}
	if err := db.Select(&dSchemas, "select * from INFORMATION_SCHEMA.SCHEMATA where SCHEMA_NAME = @SCHEMA",
		sql.Named("SCHEMA", schemaName)); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &dSchemas, nil
}

func (s *MssqlService) getTables(db *sqlx.DB, schemas *[]MsSqlSchema, dbName string) (*[]MsSqlTable, error) {
	dTables := []MsSqlTable{}
	dSchemaNames := lo.Map(*schemas, func(dSchema MsSqlSchema, _ int) string { return *dSchema.SchemaName })
	query, args, err := sqlx.In(`
		select * from INFORMATION_SCHEMA.TABLES where TABLE_TYPE = 'BASE TABLE' and TABLE_CATALOG = @CAT and TABLE_SCHEMA in (@SCHEMAS) order by TABLE_NAME`,
		sql.Named("CAT", dbName), sql.Named("SCHEMAS", strings.Join(dSchemaNames, ",")))
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	query = db.Rebind(query)
	if err := db.Select(&dTables, query, args...); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &dTables, nil
}

func (s *MssqlService) getColumns(db *sqlx.DB, schemas *[]MsSqlSchema, tables *[]MsSqlTable) (*[]MsSqlColumn, error) {
	dColumns := []MsSqlColumn{}
	dSchemaNames := lo.Map(*schemas, func(dSchema MsSqlSchema, _ int) string { return *dSchema.SchemaName })
	dTableNames := lo.Map(*tables, func(dTable MsSqlTable, _ int) string { return *dTable.TableName })
	// query, args, err := sqlx.In(
	// 	`select * from INFORMATION_SCHEMA.COLUMNS where TABLE_SCHEMA in (@SCHEMAS) and TABLE_NAME in (@TABLES) order by TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION`,
	// 	sql.Named("SCHEMAS", strings.Join(dSchemaNames, `,`)), sql.Named("TABLES", strings.Join(dTableNames, `,`)))
	// if err != nil {
	// 	return nil, tracerr.Wrap(err)
	// }
	// fmt.Println(query)
	// query = db.Rebind(query)
	// if err := db.Select(&dColumns, query, args...); err != nil {
	// 	return nil, tracerr.Wrap(err)
	// }
	if err := db.Select(&dColumns,
		`select * from INFORMATION_SCHEMA.COLUMNS where TABLE_SCHEMA in ($1) and TABLE_NAME in ($2) order by TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION`,
		strings.Join(dSchemaNames, `,`), strings.Join(dTableNames, `,`)); err != nil {
		return nil, tracerr.Wrap(err)
	}

	// spew.Dump(args)
	// spew.Dump(dSchemaNames)
	// spew.Dump(dTableNames)
	spew.Dump(dColumns)
	return &dColumns, nil
}

func (s *MssqlService) getKeyColumnUsages(db *sqlx.DB, schemas *[]MsSqlSchema, tables *[]MsSqlTable) (*[]MsSqlKeyColumnUsage, error) {
	dColumnUsages := []MsSqlKeyColumnUsage{}
	dSchemaNames := lo.Map(*schemas, func(dSchema MsSqlSchema, _ int) string { return *dSchema.SchemaName })
	dTableNames := lo.Map(*tables, func(dTable MsSqlTable, _ int) string { return *dTable.TableName })
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

// func ReadMssql(dataSourceName string) (*model.DataDef, error) {
// 	dataDef := model.DataDef{}

// 	db, err := sqlx.Connect("sqlserver", dataSourceName)
// 	if err != nil {
// 		return nil, tracerr.Wrap(err)
// 	}

// 	mSchemas, err := mssqlSchemas(db)
// 	if err != nil {
// 		return nil, tracerr.Wrap(err)
// 	}
// 	dataDef.Schemas = *mSchemas
// 	return &dataDef, nil
// }

// func mssqlSchemas(db *sqlx.DB) (*[]model.Schema, error) {
// 	mSchemas := []model.Schema{}
// 	dSchemas := []MsSqlSchema{}
// 	if err := db.Select(&dSchemas, "select * from INFORMATION_SCHEMA.SCHEMATA where SCHEMA_NAME = 'dbo'"); err != nil {
// 		return nil, tracerr.Wrap(err)
// 	}
// 	for _, dSchema := range dSchemas {
// 		mTables, err := mssqlTable(db, dSchema)
// 		if err != nil {
// 			return nil, tracerr.Wrap(err)
// 		}
// 		mSchemas = append(mSchemas, model.Schema{
// 			Name:   dSchema.SchemaName,
// 			Tables: *mTables,
// 		})
// 	}
// 	return &mSchemas, nil
// }

// func mssqlTable(db *sqlx.DB, dSchema MsSqlSchema) (*[]model.Table, error) {
// 	mTables := []model.Table{}

// 	dTables := []MsSqlTable{}
// 	if err := db.Select(&dTables, "select * from INFORMATION_SCHEMA.TABLES where TABLE_TYPE = 'BASE TABLE' order by TABLE_SCHEMA,TABLE_NAME"); err != nil {
// 		return nil, tracerr.Wrap(err)
// 	}
// 	for _, dTable := range dTables {
// 		mColumns, err := mssqlColumn(db, dTable)
// 		if err != nil {
// 			return nil, tracerr.Wrap(err)
// 		}
// 		mTables = append(mTables, model.Table{
// 			Name:    dTable.TableName,
// 			Columns: *mColumns,
// 		})
// 	}
// 	return &mTables, nil
// }

// func mssqlColumn(db *sqlx.DB, dTable MsSqlTable) (*[]model.Column, error) {
// 	mColumns := []model.Column{}

// 	dColumns := []MsSqlColumn{}
// 	if err := db.Select(&dColumns, fmt.Sprintf("select * from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME = '%s' order by ORDINAL_POSITION", dTable.TableName)); err != nil {
// 		return nil, tracerr.Wrap(err)
// 	}
// 	for _, dColumn := range dColumns {
// 		mColumn := model.Column{
// 			Name:     dColumn.ColumnName,
// 			DataType: dColumn.DataType,
// 		}
// 		if dColumn.NumericPrecision != nil && *dColumn.NumericPrecision > 0 {
// 			mColumn.DataType = fmt.Sprintf("%s(%d)", dColumn.DataType, *dColumn.NumericPrecision)
// 			if dColumn.NumericScale != nil && *dColumn.NumericScale > 0 {
// 				mColumn.DataType = fmt.Sprintf("%s(%d,%d)", dColumn.DataType, *dColumn.NumericPrecision, *dColumn.NumericScale)
// 			}
// 		}
// 		mColumns = append(mColumns, mColumn)
// 	}
// 	return &mColumns, nil
// }
