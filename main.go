package main

import (
	"bytes"
	"database/sql"
	"flag"
	"io/ioutil"
	"path/filepath"

	"go/format"
	"log"
	"os"
	"strings"
	"text/template"

	_ "embed"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hongshengjie/crud/mytable"
)

//go:embed "templates/model.tmpl"
var modelTmpl []byte

//go:embed "templates/builder.tmpl"
var crudTmpl []byte

//go:embed "templates/where.tmpl"
var whereTmpl []byte

var database string
var dsn string
var table string
var path string

//var fields string

func init() {
	flag.StringVar(&database, "database", "mysql", "mysql or postgres")
	flag.StringVar(&dsn, "dsn", "", "mysql connection url")
	flag.StringVar(&table, "table", "", "table name")
	flag.StringVar(&path, "path", "", ".sql file path or dir generate code from DDL sql file")
}

func main() {

	flag.Parse()
	switch database {
	case "mysql":
	case "postgres":
	default:
		log.Fatalln("database not right")
	}

	var tableObjs []*mytable.Table
	if path != "" {
		tableObjs = append(tableObjs, tableFromSql(path)...)
	} else {
		tableObjs = append(tableObjs, tableFromDB()...)
	}
	for _, v := range tableObjs {
		generateFiles(v)
	}
}

func tableFromSql(path string) (tableObjs []*mytable.Table) {
	info, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	if info.IsDir() {
		fs, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range fs {
			if !v.IsDir() && strings.HasSuffix(strings.ToLower(v.Name()), ".sql") {
				obj := mytable.MytableFromSqlFile(database, filepath.Join(path, v.Name()))
				if obj != nil {
					tableObjs = append(tableObjs, obj)
				}

			}

		}
	} else {
		tableObjs = append(tableObjs, mytable.MytableFromSqlFile(database, path))
	}
	return tableObjs
}

func tableFromDB() (tableObjs []*mytable.Table) {
	if dsn == "" || table == "" {
		log.Fatalln("dns or schema or table is empty")
	}

	temps := strings.Split(dsn, "/")
	if len(temps) < 2 {
		log.Fatalln("dsn not hava /")
	}
	temps2 := strings.Split(temps[1], "?")
	if len(temps2) < 2 {
		log.Fatalln("dsn not hava ?")
	}
	schema := temps2[0]

	db, err := sql.Open(database, dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	tbs := strings.Split(table, ",")
	for _, v := range tbs {
		t := mytable.NewTable(db, database, schema, v, []string{}, true)
		tableObjs = append(tableObjs, t)
	}
	return tableObjs

}

func generateFiles(tableObj *mytable.Table) {
	f := template.FuncMap{
		"sqltool":  mytable.SQLTool,
		"isnumber": mytable.IsNumber,
	}
	//????????????
	os.Mkdir(tableObj.PackageName, os.ModePerm)
	generateFile("model", string(modelTmpl), f, tableObj)
	generateFile("where", string(whereTmpl), f, tableObj)
	generateFile("builder", string(crudTmpl), f, tableObj)

}

func generateFile(name, tmpl string, f template.FuncMap, table *mytable.Table) {
	tpl, err := template.New(name).Funcs(f).Parse(string(tmpl))
	if err != nil {
		log.Fatalln(err)
	}
	bs := bytes.NewBuffer(nil)
	err = tpl.Execute(bs, table)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := format.Source(bs.Bytes())
	if err != nil {
		log.Fatalln(err)
	}
	//?????????
	fileName := filepath.Join(table.PackageName, name+".go")
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		log.Fatalln(err)
	}
	file.Write(result)
	file.Close()
}
