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
	"github.com/hongshengjie/crud/internal/model"
)

//go:embed "internal/templates/model.tmpl"
var modelTmpl []byte

//go:embed "internal/templates/builder.tmpl"
var crudTmpl []byte

//go:embed "internal/templates/where.tmpl"
var whereTmpl []byte

//go:embed "internal/templates/proto.tmpl"
var protoTmpl []byte

//go:embed "internal/templates/service.tmpl"
var serviceTmpl []byte

//go:embed "internal/templates/client.tmpl"
var clientTmpl []byte

var database string
var path string
var service bool
var protopkg string

//var fields string
const defaultDir = "crud"

func init() {
	flag.StringVar(&path, "path", "", ".sql file path or folder")
	flag.BoolVar(&service, "service", false, "-service  generate GRPC proto message and service implementation")
	flag.StringVar(&protopkg, "protopkg", "", "-protopkg  proto package field value")
}

func main() {

	flag.Parse()

	if len(os.Args) == 1 {
		info, err := os.Stat(defaultDir)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatal("crud dir is not exist please exec: crud init")
				return
			}
			log.Fatal(err)
			return
		}
		if info.IsDir() {
			path = defaultDir
		}
	}
	// subcommand
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "init":
			//create crud dir
			if err := os.Mkdir(defaultDir, os.ModePerm); err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	tableObjs, isDir := tableFromSql(path)
	for _, v := range tableObjs {
		generateFiles(v)
	}
	if isDir && path == defaultDir {
		generateFile(filepath.Join(defaultDir, "aa_client.go"), string(clientTmpl), f, tableObjs)
	}

}

func tableFromSql(path string) (tableObjs []*model.Table, isDir bool) {
	relativePath := model.GetRelativePath()
	info, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	if info.IsDir() {
		isDir = true
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
	return tableObjs, isDir
}

var f = template.FuncMap{
	"sqltool":  model.SQLTool,
	"isnumber": model.IsNumber,
	"Incr":     model.Incr,
}

func generateFiles(tableObj *model.Table) {

	//创建目录
	dir := filepath.Join(defaultDir, tableObj.PackageName)
	os.Mkdir(dir, os.ModePerm)
	generateFile(filepath.Join(dir, "model.go"), string(modelTmpl), f, tableObj)
	generateFile(filepath.Join(dir, "where.go"), string(whereTmpl), f, tableObj)
	generateFile(filepath.Join(dir, "builder.go"), string(crudTmpl), f, tableObj)
	if service {
		generateService(tableObj)
	}

}
func generateService(tableObj *model.Table) {
	pkgName := tableObj.PackageName
	tableObj.Protopkg = protopkg
	os.Mkdir(filepath.Join("proto"), os.ModePerm)
	os.Mkdir(filepath.Join("service"), os.ModePerm)

	generateFile(filepath.Join("proto", pkgName+".api.proto"), string(protoTmpl), f, tableObj)
	//protoc --go_out=. --go-grpc_out=.  user.api.proto
	cmd := exec.Command("protoc", "-I.", "-I/usr/local/include", "--go_out=.", "--go-grpc_out=.", filepath.Join("proto", pkgName+".api.proto"))
	cmd.Dir = filepath.Join(model.GetCurrentPath())
	log.Println(cmd.Dir, "exec:", cmd.String())
	s, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(s), err)
	}

	generateFile(filepath.Join("service", pkgName+".service.go"), string(serviceTmpl), f, tableObj)
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
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	file.Write(result)
	file.Close()
}
