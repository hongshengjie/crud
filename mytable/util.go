package mytable

import (
	"strings"
)

const (
	bigtypeCompare       = 1
	bigtypeCompareString = 2
	bigtypeCompareTime   = 3
	bigtypeCompareBit    = 4
)

// MysqlToGoFieldType MysqlToGoFieldType
func MysqlToGoFieldType(dt, ct string) (string, int) {
	var unsigned bool
	if strings.Contains(ct, "unsigned") {
		unsigned = true
	}
	var typ string
	var gtp int
	switch dt {
	case "bit":
		typ = "[]byte"
		gtp = bigtypeCompareBit
	case "bool", "boolean":
		typ = "bool"
	case "char", "varchar":
		typ = "string"
		gtp = bigtypeCompareString
	case "tinytext", "text", "mediumtext", "longtext", "json":
		typ = "string"
	case "tinyint":
		typ = "int8"
		if unsigned {
			typ = "uint8"
		}
		gtp = bigtypeCompare
	case "smallint":
		typ = "int16"
		if unsigned {
			typ = "uint16"
		}
		gtp = bigtypeCompare
	case "mediumint", "int", "integer":
		typ = "int32"
		if unsigned {
			typ = "uint32"
		}
		gtp = bigtypeCompare
	case "bigint":
		typ = "int64"
		if unsigned {
			typ = "uint64"
		}
		gtp = bigtypeCompare
	case "float":
		typ = "float32"
		gtp = bigtypeCompare
	case "decimal", "double":
		typ = "float64"
		gtp = bigtypeCompare
	case "binary", "varbinary":
		typ = "[]byte"
		gtp = bigtypeCompare
	case "tinyblob", "blob", "mediumblob", "longblob":
		typ = "[]byte"
	case "timestamp", "datetime", "date":
		typ = "time.Time"
		gtp = bigtypeCompareTime
	case "time", "year", "enum", "set":
		typ = "string"
		gtp = bigtypeCompare
	default:
		typ = "UNKNOWN"
	}
	return typ, gtp
}

//SQLTool SQLTool
func SQLTool(t *Table, omit bool, flag string) string {
	var ns []string
	for _, v := range t.Fields {
		if omit {
			if v.IsAutoIncrment || v.IsDefaultCurrentTimestamp {
				continue
			}
		}
		switch flag {
		case "field":
			ns = append(ns, "`"+v.ColumnName+"`")
		case "?":
			ns = append(ns, "?")
		case "gofield":
			ns = append(ns, "&a."+v.GoColumnName)
		case "goinfield":
			ns = append(ns, "in.a."+v.GoColumnName)
		case "goinfieldcol":
			ns = append(ns, v.GoColumnName)
		case "goinfieldcolbulk":
			ns = append(ns, "a."+v.GoColumnName)
		case "set":
			ns = append(ns, v.ColumnName+" = ? ")
		default:
			ns = append(ns, flag)
		}

	}
	return strings.Join(ns, ",")
}

func IsNumber(arg string) bool {
	switch arg {
	case "int8", "int16", "int", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return true
	}
	return false
}
