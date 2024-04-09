package dbm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-sqlx/sqlx"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
)

func ReadMYSQL(dataSourceName string, schemaName string) (*model.DataDef, error) {
	dataDef := model.DataDef{}

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
	if err := db.Select(&dSchemas, "select * from INFORMATION_SCHEMA.SCHEMATA where SCHEMA_NAME = '?'", schemaName); err != nil {
		return nil, tracerr.Wrap(err)
	}
	for _, dSchema := range dSchemas {
		mTables, err := mysqlTable(db, dSchema)
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

func mysqlTable(db *sqlx.DB, dSchema MySqlSchema) (*[]model.Table, error) {
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
