package dbm

import (
	"fmt"

	"github.com/go-sqlx/sqlx"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

func ReadMSSQL(in string) (*model.DataDef, error) {
	dataDef := model.DataDef{}

	db, err := sqlx.Connect("sqlserver", in)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	mSchemas, err := mssqlSchemas(db)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	dataDef.Schemas = *mSchemas
	return &dataDef, nil
}

func mssqlSchemas(db *sqlx.DB) (*[]model.Schema, error) {
	mSchemas := []model.Schema{}
	dSchemas := []MsSqlSchema{}
	if err := db.Select(&dSchemas, "select * from INFORMATION_SCHEMA.SCHEMATA where SCHEMA_NAME = 'dbo'"); err != nil {
		return nil, tracerr.Wrap(err)
	}
	for _, dSchema := range dSchemas {
		mTables, err := mssqlTable(db, dSchema)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		mSchemas = append(mSchemas, model.Schema{
			Name:   dSchema.SchemaName,
			Tables: *mTables,
		})
	}
	return &mSchemas, nil
}

func mssqlTable(db *sqlx.DB, dSchema MsSqlSchema) (*[]model.Table, error) {
	mTables := []model.Table{}

	dTables := []MsSqlTable{}
	if err := db.Select(&dTables, "select * from INFORMATION_SCHEMA.TABLES where TABLE_TYPE = 'BASE TABLE' order by TABLE_SCHEMA,TABLE_NAME"); err != nil {
		return nil, tracerr.Wrap(err)
	}
	for _, dTable := range dTables {
		mColumns, err := mssqlColumn(db, dTable)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		mTables = append(mTables, model.Table{
			Name:    dTable.TableName,
			Columns: *mColumns,
		})
	}
	return &mTables, nil
}

func mssqlColumn(db *sqlx.DB, dTable MsSqlTable) (*[]model.Column, error) {
	mColumns := []model.Column{}

	dColumns := []MsSqlColumn{}
	if err := db.Select(&dColumns, fmt.Sprintf("select * from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME = '%s' order by ORDINAL_POSITION", dTable.TableName)); err != nil {
		return nil, tracerr.Wrap(err)
	}
	for _, dColumn := range dColumns {
		mColumn := model.Column{
			Name:     dColumn.ColumnName,
			DataType: dColumn.DataType,
		}
		if dColumn.NumericPrecision != nil && *dColumn.NumericPrecision > 0 {
			mColumn.DataType = fmt.Sprintf("%s(%d)", dColumn.DataType, *dColumn.NumericPrecision)
			if dColumn.NumericScale != nil && *dColumn.NumericScale > 0 {
				mColumn.DataType = fmt.Sprintf("%s(%d,%d)", dColumn.DataType, *dColumn.NumericPrecision, *dColumn.NumericScale)
			}
		}
		mColumns = append(mColumns, mColumn)
	}
	return &mColumns, nil
}