package mytable

import (
	"strings"
)

const (
	bigtypeCompare       = 1
	bigtypeCompareString = 2
	bigtypeCompareTime   = 3
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
		typ = "[]uint8"
	case "bool", "boolean":
		typ = "bool"
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext", "json":
		typ = "string"
		gtp = bigtypeCompareString
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
	case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		typ = "[]byte"
		gtp = bigtypeCompare
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
		case "set":
			ns = append(ns, v.ColumnName+" = ? ")
		default:
			ns = append(ns, flag)
		}

	}
	return strings.Join(ns, ",")
}

// SQLUpdate sql update
func SQLUpdate(t *Table, omit bool) string {
	var ns []string
	var prime *Column
	for _, v := range t.Fields {
		if omit {
			if v.IsDefaultCurrentTimestamp {
				continue
			}
		}
		if v.IsPrimaryKey {
			prime = v
			continue
		}
		ns = append(ns, "a."+v.GoColumnName)
	}
	ns = append(ns, "a."+prime.GoColumnName)

	return strings.Join(ns, ",")
}

// SQLIndexParamList SQLIndexParamList
func SQLIndexParamList(index *Index, needType bool) string {
	var ns []string
	for _, v := range index.IndexColumns {
		s := v.ColumnName
		if needType {
			s = s + " " + v.GoColumnType
		}
		ns = append(ns, s)
	}
	return strings.Join(ns, ",")
}

// SQLIndexQuery SQLIndexQuery
func SQLIndexQuery(index *Index) string {
	var ns []string
	for _, v := range index.IndexColumns {
		s := v.ColumnName + " = ? "

		ns = append(ns, s)
	}
	return strings.Join(ns, "AND ")
}
