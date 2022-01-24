package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hongshengjie/crud/example/user"
	"github.com/hongshengjie/crud/example/user/api"
	"github.com/hongshengjie/crud/example/user/service"
	"github.com/hongshengjie/crud/xsql"

	_ "github.com/go-sql-driver/mysql"
)

//go:generate protoc --go_out=. --go-grpc_out=.  alltypetable.api.proto
//go:generate protoc --go_out=. --go-grpc_out=.  user.api.proto
//linux   protoc  -I . -I /usr/local/include --go_out=. --go-grpc_out=.  user.api.proto
var db *sql.DB
var dsn = "root:123456@tcp(127.0.0.1:3306)/test?parseTime=true"
var ctx = context.Background()

func InitDB() {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
}

func UserExample() {

	u := &user.User{
		Id:    0,
		Name:  "shengjie",
		Age:   18,
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	u1 := &user.User{
		Id:    1,
		Name:  "shengjie2",
		Age:   22,
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	_, err := user.
		Create(xsql.Debug(db)).
		SetUser(u).
		Save(ctx)

	fmt.Println(err)

	_, err = user.
		Create(db).
		SetUser(u1, u).
		Save(ctx)

	fmt.Println(err)

	_, err = user.
		Create(db).
		SetUser(u1, u).
		Upsert(ctx)

	fmt.Println(err)

	au, err := user.
		Find(db).
		Where(user.Or(
			user.IdGT(10),
			user.NameEQ("testb"),
		)).
		Offset(0).
		Limit(3).
		OrderAsc("name").
		All(ctx)

	fmt.Printf("%+v %v", au, err)

	c, err := user.
		Find(db).
		Count(xsql.Distinct(user.Name)).
		Where(user.Or(
			user.IdGT(10),
			user.NameEQ("testb"),
		)).
		Int64(ctx)

	fmt.Println(c, err)

	c1, err := user.
		Find(db).
		Select(user.Columns()...).
		Where(user.Or(
			user.IdGT(10),
			user.NameEQ("testb"),
		)).
		All(ctx)

	fmt.Println(c1, err)

	c2, err := user.
		Find(db).
		Select(xsql.Sum(user.Age)).
		Where(user.Or(
			user.IdGT(10),
			user.NameEQ("testb"),
		)).
		Int64(ctx)

	fmt.Println(c2, err)

	effect, err := user.
		Update(db).
		SetAge(100).
		SetName("java").
		Where(user.IdEQ(1)).
		Save(ctx)
	fmt.Println(effect, err)

	effect, err = user.
		Update(db).
		AddAge(100).
		SetName("java").
		Where(user.IdEQ(5)).
		Save(ctx)

	fmt.Println(effect, err)

	effect, err = user.
		Delete(db).
		Where(user.And(
			user.IdEQ(3),
			user.IdIn(1, 3),
		)).
		Exec(ctx)

	fmt.Println(effect, err)

	effect, err = user.
		Delete(db).
		Where(user.IdEQ(2)).
		Exec(ctx)

	fmt.Println(effect, err)
	tx, _ := db.Begin()
	u2 := &user.User{
		Id:    0,
		Name:  "foo",
		Age:   2,
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	_, err = user.
		Create(tx).
		SetUser(u2).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return
	}

	effect, err = user.
		Update(tx).
		SetAge(100).
		Where(user.IdEQ(1)).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	fmt.Println(effect, err)

}
func UserSelect() {
	us, _ := user.Find(db).
		Select().
		Where(
			user.AgeGT(10),
		).
		All(ctx)

	us2, _ := user.Find(db).
		Select(user.Columns()...).
		Where(
			user.AgeGT(10),
		).
		All(ctx)
	fmt.Println(us, us2)
}

func main() {

	InitDB()
	ListUsers()
}

func ListUsers() {
	s := service.UserServiceImpl{}
	s.SetDB(xsql.Debug(db))
	r, err := s.ListUsers(ctx, &api.ListUsersReq{
		Page:     1,
		PageSize: 20,
		Orderby:  "-id",
		Filter: []*api.UserFilter{
			{
				Field: "name",
				Op:    "in",
				Value: "java,shengjie",
			},
		},
	})
	rr, _ := json.Marshal(r)
	fmt.Println(string(rr), err)

}
