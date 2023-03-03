package model

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"strings"
)

type PbMessage struct {
	PbName string
	GoName string
	Fields []*PbField
}

type PbField struct {
	PbName string
	PbType string
	GoName string
	GoType string
}

func ParseStruct(filePath, structName string) *PbMessage {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	ts, ok := f.Scope.Objects[structName].Decl.(*ast.TypeSpec)
	if !ok {
		log.Fatalf("not find struct : %s on %s go source file", structName, filePath)
	}

	d := &PbMessage{
		PbName: structName,
		GoName: structName,
	}
	st, ok := ts.Type.(*ast.StructType)
	if !ok {
		log.Fatalf("%s is not a struct ", structName)
	}

	for _, field := range st.Fields.List {
		gotype := GetGoType(field.Type)
		f := &PbField{
			PbName: strings.ToLower(field.Names[0].Name),
			PbType: GoTypeToProtoType(gotype),
			GoName: field.Names[0].Name,
			GoType: gotype,
		}
		if f.PbType == "" {
			f.PbType = "string"
		}
		d.Fields = append(d.Fields, f)

	}
	return d

}

func GetGoType(exp ast.Expr) string {
	var gotype string
	switch reflect.TypeOf(exp) {
	case reflect.TypeOf(&ast.SelectorExpr{}):
		vv := exp.(*ast.SelectorExpr)
		pkg := vv.X.(*ast.Ident)
		gotype = pkg.String() + "." + vv.Sel.String()
	case reflect.TypeOf(&ast.Ident{}):
		vv := exp.(*ast.Ident)
		gotype = vv.String()
	case reflect.TypeOf(&ast.ArrayType{}):
		vv := exp.(*ast.ArrayType)
		gotype = "[]" + GetGoType(vv.Elt)

	case reflect.TypeOf(&ast.MapType{}):
		vv := exp.(*ast.MapType)
		key := GetGoType(vv.Key)
		value := GetGoType(vv.Value)
		gotype = fmt.Sprintf("map[%s]%s", key, value)
	default:
		panic("not support embed field or include other struct ")
	}
	return gotype
}
