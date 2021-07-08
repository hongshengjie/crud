# crud    a crud code generater for mysql 

## Feature

this tool can help you generate some tedious crud code.

- generate go model struct and fields that corresponding to table structure.
- generate elegent insert,select,update,delete sql builder api.
- generate elegent where condition sql statement construction function.
- support transaction.


## Origion 

Previously, database/sql package in standard library was used to write Dao layer code,
Handwritten SQL statements and scan methods make me tired and annoyed,
Then when self-test, there are still many problems, there is no development efficiency, also very tired.
Gorm these advanced framework,but feel API is not elegant.
This makes me very sad. So there's crud.



## Supported database

- maridb、mysql8、mysql5

## Install

```shell
git clone git@github.com:hongshengjie/crud.git
cd crud
go install 
```

## How to use

for example: 
```sql
CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id字段',
  `name` varchar(100) NOT NULL COMMENT '名称',
  `age` int(11) NOT NULL DEFAULT 0 COMMENT '年龄',
  `ctime` datetime NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_name_IDX` (`name`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4

```
attention: table must hava a prime key ,and every field must be not null.

```shell
action:

cd you_project/model

crud -path user.sql

OR

crud -path ./tables/

OR

crud -dsn='user:password@tcp(127.0.0.1:3306)/example?parseTime=true' -table=user

result: 

user                    // directory name (alse package name)
├── builder.gen.go    // include insert,select, delete, update sql builder
├── custom.go         // custom code, file, if you have a custom method, write it here, it will not be covered when it is regenerated
├── model.gen.go      // the golang structure corresponding to the table structure, each field name constant and field name slice
└── where.gen.go      // some comparison methods corresponding to each field, such as IDEQ()
```


The generated package has no state;

The business code is synchronously transferred to *sql.DB or *sql.Tx to use the generated method;

For example, the following insertion method:

```go
err := user.Create(db).SetUser(u).Save(ctx);
```

Directly reference the Create method under the user package, and the parameter DB is passed in from outside;

The advantage of no state is that the method under the package can be directly referenced without creating an instance;  

Command parameter description

```shell
crud 

  -dsn string, Database connection
        
  -table string, table name

  -path string, generate code from DDL sql file or directory contain .sql extension files
        
```



## Example

- Init DB 

```go
db, _ = sql.Open("mysql","user:password@tcp(127.0.0.1:3306)/example?parseTime=true")
```

- Insert 
```go
u := &user.User{
	ID: 0,
	Name:  "testA",
	Age:   22,
	Ctime: time.Now(),
}
// If the ID of u is auto incremental and save is successful, the returned ID field will be automatically set to u.ID
err := user.Create(db).SetUser(u).Save(ctx)
```

- Update 

```go 
user.Update(db).SetAge(100).SetName("java").Where(user.IDEQ(1)).Save(ctx)
// update `user` set `age` = `age` + -100, `name` = 'java' where `id` = 5
user.Update(db).AddAge(-100).SetName("java").ByID(5).Save(ctx)
```

- Delete 

```go
user.Delete(db).Where(user.And(user.IDEQ(3), user.IDIn(1, 3))).Exec(ctx)
user.Delete(db).ByID(2).Exec(ctx)
```

- Query 

```go
au, err := user.Find(db).Where(
	user.Or(
		user.IDEQ(10),
		user.NameEQ("testb"),
	)).
	Offset(0).
	Limit(3).
	OrderAsc(user.Name).
	One(ctx)
fmt.Println(au, err)

// Count() method
c, err := user.Find(db).Count(xsql.Distinct(user.ID)).Where(user.IDGT(0)).Int64(ctx)
fmt.Println(c, err)

```

- transaction

```go
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
```

## TODO
- [ ] (support PostgreSQL)


## Reference
[ent](https://github.com/ent/ent)