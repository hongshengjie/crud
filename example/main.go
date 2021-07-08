package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hongshengjie/crud/example/user"
	"github.com/hongshengjie/xsql"

	_ "github.com/go-sql-driver/mysql"
)

var db *xsql.DB
var dsn = "root:root@tcp(127.0.0.1:3306)/example?parseTime=true"
var ctx = context.Background()

func InitDB() {
	var err error
	db, err = xsql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
}

func UserExample() {

	// u := &user.User{
	// 	Name:  "testA",
	// 	Age:   22,
	// 	Ctime: time.Now(),
	// }
	// u1 := &user.User{
	// 	Name:  "testb",
	// 	Age:   22,
	// 	Ctime: time.Now(),
	// }
	// err := user.Create(db).SetUser(u).Save(ctx)
	// fmt.Println(err)

	// err = user.Create(db).SetUser(u1).Save(ctx)
	// fmt.Println(err)

	au, err := user.Find(db).Where(
		user.Or(
			user.IDGT(10),
			user.NameEQ("testb"),
		)).
		Offset(0).
		Limit(3).
		OrderAsc("name").
		All(ctx)
	fmt.Printf("%+v", au)
	fmt.Println(err)

	c, err := user.Find(db).Count(xsql.Distinct(user.Name)).Where(user.Or(
		user.IDGT(10),
		user.NameEQ("testb"),
	)).Int64(ctx)
	fmt.Println(c, err)

	c1, err := user.Find(db).Select(user.Columns...).Where(user.Or(
		user.IDGT(10),
		user.NameEQ("testb"),
	)).All(ctx)
	fmt.Println(c1, err)

	c2, err := user.Find(db).Select(xsql.Sum(user.Age)).Where(user.Or(
		user.IDGT(10),
		user.NameEQ("testb"),
	)).Int64(ctx)
	fmt.Println(c2, err)

	// effect, err := user.Update(db).SetAge(10).WhereP(xsql.EQ(user.ID, 1)).Save(ctx)

	// effect, err = user.Update(db).SetAge(100).SetName("java").Where(user.IDEQ(1)).Save(ctx)

	// effect, err = user.Update(db).AddAge(-100).SetName("java").ByID(5).Save(ctx)

	// effect, err = user.Delete(db).Where(user.And(user.IDEQ(3), user.IDIn(1, 3))).Exec(ctx)

	// effect, err = user.Delete(db).WhereP(xsql.EQ(user.ID, 32)).Exec(ctx)

	// effect, err = user.Delete(db).ByID(2).Exec(ctx)

	tx, _ := db.Begin()
	u2 := &user.User{
		ID:    0,
		Name:  "foo",
		Age:   2,
		Ctime: time.Now(),
	}
	_, err = user.Create(tx).SetUser(u2).Save(ctx)
	if err != nil {
		tx.Rollback()
		return
	}

	effect, err := user.Update(tx).SetAge(100).Where(user.IDEQ(1)).Save(ctx)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	fmt.Println(effect, err)

}
func main() {
	InitDB()
	xsql.Debug()
	//sqlbuild()
	UserExample()
}
