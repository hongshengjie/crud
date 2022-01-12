package model

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/xwb1989/sqlparser"
)

func MysqlColumn(ddl *sqlparser.DDL) ([]*Column, error) {

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

		c.GoColumnName = GoCamelCase(c.ColumnName)
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
		v.ProtoType = GoTypeToProtoType(v.GoColumnType)
	}
	return res, nil
}

func MysqlTable(db string, path string) *Table {
	sql, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := sqlparser.ParseStrictDDL(trimTimeStampFunc(string(sql)))
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

	gotableName := GoCamelCase(tableName)
	mytable := &Table{
		Database:    db,
		TableName:   tableName,
		GoTableName: gotableName,
		PackageName: strings.ToLower(gotableName),
	}
	columns, err := MysqlColumn(ddl)
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
		v.ProtoType = GoTypeToProtoType(v.GoColumnType)
	}

	mytable.GenerateWhereCol = mytable.Fields
	return mytable

}
