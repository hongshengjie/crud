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

//go:embed "internal/templates/react-grommet.tmpl"
var reactGrommetTmpl []byte

//go:embed "internal/templates/builder_mgo.tmpl"
var crudMgo []byte

//go:embed "internal/templates/struct2pb.tmpl"
var struct2PB []byte

var database string
var path string
var service bool
var http bool
var protopkg string
var reactgrommet bool
var mgo string
var struct2pb string

// var fields string
const defaultDir = "crud"

func init() {
	//flag.StringVar(&path, "path", "cr", ".sql file path or folder")
	flag.BoolVar(&service, "service", false, "-service  generate GRPC proto message and service implementation")
	flag.BoolVar(&http, "http", false, "-http  generate Gin controller")
	flag.BoolVar(&reactgrommet, "reactgrommet", false, "-reactgrommet  generate reactgrommet tsx code work with -service")
	flag.StringVar(&protopkg, "protopkg", "", "-protopkg  proto package field value")
	flag.StringVar(&mgo, "mgo", "", "-mgo find struct from file and generate crud method example  ./user.go:User  User struct in ./user.go file ")
	flag.StringVar(&struct2pb, "struct2pb", "", "-struct2pb find struct from file and generate corresponding proto message  ./user.go:User  User struct in ./user.go file ")
}

func main() {

	flag.Parse()

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

	if mgo != "" {
		pathName := strings.Split(mgo, ":")
		if len(pathName) != 2 {
			log.Fatalf("-mgo not right example ./user.go:User")
		}
		filePath := pathName[0]
		structName := pathName[1]
		doc := model.ParseMongoStruct(filePath, structName)
		generateFile(filePath, string(crudMgo), nil, doc)
		return
	}
	if struct2pb != "" {
		pathName := strings.Split(struct2pb, ":")
		if len(pathName) != 2 {
			log.Fatalf("-struct2pb not right example ./user.go:User")
		}
		filePath := pathName[0]
		structName := pathName[1]
		message := model.ParseStruct(filePath, structName)
		tpl, err := template.New("").Funcs(f).Parse(string(struct2PB))
		if err != nil {
			log.Fatalln(err)
		}
		err = tpl.Execute(os.Stdout, message)
		if err != nil {
			log.Fatalln(err)
		}

		return
	}
	if path == "" {
		path = defaultDir
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
	"sqltool":                        model.SQLTool,
	"isnumber":                       model.IsNumber,
	"Incr":                           model.Incr,
	"GoTypeToTypeScriptDefaultValue": model.GoTypeToTypeScriptDefaultValue,
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

	// proto-go  grpc
	var cmd *exec.Cmd
	if http {
		cmd = exec.Command("protoc", "-I.", "--go_out=.", "--go-grpc_out=.", "--go-gin_out=.", filepath.Join("proto", pkgName+".api.proto"))
	} else {
		cmd = exec.Command("protoc", "-I.", "--go_out=.", "--go-grpc_out=.", filepath.Join("proto", pkgName+".api.proto"))
	}

	cmd.Dir = filepath.Join(model.GetCurrentPath())
	log.Println(cmd.Dir, "exec:", cmd.String())
	s, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(s), err)
	}
	// inject-tag
	cmd = exec.Command("protoc-go-inject-tag", "-input", filepath.Join("api", pkgName+".api.pb.go"))
	cmd.Dir = filepath.Join(model.GetCurrentPath())
	log.Println(cmd.Dir, "exec:", cmd.String())
	s, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(string(s), err)
	}

	generateFile(filepath.Join("service", pkgName+".service.go"), string(serviceTmpl), f, tableObj)

	if reactgrommet {
		generateFile(filepath.Join("web", "src", "pages", pkgName+".tsx"), string(reactGrommetTmpl), f, tableObj)
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
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	file.Write(result)
	file.Close()
}
