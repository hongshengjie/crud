package model

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"strings"
)

type Document struct {
	Package       string
	Name          string
	GoName        string
	ImportTime    bool
	Fields        []*Field
	ObjectIDField *Field
}

type Field struct {
	Name   string
	GoName string
	GoType string
	Tag    string
}

func ParseMongoStruct(filePath, structName string) *Document {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	ts, ok := f.Scope.Objects[structName].Decl.(*ast.TypeSpec)
	if !ok {
		log.Fatalf("not find struct : %s on %s go source file", structName, filePath)
	}

	d := &Document{
		Package: f.Name.String(),
		Name:    strings.ToLower(structName),
		GoName:  structName,
	}
	st, ok := ts.Type.(*ast.StructType)
	if !ok {
		log.Fatalf("%s is not a struct ", structName)
	}

	for _, field := range st.Fields.List {
		f := &Field{
			GoName: field.Names[0].Name,
		}
		f.GoType = GetGoType(field.Type)
		if f.GoType == "time.Time" {
			d.ImportTime = true
		}
		if field.Tag != nil {
			f.Tag = field.Tag.Value

			v1, ok := reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Lookup("bson")
			if ok {
				if v1 == "_id,omitempty" {
					f.Name = "_id"
					d.ObjectIDField = f
				} else {
					f.Name = strings.Split(v1, ",")[0]
				}
			}

		} else {
			f.Name = strings.ToLower(f.GoName)
		}
		d.Fields = append(d.Fields, f)

	}
	if d.ObjectIDField == nil {
		log.Fatalf("struct must have ID field with bson tag  `bson:\"_id,omitempty\"` ")
	}

	return d

}
