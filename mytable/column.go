package mytable

import (
	"database/sql"
	"strings"

	"github.com/hongshengjie/crud/snaker"
)

// Column Column
type Column struct {
	OrdinalPosition           int    // field_ordinal
	ColumnName                string // column_name
	DataType                  string // data_type
	ColumnType                string // column_type
	ColumnComment             string // column_comment,
	NotNull                   bool   // not_null
	IsPrimaryKey              bool   // is_primary_key
	IsAutoIncrment            bool   // is_auto_incrment
	IsDefaultCurrentTimestamp bool   // is_default_currenttimestamp
	GoColumnName              string // go field name
	GoColumnType              string // go field type
	BigType                   int    // 0 表示不生成where 1 表示比较类型 2表示比较类型+字符串 3表示比较类型，修改传入参数
	GoConditionType           string // 生成where 的类型参数
}

// MyTableColumns MyTableColumns
func MyTableColumns(db *sql.DB, schema string, table string) ([]*Column, error) {
	var err error

	const sqlstr = `SELECT ` +
		`ordinal_position,` +
		`column_name, ` +
		`data_type, ` +
		`column_type,` +
		`column_comment, ` +
		`IF(is_nullable = 'YES', false, true) AS not_null, ` +
		`IF(column_key = 'PRI', true, false) AS is_primary_key, ` +
		`IF(extra = 'auto_increment',true,false) AS is_auto_incrment,` +
		`IF (column_default = 'CURRENT_TIMESTAMP' or column_default = 'current_timestamp()',true,false) AS is_default_currenttimestamp ` +
		`FROM information_schema.columns ` +
		`WHERE table_schema = ? AND table_name = ? ` +
		`ORDER BY ordinal_position`

	q, err := db.Query(sqlstr, schema, table)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	res := []*Column{}
	for q.Next() {
		c := Column{}
		err = q.Scan(&c.OrdinalPosition, &c.ColumnName, &c.DataType, &c.ColumnType, &c.ColumnComment, &c.NotNull, &c.IsPrimaryKey, &c.IsAutoIncrment, &c.IsDefaultCurrentTimestamp)
		if err != nil {
			return nil, err
		}
		c.GoColumnName = snaker.SnakeToCamelIdentifier(c.ColumnName)
		c.GoColumnType, c.BigType = MysqlToGoFieldType(c.DataType, c.ColumnType)
		c.GoConditionType = c.GoColumnType
		if c.BigType == bigtypeCompareTime {
			c.GoConditionType = "string"
		}
		if strings.Contains(c.GoColumnType, "int") {
			c.GoColumnType = "int64"
		}
		res = append(res, &c)
	}
	if q.Err() != nil {
		return nil, q.Err()
	}

	return res, nil
}
