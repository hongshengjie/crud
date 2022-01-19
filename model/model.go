package model

// Table Table
type Table struct {
	Database         string
	TableName        string    // table name
	GoTableName      string    // go struct name
	PackageName      string    // package name
	Fields           []*Column // columns
	GenerateWhereCol []*Column // GenerateWhereCol 生成where字段比较方法的列
	PrimaryKey       *Column   // priomary_key column
	ImportTime       bool      // is need import time
	RelativePath     string
}

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
	ProtoType                 string // protoType
}
