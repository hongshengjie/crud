package mytable

import (
	"database/sql"
	"io/ioutil"

	"log"
	"sort"
	"strings"

	"github.com/hongshengjie/crud/snaker"
	"github.com/xwb1989/sqlparser"
)

// Table Table
type Table struct {
	Database         string
	TableName        string    // table name
	GoTableName      string    // go struct name
	PackageName      string    // package name
	Fields           []*Column // columns
	Indexes          []*Index  // indexes
	GenerateWhereCol []*Column // GenerateWhereCol 生成where字段比较方法的列
	ConditionsFields []string  // ConditionsFields  用户指定生成的条件方法
	PrimaryKey       *Column   // priomary_key column
	ImportTime       bool      // is need import time
}

// NewTable NewTable
func NewTable(db *sql.DB, database, schema, table string, conditionsFields []string, isAll bool) *Table {
	gotableName := snaker.SnakeToCamelIdentifier(table)
	mytable := &Table{
		Database:    database,
		TableName:   table,
		GoTableName: gotableName,
		PackageName: strings.ToLower(gotableName),
	}
	columns, err := MyTableColumns(db, schema, table)
	if err != nil {
		log.Fatal(err)
	}
	if len(columns) <= 0 {
		log.Fatal("schema or table not exist")
	}
	mytable.Fields = columns
	for _, v := range columns {
		if v.IsPrimaryKey {
			mytable.PrimaryKey = v
		}
		if v.GoColumnType == "time.Time" {
			mytable.ImportTime = true
		}
	}
	if mytable.PrimaryKey == nil {
		log.Fatal("table do not have a primary key")
	}
	if isAll {
		mytable.GenerateWhereCol = mytable.Fields
		return mytable
	}

	// 搞索引
	indexes, err := MyTableIndexes(db, schema, table)
	if err != nil {
		log.Fatal(err)
	}
	indexcols := make(map[string]*Column)
	for _, v := range indexes {
		for _, fieldName := range v.IndexFields {
			for _, c := range columns {
				if c.ColumnName == fieldName {
					v.IndexColumns = append(v.IndexColumns, c)
					indexcols[c.ColumnName] = c
					break
				}
			}
		}
	}
	// 添加索引列之外 指定的field
	for _, v := range conditionsFields {
		for _, c := range columns {
			if c.ColumnName == v {
				indexcols[c.ColumnName] = c
				break
			}
		}
	}
	generateCol := make([]*Column, 0, len(indexcols))
	for _, v := range indexcols {
		generateCol = append(generateCol, v)
	}

	sort.Slice(generateCol, func(i, j int) bool {
		return generateCol[i].OrdinalPosition < generateCol[j].OrdinalPosition
	})
	mytable.GenerateWhereCol = generateCol
	mytable.Indexes = indexes

	return mytable

}

func MytableFromSqlFile(db string, path string) *Table {
	sql, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := sqlparser.ParseStrictDDL(string(sql))
	if err != nil {
		log.Fatal(err)
	}
	ddl, ok := stmt.(*sqlparser.DDL)
	if !ok {
		log.Fatal("please check sql file statement is DDL ")
	}
	if ddl.Action != sqlparser.CreateStr {
		log.Fatal("please check sql file statement is DDL and action is create  ")
	}
	tableName := ddl.NewName.Name.String()

	gotableName := snaker.SnakeToCamelIdentifier(tableName)
	mytable := &Table{
		Database:    db,
		TableName:   tableName,
		GoTableName: gotableName,
		PackageName: strings.ToLower(gotableName),
	}
	columns, err := MyTableColumnsFromDDL(ddl)
	if err != nil {
		log.Fatal(err)
	}
	if len(columns) <= 0 {
		log.Fatal("schema or table not exist")
	}
	mytable.Fields = columns
	for _, v := range columns {
		if v.IsPrimaryKey {
			mytable.PrimaryKey = v
		}
		if v.GoColumnType == "time.Time" {
			mytable.ImportTime = true
		}
	}
	if mytable.PrimaryKey == nil {
		log.Fatal("table do not have a primary key")
	}

	mytable.GenerateWhereCol = mytable.Fields
	return mytable

}

func MyTableColumnsFromDDL(ddl *sqlparser.DDL) ([]*Column, error) {

	res := []*Column{}
	for k, v := range ddl.TableSpec.Columns {

		var ct string
		if v.Type.Unsigned {
			ct = "unsigned"
		}
		var dct bool
		if v.Type.Default != nil {
			if string(v.Type.Default.Val) == "current_timestamp" || string(v.Type.Default.Val) == "current_timestamp()" {
				dct = true
			}
		}
		var comment string
		if v.Type.Comment != nil {
			comment = string(v.Type.Comment.Val)
		}

		c := &Column{
			OrdinalPosition:           k,
			ColumnName:                v.Name.String(),
			DataType:                  v.Type.Type,
			ColumnType:                ct,
			ColumnComment:             comment,
			NotNull:                   bool(v.Type.NotNull),
			IsPrimaryKey:              false,
			IsAutoIncrment:            bool(v.Type.Autoincrement),
			IsDefaultCurrentTimestamp: dct,
			GoColumnName:              "",
			GoColumnType:              "",
			BigType:                   0,
			GoConditionType:           "",
		}

		c.GoColumnName = snaker.SnakeToCamelIdentifier(c.ColumnName)
		c.GoColumnType, c.BigType = MysqlToGoFieldType(c.DataType, c.ColumnType)
		if strings.Contains(c.GoColumnType, "int") {
			c.GoColumnType = "int64"
		}
		c.GoConditionType = c.GoColumnType
		if c.BigType == bigtypeCompareTime {
			c.GoConditionType = "string"
		}
		res = append(res, c)

	}
	var primaryKey string
	for _, v := range ddl.TableSpec.Indexes {
		if v.Info.Primary {
			if len(v.Columns) != 1 {
				log.Fatal("primary key must be one column")
			}
			primaryKey = v.Columns[0].Column.String()
		}
	}
	for _, v := range res {
		if v.ColumnName == primaryKey {
			v.IsPrimaryKey = true
		}
	}
	return res, nil
}
