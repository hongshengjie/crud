package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var ctx = context.Background()

type User struct {
	Id    int64     `json:"id"`    // id字段
	Name  string    `json:"name"`  // 名称
	Age   int64     `json:"age"`   // 年龄
	Ctime time.Time `json:"ctime"` // 创建时间
	Mtime time.Time `json:"mtime"` // 更新时间
}

const dsn = "root:123456@tcp(127.0.0.1:3306)/test?parseTime=true"

func main() {
	db, _ = sql.Open("mysql", dsn)
	fmt.Println(FindUserReflect())
}

func FindUser() ([]*User, error) {
	rows, err := db.QueryContext(ctx, "SELECT `id`,`name`,`age`,`ctime`,`mtime` FROM user WHERE `age`<? ORDER BY `id` DESC LIMIT 20 ", 20)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []*User{}
	for rows.Next() {
		a := &User{}
		if err := rows.Scan(&a.Id, &a.Name, &a.Age, &a.Ctime, &a.Mtime); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return result, nil
}

func FindUserReflect() ([]*User, error) {
	b := SelectBuilder{builder: &strings.Builder{}}
	sql, args := b.
		Select("id", "name", "age", "ctime", "mtime").
		From("user").
		Where(GT("id", 0), GT("age", 0)).
		OrderBy("id").
		Limit(0, 20).
		Query()
	fmt.Println(sql, args)
	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	result := []*User{}
	err = ScanSlice(rows, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type SelectBuilder struct {
	builder   *strings.Builder
	column    []string
	tabelName string
	where     []func(s *SelectBuilder)
	args      []interface{}
	orderby   string
	offset    *int64
	limit     *int64
}

func (s *SelectBuilder) Select(field ...string) *SelectBuilder {
	s.column = append(s.column, field...)
	return s
}
func GT(field string, arg interface{}) func(s *SelectBuilder) {
	return func(s *SelectBuilder) {
		s.builder.WriteString("`" + field + "`" + " > ?")
		s.args = append(s.args, arg)
	}
}
func (s *SelectBuilder) From(name string) *SelectBuilder {
	s.tabelName = name
	return s
}
func (s *SelectBuilder) Where(f ...func(s *SelectBuilder)) *SelectBuilder {
	s.where = append(s.where, f...)
	return s
}
func (s *SelectBuilder) OrderBy(field string) *SelectBuilder {
	s.orderby = field
	return s
}
func (s *SelectBuilder) Limit(offset, limit int64) *SelectBuilder {
	s.offset = &offset
	s.limit = &limit
	return s
}
func (s *SelectBuilder) Query() (string, []interface{}) {
	s.builder.WriteString("SELECT ")
	for k, v := range s.column {
		if k > 0 {
			s.builder.WriteString(",")
		}
		s.builder.WriteString("`" + v + "`")
	}
	s.builder.WriteString(" FROM ")
	s.builder.WriteString("`" + s.tabelName + "` ")
	if len(s.where) > 0 {
		s.builder.WriteString("WHERE ")
		for k, f := range s.where {
			if k > 0 {
				s.builder.WriteString(" AND ")
			}
			f(s)
		}
	}
	if s.orderby != "" {
		s.builder.WriteString(" ORDER BY " + s.orderby)
	}
	if s.limit != nil {
		s.builder.WriteString(" LIMIT ")
		s.builder.WriteString(strconv.FormatInt(*s.limit, 10))
	}
	if s.offset != nil {
		s.builder.WriteString(" OFFSET ")
		s.builder.WriteString(strconv.FormatInt(*s.offset, 10))
	}
	return s.builder.String(), s.args
}

func ScanSlice(rows *sql.Rows, dst interface{}) error {
	defer rows.Close()
	// dst的地址
	val := reflect.ValueOf(dst) //  &[]*User
	// 判断是否是指针类型，go是值传递，只有传指针才能把更改生效
	if val.Kind() != reflect.Ptr {
		return errors.New("dst not a pointer")
	}
	// 指针指向的 Value
	val = reflect.Indirect(val) // []*User
	if val.Kind() != reflect.Slice {
		return errors.New("dst not a pointer to slice")
	}
	// 获取slice中的类型
	struPointer := val.Type().Elem() //*main.User

	// 指针指向的类型 具体结构体
	stru := struPointer.Elem() //     main.User
	//

	cols, err := rows.Columns() // [id,name,age,ctime,mtime]
	if err != nil {
		return err
	}
	// 判断查询的字段数是否大于 结构体的字段数
	if stru.NumField() < len(cols) { // 5,5
		return fmt.Errorf("cols num not match")
	}

	//结构体的json tag的value 对应 字段index
	tagIdx := make(map[string]int) //map tag -> field idx
	for i := 0; i < stru.NumField(); i++ {
		tagname := stru.Field(i).Tag.Get("json")
		if tagname != "" {
			tagIdx[tagname] = i
		}
	}
	resultType := make([]reflect.Type, 0, len(cols)) // [int64,string,int64,time.Time,time.Time]
	index := make([]int, 0, len(cols))               // [0,1,2,3,4,5]
	// 查找和列名相对应的结构体jsontag name 的字段类型，保存类型和序号 到resultType 和 index 中
	for _, v := range cols {
		if i, ok := tagIdx[v]; ok {
			resultType = append(resultType, stru.Field(i).Type)
			index = append(index, i)
		}
	}
	for rows.Next() {
		// 创建结构体指针
		obj := reflect.New(stru).Elem()                   // main.User
		result := make([]interface{}, 0, len(resultType)) //[]
		// 创建结构体字段类型实例的指针,并转化为interface{} 类型
		for _, v := range resultType {
			result = append(result, reflect.New(v).Interface()) // *Int64 ,*string ....
		}
		// 扫描结果
		err := rows.Scan(result...)
		if err != nil {
			panic(err)
		}
		fmt.Println(result...)

		for i, v := range result {
			// 找对对应的结构体index
			filedIndex := index[i]
			// 把scan 后的值通过反射得到指针指向的value，赋值给对应的结构体字段
			obj.Field(filedIndex).Set(reflect.ValueOf(v).Elem()) // 给obj 的每个字段赋值
		}
		// append 到slice
		vv := reflect.Append(val, obj.Addr()) // append到 []*main.User
		val.Set(vv)                           // &[]*main.User
	}
	return rows.Err()
}
