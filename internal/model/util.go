package model

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
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
func GoTypeToProtoType(g string) string {
	switch g {
	case "[]byte":
		return "bytes"
	case "bool":
		return "bool"
	case "string":
		return "string"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "int8", "int16", "int32":
		return "int32"
	case "uint8", "uint16", "uint32":
		return "uint32"
	case "int", "int64":
		return "int64"
	case "uint64":
		return "uint64"
	case "time.Time":
		return "string"
	default:
		return ""
	}
}

func GoTypeToTypeScriptDefaultValue(g string) string {
	switch g {
	case "[]byte":
		return "new Uint8Array()"
	case "bool":
		return "false"
	case "string":
		return "''"
	case "float32":
		return "0"
	case "float64":
		return "0"
	case "int8", "int16", "int32":
		return "0"
	case "uint8", "uint16", "uint32":
		return "0"
	case "uint64", "int64", "int":
		return "0n"
	case "time.Time":
		return "''"
	default:
		return ""
	}
}

func Incr(x int) int {
	return x + 1
}

// SQLTool SQLTool
func SQLTool(t *Table, flag string) string {
	var ns []string
	for _, v := range t.Fields {
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

func trimTimeStampFunc(raw string) string {
	s := strings.ReplaceAll(raw, "current_timestamp()", "current_timestamp")
	s = strings.ReplaceAll(s, "CURRENT_TIMESTAMP()", "CURRENT_TIMESTAMP")
	return s
}

func GoCamelCase(s string) string {
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	var b []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '.' && i+1 < len(s) && isASCIILower(s[i+1]):
			// Skip over '.' in ".{{lowercase}}".
		case c == '.':
			b = append(b, '_') // convert '.' to '_'
		case c == '_' && (i == 0 || s[i-1] == '.'):
			// Convert initial '_' to ensure we start with a capital letter.
			// Do the same for '_' after '.' to match historic behavior.
			b = append(b, 'X') // convert '_' to 'X'
		case c == '_' && i+1 < len(s) && isASCIILower(s[i+1]):
			// Skip over '_' in "_{{lowercase}}".
		case isASCIIDigit(c):
			b = append(b, c)
		default:
			// Assume we have a letter now - if not, it's a bogus identifier.
			// The next word is a sequence of characters that must start upper case.
			if isASCIILower(c) {
				c -= 'a' - 'A' // convert lowercase to uppercase
			}
			b = append(b, c)

			// Accept lower case sequence that follows.
			for ; i+1 < len(s) && isASCIILower(s[i+1]); i++ {
				b = append(b, s[i+1])
			}
		}
	}
	return string(b)
}

// JSONCamelCase converts a snake_case identifier to a camelCase identifier,
// according to the protobuf JSON specification.
func JSONCamelCase(s string) string {
	var b []byte
	var wasUnderscore bool
	for i := 0; i < len(s); i++ { // proto identifiers are always ASCII
		c := s[i]
		if c != '_' {
			if wasUnderscore && isASCIILower(c) {
				c -= 'a' - 'A' // convert to uppercase
			}
			b = append(b, c)
		}
		wasUnderscore = c == '_'
	}
	return string(b)
}

// JSONSnakeCase converts a camelCase identifier to a snake_case identifier,
// according to the protobuf JSON specification.
func JSONSnakeCase(s string) string {
	var b []byte
	for i := 0; i < len(s); i++ { // proto identifiers are always ASCII
		c := s[i]
		if isASCIIUpper(c) {
			b = append(b, '_')
			c += 'a' - 'A' // convert to lowercase
		}
		b = append(b, c)
	}
	return string(b)
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}
func isASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func GetCurrentPath() string {
	exPath, _ := os.Getwd()
	return exPath
}

func GetRelativePath() string {
	modName, rootPath := GetModuleName()
	pwd := GetCurrentPath()
	relative := strings.TrimPrefix(pwd, rootPath)
	return filepath.Join(modName, relative)
}

func GetModuleName() (string, string) {
	mod := GoModFilePath()
	if mod == "" {
		return "", ""
	}
	f, err := ioutil.ReadFile(mod)
	if err != nil {
		log.Fatal(err)
	}
	p, _ := filepath.Split(mod)
	// module name and project root path
	return modfile.ModulePath(f), filepath.Clean(p)
}

func GoModFilePath() string {
	exPath := GetCurrentPath()
	gomodPath := []string{}
	names := strings.Split(exPath, string(os.PathSeparator))

	for k := range names {
		if k == 0 {
			if strings.HasSuffix(names[0], ":") {
				names[0] = names[0] + string(os.PathSeparator)
			} else {
				names[0] = string(os.PathSeparator) + names[0]
			}
		}
		prefix := filepath.Join(names[:k+1]...)
		gomodPath = append(gomodPath, filepath.Join(prefix, "go.mod"))

	}

	for i := len(gomodPath)/2 - 1; i >= 0; i-- {
		opp := len(gomodPath) - 1 - i
		gomodPath[i], gomodPath[opp] = gomodPath[opp], gomodPath[i]
	}
	for _, v := range gomodPath {
		if _, err := os.Stat(v); os.IsNotExist(err) {
			continue
		} else {
			return v
		}
	}
	return ""

}
