package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"example/api"
	"example/crud"
	"example/crud/user"
	"example/service"

	"github.com/hongshengjie/crud/xsql"
)

//go:generate protoc --go_out=. --go-grpc_out=.  alltypetable.api.proto
//go:generate protoc --go_out=. --go-grpc_out=.  user.api.proto
//linux   protoc  -I . -I /usr/local/include --go_out=. --go-grpc_out=.  user.api.proto
var db *sql.DB

var ctx = context.Background()

func InitDB() {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

}

var client *crud.Client

var dsn = "root:123456@tcp(127.0.0.1:3306)/test?parseTime=true"

func InitDB2() {
	client, _ = crud.NewClient(&xsql.Config{
		DSN:          dsn,
		ReadDSN:      []string{dsn},
		Active:       10,
		Idle:         10,
		IdleTimeout:  time.Hour,
		QueryTimeout: time.Second,
		ExecTimeout:  time.Second,
	})
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

	us2, err := user.Find(db).
		Select(user.Columns()...).
		Where(
			user.AgeGT(10),
		).
		All(ctx)
	fmt.Println(us, us2, err)
}

func main() {

	//InitDB()
	InitDB2()
	//e, err := user.Update(db).SetAge(10).Where(user.IdEQ(4)).WithTimeOut(time.Millisecond * 250).Save(ctx)

	e, err := user.Find(db).Timeout(time.Millisecond * 30).All(ctx)
	fmt.Println(e, err)
	client.User.Find().Select().All(ctx)

	tx, _ := client.Begin(ctx)
	tx.User.Update().SetAge(1).Save(ctx)
	tx.Commit()
}

func ListUsers() {
	s := service.UserServiceImpl{Client: client}

	r, err := s.ListUsers(ctx, &api.ListUsersReq{
		Page:     1,
		PageSize: 20,
		OrderBy:  "-id",
		IdGt:     1,
	})
	rr, _ := json.Marshal(r)
	fmt.Println(string(rr), err)

}
