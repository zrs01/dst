package mariadb

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/zrs01/dst/util"
	"github.com/ztrue/tracerr"
)

type ExportService struct {
	driverName         string
	dataSourceName     string
	targetDatabaseName string
	targetTableName    string
	commonColumnFile   string
}

func NewExportService() *ExportService {
	return &ExportService{}
}

func (m *ExportService) WithDriverName(driverName string) *ExportService {
	m.driverName = driverName
	return m
}

func (m *ExportService) WithDataSourceName(dataSourceName string) *ExportService {
	m.dataSourceName = dataSourceName
	return m
}

func (m *ExportService) WithDatabaseName(targetDatabaseName string) *ExportService {
	m.targetDatabaseName = targetDatabaseName
	return m
}

func (m *ExportService) WithTableName(targetTableName string) *ExportService {
	m.targetTableName = targetTableName
	return m
}

func (m *ExportService) WithCommonColumnFile(commonColumnFile string) *ExportService {
	m.commonColumnFile = commonColumnFile
	return m
}

func (m *ExportService) ToModel() (*model.Schema, error) {
	if m.targetDatabaseName == "" {
		return nil, tracerr.Errorf("database name is required")
	}

	excludeColumns := &[]*model.Column{}
	if m.commonColumnFile != "" {
		if _, err := os.Stat(m.commonColumnFile); os.IsNotExist(err) {
			return nil, tracerr.Errorf("file '%s' does not exist", m.commonColumnFile)
		}
		ec, err := util.UnmarshalYml(m.commonColumnFile, &model.Table{})
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		excludeColumns = &ec.Columns
	}

	db, err := util.OpenSqlDb(m.driverName, m.dataSourceName)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer db.Close()

	/* ----------------------- Get the database definition ---------------------- */
	scrDbm := NewDatabaseManager()
	scrDbs, error := scrDbm.WithSchemaName(m.targetDatabaseName).GetDef(db)
	if error != nil {
		return nil, tracerr.Wrap(error)
	}
	if len(scrDbs) == 0 {
		return nil, tracerr.Errorf("The database (%s) was not found", m.targetDatabaseName)
	}

	desDb := &model.Schema{}
	desDb.Name = scrDbs[0].Name
	desDb.Desc = scrDbs[0].Comment

	/* ------------------------ Get the table definition ------------------------ */
	{
		srcTbm := NewTableManager()
		if m.targetTableName != "" {
			srcTbm.WithTableName(m.targetTableName)
		}
		srcTbs, err := srcTbm.WithTableSchema(m.targetDatabaseName).WithTableType("BASE TABLE").GetDef(db)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		desDb.Tables = []*model.Table{}
		for _, stb := range srcTbs {
			desTb := &model.Table{}
			desTb.Name = stb.TableName
			desTb.Desc = stb.Comment
			desTb.Version = stb.TableType == "VERSIONED"

			desDb.Tables = append(desDb.Tables, desTb)
		}
	}

	/* ------------------------ Get the column definition ----------------------- */
	{
		scrColm := NewColumnManager()
		scrCols, err := scrColm.WithTableSchema(m.targetDatabaseName).GetDef(db)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		for _, t := range desDb.Tables {

			cids := []string{} // storage of column ID of the table
			cidRe := regexp.MustCompile(`\[(\d+)\]`)

			t.Columns = []*model.Column{}
			for _, sc := range scrCols {
				if sc.TableName == t.Name {
					// check excluded column
					isExclude := lo.ContainsBy(*excludeColumns, func(x *model.Column) bool {
						return x.Name == sc.ColumnName && x.DataType == sc.ColumnType
					})

					if !isExclude {
						desCol := &model.Column{}

						desCol.Name = sc.ColumnName
						desCol.DataType = sc.ColumnType
						desCol.Identity = lo.If(sc.ColumnKey == "PRI", "Y").Else("N")
						desCol.NotNull = lo.If(sc.IsNullable == "YES", "N").Else("Y")
						desCol.Unique = lo.If(sc.ColumnKey == "UNI" || sc.ColumnKey == "PRI", "Y").Else("N")
						desCol.Value = lo.If(sc.ColumnDefault != "" && sc.ColumnDefault != "NULL", sc.ColumnDefault).Else("")
						desCol.Desc = strings.TrimSpace(sc.Comment)
						if sc.IsGenerated == "ALWAYS" {
							desCol.Compute = sc.GenerationExpression
							if strings.Contains(sc.Extra, "STORED") {
								desCol.ComputeType = "stored"
							}
						}

						// collect column ID
						match := cidRe.FindStringSubmatch(desCol.Desc)
						if len(match) > 1 { // column ID found
							cids = append(cids, match[1])
						}
						t.Columns = append(t.Columns, desCol)
					}
				}
			}

			// fix column ID
			if len(cids) != len(t.Columns) {
				count := 1
				for i := 0; i < len(t.Columns); i++ {
					dc := t.Columns[i]
					match := cidRe.FindStringSubmatch(dc.Desc)
					if len(match) < 2 { // column ID doesn't found
						countStr := fmt.Sprintf("%03d", count)
						for lo.Contains(cids, countStr) {
							count++
							countStr = fmt.Sprintf("%03d", count)
						}
						// add ID to desc
						dc.Desc = fmt.Sprintf("[%s] %s", countStr, dc.Desc)
						cids = append(cids, countStr)
					}
				}
			}
		}
	}

	/* ------------------------ Get forign key definition ----------------------- */
	{
		srcConstm := NewConstraintManager()
		for _, desTb := range desDb.Tables {
			// Get constraints of the table
			scrConsts, err := srcConstm.WithTableSchema(m.targetDatabaseName).WithTableName(desTb.Name).GetDef(db)
			if err != nil {
				return nil, tracerr.Wrap(err)
			}

			if len(scrConsts) > 0 {
				for _, srcCol := range desTb.Columns {
					// Get foreign key
					found, ok := lo.Find(scrConsts, func(i *Constraint) bool {
						return i.ConstraintName != "PRIMARY" && i.TableSchema == m.targetDatabaseName && i.TableName == desTb.Name && i.ColumnName == srcCol.Name
					})
					if ok {
						srcCol.ForeignKey = fmt.Sprintf("%s.%s", found.ReferencedTableName, found.ReferencedColumnName)
					}
				}
			}
		}
	}

	/* -------------------------- Get index definition -------------------------- */
	{
		srcIdxm := NewIndexManager()
		for _, desTb := range desDb.Tables {
			// Get indexes of the table
			scrIdxs, err := srcIdxm.WithTableSchema(m.targetDatabaseName).WithTableName(desTb.Name).GetDef(db)
			if err != nil {
				return nil, tracerr.Wrap(err)
			}

			if len(scrIdxs) > 0 {
				for _, srcCol := range desTb.Columns {
					_, ok := lo.Find(scrIdxs, func(i *Index) bool {
						return i.TableSchema == m.targetDatabaseName && i.TableName == desTb.Name && i.IndexName != "PRIMARY" && len(i.Columns) == 1 && i.Columns[0] == srcCol.Name
					})
					if ok {
						srcCol.Index = "Y"
					}
				}
				// multi column index
				desTb.Indexes = lo.FilterMap(scrIdxs, func(i *Index, _ int) (*model.Index, bool) {
					if i.TableSchema == m.targetDatabaseName && i.TableName == desTb.Name && len(i.Columns) > 1 {
						return &model.Index{
							Name:    i.IndexName,
							Unique:  lo.If(i.Unique, "Y").Else("N"),
							Columns: i.Columns,
						}, true
					}
					return nil, false
				})
			}
		}
	}

	return desDb, nil
}

func (m *ExportService) ToRoutines() ([]*model.Routine, error) {
	db, err := util.OpenSqlDb(m.driverName, m.dataSourceName)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer db.Close()

	routineManager := NewRoutineManager()
	srcRoutines, err := routineManager.WithRoutineSchema(m.targetDatabaseName).GetDef(db)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	desRoutines := []*model.Routine{}
	for _, srcRt := range srcRoutines {
		desRoutines = append(desRoutines, &model.Routine{
			Name: srcRt.RoutineName,
			Code: srcRt.RoutineDefinition,
		})
	}
	return desRoutines, nil
}
