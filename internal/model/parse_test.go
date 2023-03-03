package model

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"text/template"
)

const tmp = `
message {{.PbName}} {
	{{- range $index,$field:=.Fields}}
	{{$field.PbType}} {{$field.PbName}}   = {{Incr $index}} ;
	{{- end }}
}

`

var f = template.FuncMap{

	"Incr": Incr,
}

func TestParseStruct(t *testing.T) {
	//temp, _ := os.ReadFile("../templates/builder_mgo.tmpl")
	r, _ := template.New("").Funcs(f).Parse(string(tmp))

	m := ParseStruct("../mgo/test/user.go", "User")
	b, _ := json.Marshal(m)
	fmt.Println(string(b))

	r.Execute(os.Stdout, m)
}
