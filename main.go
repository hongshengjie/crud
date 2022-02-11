package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"go/format"
	"log"
	"os"
	"strings"
	"text/template"

	_ "embed"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hongshengjie/crud/model"
)

//go:embed "templates/model.tmpl"
var modelTmpl []byte

//go:embed "templates/builder.tmpl"
var crudTmpl []byte

//go:embed "templates/where.tmpl"
var whereTmpl []byte

//go:embed "templates/proto.tmpl"
var protoTmpl []byte

//go:embed "templates/service.tmpl"
var serviceTmpl []byte

//go:embed "templates/client.tmpl"
var clientTmpl []byte

var database string
var path string
var service bool
var protopkg string

//var fields string

func init() {
	flag.StringVar(&path, "path", "", ".sql file path or folder")
	flag.BoolVar(&service, "service", false, "-service  generate GRPC proto message and service implementation")
	flag.StringVar(&protopkg, "protopkg", "", "-protopkg  proto package field value")
}

func main() {

	flag.Parse()
	var tableObjs []*model.Table

	tableObjs = append(tableObjs, tableFromSql(path)...)

	for _, v := range tableObjs {
		generateFiles(v)
	}
	if len(tableObjs) > 1 {
		generateFile("client.go", string(clientTmpl), f, tableObjs)
	}

}

func tableFromSql(path string) (tableObjs []*model.Table) {
	relativePath := model.GetRelativePath()
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
				obj := model.MysqlTable(database, filepath.Join(path, v.Name()), relativePath)
				if obj != nil {
					tableObjs = append(tableObjs, obj)
				}

			}

		}
	} else {
		tableObjs = append(tableObjs, model.MysqlTable(database, path, relativePath))
	}
	return tableObjs
}

var f = template.FuncMap{
	"sqltool":  model.SQLTool,
	"isnumber": model.IsNumber,
	"Incr":     model.Incr,
}

func generateFiles(tableObj *model.Table) {

	pkgName := tableObj.PackageName
	//创建目录
	os.Mkdir(tableObj.PackageName, os.ModePerm)
	generateFile(filepath.Join(pkgName, "model.go"), string(modelTmpl), f, tableObj)
	generateFile(filepath.Join(pkgName, "where.go"), string(whereTmpl), f, tableObj)
	generateFile(filepath.Join(pkgName, "builder.go"), string(crudTmpl), f, tableObj)
	if service {
		tableObj.Protopkg = protopkg
		os.Mkdir(filepath.Join(tableObj.PackageName, "api"), os.ModePerm)
		generateFile(filepath.Join(pkgName, pkgName+".api.proto"), string(protoTmpl), f, tableObj)
		//protoc --go_out=. --go-grpc_out=.  user.api.proto
		cmd := exec.Command("protoc", "-I.", "-I/usr/local/include", "--go_out=.", "--go-grpc_out=.", pkgName+".api.proto")
		cmd.Dir = filepath.Join(model.GetCurrentPath(), pkgName)
		log.Println(cmd.Dir, "exec:", cmd.String())
		s, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(string(s), err)
		}
		os.Mkdir(filepath.Join(pkgName, "service"), os.ModePerm)
		generateFile(filepath.Join(pkgName, "service", pkgName+".service.go"), string(serviceTmpl), f, tableObj)
	}

}

func generateFile(filename, tmpl string, f template.FuncMap, data interface{}) {
	tpl, err := template.New(filename).Funcs(f).Parse(string(tmpl))
	if err != nil {
		log.Fatalln(err)
	}
	bs := bytes.NewBuffer(nil)
	err = tpl.Execute(bs, data)
	if err != nil {
		log.Fatalln(err)
	}

	result := bs.Bytes()
	if strings.HasSuffix(filename, ".go") {
		result, err = format.Source(bs.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		log.Fatalln(err)
	}
	file.Write(result)
	file.Close()
}
