package mytable

import (
	"database/sql"

	"log"
	"sort"
	"strings"

	"github.com/hongshengjie/crud/snaker"
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
