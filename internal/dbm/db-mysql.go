package dbm

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-sqlx/sqlx"
	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

func ReadMYSQL(dataSourceName string, schemaName string) (*model.DataDef, error) {
	dataDef := model.DataDef{}

	fmt.Println("Connecting to MySQL ...")
	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	mSchemas, err := mysqlSchemas(db, schemaName)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	dataDef.Schemas = *mSchemas
	return &dataDef, nil
}

func mysqlSchemas(db *sqlx.DB, schemaName string) (*[]model.Schema, error) {
	mSchemas := []model.Schema{}
	dSchemas := []MySqlSchema{}
	if err := db.Select(&dSchemas, "select * from INFORMATION_SCHEMA.SCHEMATA where SCHEMA_NAME = ?", schemaName); err != nil {
		return nil, tracerr.Wrap(err)
	}
	for _, dSchema := range dSchemas {
		mTables, err := mysqlTable(db, dSchema)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		mSchemas = append(mSchemas, model.Schema{
			Name:   *dSchema.SchemaName,
			Tables: *mTables,
		})
	}
	return &mSchemas, nil
}

func mysqlTable(db *sqlx.DB, dSchema MySqlSchema) (*[]model.Table, error) {
	dTables := []MySqlTable{}
	if err := db.Select(&dTables,
		"select * from INFORMATION_SCHEMA.TABLES where TABLE_TYPE = 'BASE TABLE' and TABLE_SCHEMA = ? order by TABLE_NAME",
		dSchema.SchemaName); err != nil {
		return nil, tracerr.Wrap(err)
	}

	mTables := make([]model.Table, len(dTables))

	for index, dTable := range dTables {
		mColumns, err := mysqlColumn(db, dTable)
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

func mysqlColumn(db *sqlx.DB, dTable MySqlTable) (*[]model.Column, error) {
	dColumns := []MySqlColumn{}
	if err := db.Select(&dColumns,
		"select * from INFORMATION_SCHEMA.COLUMNS where TABLE_SCHEMA = ? and TABLE_NAME = ? order by ORDINAL_POSITION",
		dTable.TableSchema, dTable.TableName); err != nil {
		return nil, tracerr.Wrap(err)
	}

	mColumns := make([]model.Column, len(dColumns))

	for index, dColumn := range dColumns {
		mColumn := model.Column{
			Name:     *dColumn.ColumnName,
			DataType: *dColumn.DataType,
			NotNull:  lo.Ternary(*dColumn.IsNullable == "YES", "Y", "N"),
		}
		if dColumn.NumericPrecision != nil && *dColumn.NumericPrecision > 0 {
			mColumn.DataType = fmt.Sprintf("%s(%d)", *dColumn.DataType, *dColumn.NumericPrecision)
			if dColumn.NumericScale != nil && *dColumn.NumericScale > 0 {
				mColumn.DataType = fmt.Sprintf("%s(%d,%d)", *dColumn.DataType, *dColumn.NumericPrecision, *dColumn.NumericScale)
			}
		}
		mColumns[index] = mColumn
	}
	return &mColumns, nil
}
