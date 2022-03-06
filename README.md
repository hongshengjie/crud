
# crud is a mysql crud code generate tool

## [中文文档](README_zh.md)


## Overview

Crud is a very easy to learn and easy to use ORM framework. Using crud can enable you to complete business requirements quickly, gracefully and with high performance. Currently, MariaDB and MySQL are supported.


- From SQL DDL table structure design to corresponding model and service generation, it conforms to the process of creating tables before writing code

- Supports transactions, row-level locking, for update, lock in share mode

- Elegant API, no ugly hard coding, SQL fragments, all static method calls, and automatic prompt of IDE

- It supports batch insertion, upsert, and automatic assignment of self incrementing ID to structure

- Support context

- High performance. When querying all fields in the table, no reflection is used to create objects, and the performance is consistent with that of native

- Query support forceindex

- Query supports flexible setting of query criteria

- Query supports group by and having

- Query supports scan query results to user-defined structures (using reflection)

- Server code standardization

- Support the generation of proto files and service semi implementation codes containing grpc interface definitions according to SQL DDL table structure definition files


## [example](https://github.com/hongshengjie/crud-exmaple)

## Getting Started 

### install

```bash

go install  github.com/hongshengjie/crud@latest

```
### Using the command line

```bash
crud -h 

Usage of crud:
  -protopkg string
        -protopkg  proto package field value
  -service
        -service  generate GRPC proto message and service implementation
```

```example
#  generation crud directory
crud init

# Put user.sql In the crud directory sql


# According to the table structure, generate the proto file of grpc interface and service semi implementation code for the CRUD of the table
crud  -service -protopkg example

```

## Init


### Init db
```go

db, _ = sql.Open("mysql","user:pwd@tcp(127.0.0.1:3306)/example?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8")

```

### Or the client wrapped in curd has read-write separation and context read-write timeout configuration ability

```go

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
```


### As user SQL table creation file as an example

```SQL
CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id字段',
  `name` varchar(100) NOT NULL COMMENT '名称',
  `age` int(11) NOT NULL DEFAULT '0' COMMENT '年龄',
  `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `ix_name` (`name`) USING BTREE,
  KEY `ix_mtime` (`mtime`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4
```

```bash
# exec bcurd under example
crud 

# The following user directories and files will be generated
example/
├── crud
│   ├── user
│   │   ├── builder.go
│   │   ├── model.go
│   │   └── where.go
│   └── user.sql

```
> The user directory is generated above, and the package name is user.

## CRUD API

### Create

#### Single insert
```go
u := &user.User{
	ID:    0,
	Name:  "shengjie",
	Age:   18,
	Ctime: time.Now(),
	Mtime: time.Now(),
}
effect, err := user.
	Create(db).
	SetUser(u).
	Save(ctx)

fmt.Println(err, u, effect)
```
> Insert a single record. Before inserting the above code, set id = 0 and ID field as auto_increment, crud will assign the self increasing ID generated by the database to u.ID, and the u.ID after insertion is the ID generated by DB.


#### Batch insert

```go
u1 := &user.User{
	ID:   0,
	Name: "shengjie",
	Age:  22,
	Ctime: time.Now(),
	Mtime: time.Now(),
}
u2 := &user.User{
	ID:   0,
	Name: "shengjie2",
	Age:  22,
	Ctime: time.Now(),
	Mtime: time.Now(),
}
effect, err = user.
	Create(db).
	SetUser(u1,u2).
	Save(ctx)
fmt.Println(effect, err, u1, u2)
```
> The above two records will be inserted. The lastinsertid returned by each record cannot be obtained during batch insertion, so the ID of U1 and U2 after insertion are 0.

#### Upsert

```go
a := &user.User{
	ID:   1,
	Name: "shengjie",
	Age:  19,
}
effect, err := user.
	Create(db).
	SetUser(a).
	Upsert(ctx)

fmt.Println(effect, err, a)
```

> If a unique key conflict is encountered during insertion, all fields will be updated with the new value passed in.

#### Attention
1. During batch insertion, the structure will not take the lastinsertid returned by the database.

2. If the default value of the database is not the zero value of its type, and the corresponding structure does not set the value of this field in the insertion operation, crud will insert dB with the zero value of its type.

3. It is strongly recommended that the value type must use: not null default 0, and the string type must use: not null default ""


### Query

#### Query a single record
```go
u, err = user.
	Find(db).
	Where(user.IDEQ(1)).
	One(ctx)

fmt.Println(u, err)
```
> One(ctx) will automatically set the query statement limit = 1.


#### Query multiple records
```go
list, err := user.
	Find(db).
	Where(
		user.AgeIn(18, 20, 30),
		).
	All(ctx)

liststr, _ := json.Marshal(list)
fmt.Printf("%+v %+v \n", string(liststr), err)
```
> Query all records with ages of 18, 20 and 30, and All(ctx) returns []*user.User .

```go
list, err := user.Find(db)).
	Where(user.Or(
		user.IDGT(97),
		user.AgeIn(10, 20, 30),
		)).
	OrderAsc(user.Age).
	Offset(2).
	Limit(20).
	All(ctx)
fmt.Printf("%+v %+v \n", list, err)
```
> Rich query criteria expression support

```go
list, err := user.
	Find(db).
	Where(
		user.NameContains("java"),
		).
	All(ctx)

list, err = user.
	Find(db).
	Where(
		user.NameHasPrefix("java"),
		).
	All(ctx)
```
> String field fuzzy query and prefix matching.


#### The query result is a single column
```go
count, err := user.
	Find(db).
	Count().
	Where(user.IDGT(0)).
	Int64(ctx)

fmt.Println(count, err)

names, err := user.
	Find(db).
	Select(user.Name).
	Limit(2).
	Where(
		user.IDIn(1, 2, 3, 4),
		).
	Strings(ctx)
fmt.Println(names, err)
```
> Count() query the quantity of qualified records; If the returned result contains only one column and only one row, Int64 and String can be used; If the returned result contains only one column and multiple rows, you can use Int64s and Strings to get the list.

#### Select () parameter description

```go
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

```
> The SQL statements and results generated by the above two queries are the same, but they are very different internally.
> When Select() does not specify parameters, crud will find all fields corresponding to the model. When returning results, it does not use reflection to create objects, If the return value has a null value, an error will be returned.
> When Select(user.Columns()...) When all column names are specified, the returned results will use reflection to create objects. If the return value has a null value, no error will be reported, and the default value of this field is zero

### Transaction support

```go
tx, err := db.Begin(ctx)
if err != nil {
	return err
}
u1 := &user.User{
	ID:   0,
	Name: "shengjie",
	Age:  18,
}
_, err = user.
	Create(tx).
	SetUser(u1).
	Save(ctx)
if err != nil {
	return tx.Rollback()
}
effect, err := user.
	Update(tx).
	SetAge(100).
	Where(
		user.IDEQ(u1.ID)
		).
	Save(ctx)

if err != nil {
	return tx.Rollback()
}
fmt.Println(effect, err)
return tx.Commit()
```



### Advanced Query

#### Custom query result acquisition
```go
type GroupResutl struct {
	Name string `json:"name"` 
	Cnt  int64  `json:"cnt"`
}

result := []*GroupResutl{}
err := user.Find(db).
	Select(
		user.Name,
		xsql.As(xsql.Count("*"), "cnt"),
		).
	ForceIndex(`ix_name`).
	GroupBy(user.Name).
	Having(xsql.GT(`cnt`, 1)).
	Slice(ctx, &result)
// SELECT `name`, COUNT(*) AS `cnt` FROM `user` FORCE INDEX (`ix_name`) GROUP BY `name` HAVING `cnt` > ? 
fmt.Println(err, result)
b, _ := json.Marshal(result)
fmt.Println(string(b))

```
> The above uses force index, groupby, having, count and as to scan the user-defined query results into the user-defined structure. The JSON tag of the structure needs to be consistent with the column name returned from the query results, and the fields in the structure need to be capitalized.

> Slice(context,interface{}):The second parameter of the method needs to be passed in: a pointer to a structure slice


### Update
```go

effect, err := user.
	Update(db).
	SetAge(10).
	Where(user.NameEQ("java")).
	Save(ctx)

fmt.Println(effect, err)


effect, err = user.
	Update(db).
	SetAge(100).
	SetName("java").
	SetName("python").
	Where(user.IDEQ(97)).
	Save(ctx)

fmt.Println(effect, err)

// update `user` set `age` = COALESCE(`age`, 0) + -100, `name` = 'java' where `id` = 5
effect, err = user.
	Update(db).
	AddAge(-100).
	SetName("java").
	Where(user.IDEQ(97)).
	Save(ctx)
fmt.Println(effect, err)

```
### Delete
```go

effect, err = user.
	Delete(db).
	Where(
		user.And(
			user.IDEQ(3), 
			user.IDIn(1, 3),
		)).
	Exec(ctx)

```
> It is only executed when the Exec method is called


### Debug Log

```go
_, err := user.
	Create(xsql.Debug(db)).
	SetUser(u).
	Save(ctx)

fmt.Println(err)
```
> The generated SQL statement and parameters will be printed

## Generate grpc interface definition proto file and service implementation code

This function helps us generate a lot of cumbersome code that needs to be written by ourselves. For example, a project needs to manage the background, and the interfaces for adding, deleting, modifying and querying need to be built. If we can complete the interface writing with a little modification on the basis of the generated code, the business interface will be realized quickly and with quality.

### Dependencies

1. protoc
2. protoc-gen-go
3. protoc-gen-go-grpc
4. make sure /usr/local/include have google/protobuf/empty.proto file


```
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

### usage
```bash

crud -service -protopkg example


example/
├── api
│   ├── user.api_grpc.pb.go
│   └── user.api.pb.go
├── crud
│   ├── aa_client.go
│   ├── user
│   │   ├── builder.go
│   │   ├── model.go
│   │   └── where.go
│   └── user.sql
├── proto
│   └── user.api.proto
└── service
    └── user.service.go

```
> There are more api and service directories and proto files.

### proto example 
usr.api.proto
```proto
syntax="proto3";
package user;
option go_package = "/api";

import "google/protobuf/empty.proto";

service UserService { 
    rpc CreateUser(User)returns(User);
    rpc DeleteUser(UserId)returns(google.protobuf.Empty);
    rpc UpdateUser(UpdateUserReq)returns(User);
    rpc GetUser(UserId)returns(User);
    rpc ListUsers(ListUsersReq)returns(ListUsersResp);
}

message User {
    //id字段
    int64	id = 1 ;
    //名称
    string	name = 2 ;
    //年龄
    int64	age = 3 ;
    //创建时间
    string	ctime = 4 ;
    //更新时间
    string	mtime = 5 ;  
}

message UserId{
    int64 id = 1 ;
}

message UpdateUserReq{

    User user = 1 ;

    repeated string update_mask  = 2 ;
}


message ListUsersReq{
    // 
    int64 page = 1 ;
    // default 20
    int64 page_size = 2 ;
    // order by  for example :  [-id]  -表示：倒序排序
    string orderby = 3 ; 
     // filter
    repeated UserFilter filter = 4 ;
}

message UserFilter{
    string field = 1;
    string op = 2;
    string value = 3;
}


message ListUsersResp{

    repeated User users = 1 ;

    int64 total_count = 2 ;
    
    int64 page_count = 3 ;
}


```
> Generate a proto message corresponding to the table structure, and the generated API file conforms to Google API design specification.

### service example 
user.service.go
```go
package service

import (
	"context"
	"errors"
	"github.com/hongshengjie/crud/internal/example/api"
	"github.com/hongshengjie/crud/internal/example/crud"
	"github.com/hongshengjie/crud/internal/example/crud/user"
	"github.com/hongshengjie/crud/xsql"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"strings"
	"time"
)

// UserServiceImpl UserServiceImpl
type UserServiceImpl struct {
	Client *crud.Client
}

// CreateUser CreateUser
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *api.User) (*api.User, error) {

	// do some parameter check
	// if req.GetXXXX() != 0 {
	// 	return nil, errors.New(-1, "参数错误")
	// }
	a := &user.User{
		Id:    0,
		Name:  req.GetName(),
		Age:   req.GetAge(),
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	var err error
	_, err = s.Client.User.
		Create().
		SetUser(a).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// query after create and return
	a2, err := s.Client.Master.User.
		Find().
		Where(
			user.IdEQ(a.Id),
		).
		One(ctx)
	if err != nil {
		return nil, err
	}
	return convertUser(a2), nil
}

// DeleteUser DeleteUser
func (s *UserServiceImpl) DeletesUser(ctx context.Context, req *api.UserId) (*emptypb.Empty, error) {
	_, err := s.Client.User.
		Delete().
		Where(
			user.IdEQ(req.GetId()),
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Updateuser UpdateUser
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *api.UpdateUserReq) (*api.User, error) {

	if len(req.GetUpdateMask()) == 0 {
		return nil, errors.New("update_mask empty")
	}
	update := s.Client.User.Update()
	for _, v := range req.GetUpdateMask() {
		switch v {
		case "user.name":
			update.SetName(req.GetUser().GetName())
		case "user.age":
			update.SetAge(req.GetUser().GetAge())
		case "user.ctime":
			t, err := time.ParseInLocation("2006-01-02 15:04:05", req.GetUser().GetCtime(), time.Local)
			if err != nil {
				return nil, err
			}
			update.SetCtime(t)
		case "user.mtime":
			t, err := time.ParseInLocation("2006-01-02 15:04:05", req.GetUser().GetMtime(), time.Local)
			if err != nil {
				return nil, err
			}
			update.SetMtime(t)
		}
	}
	_, err := update.
		Where(
			user.IdEQ(req.GetUser().GetId()),
		).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// query after update and return
	a, err := s.Client.Master.User.
		Find().
		Where(
			user.IdEQ(req.GetUser().GetId()),
		).
		One(ctx)
	if err != nil {
		return nil, err
	}
	return convertUser(a), nil
}

// GetUser GetUser
func (s *UserServiceImpl) GetUser(ctx context.Context, req *api.UserId) (*api.User, error) {
	a, err := s.Client.User.
		Find().
		Where(
			user.IdEQ(req.GetId()),
		).
		One(ctx)
	if err != nil {
		return nil, err
	}
	return convertUser(a), nil
}

// ListUsers ListUsers
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *api.ListUsersReq) (*api.ListUsersResp, error) {
	page := req.GetPage()
	size := req.GetPageSize()
	if size <= 0 {
		size = 20
	}
	offset := size * (page - 1)
	if offset < 0 {
		offset = 0
	}
	finder := s.Client.User.
		Find().
		Offset(offset).
		Limit(size)

	if req.GetOrderby() != "" {
		odb := strings.TrimPrefix(req.GetOrderby(), "-")
		if odb == req.GetOrderby() {
			finder.OrderAsc(odb)
		} else {
			finder.OrderDesc(odb)
		}
	}
	counter := s.Client.User.
		Find().
		Count()

	var ps []*xsql.Predicate
	for _, v := range req.GetFilter() {
		p, err := xsql.GenP(v.Field, v.Op, v.Value)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}

	list, err := finder.WhereP(ps...).All(ctx)
	if err != nil {
		return nil, err
	}

	count, err := counter.WhereP(ps...).Int64(ctx)
	if err != nil {
		return nil, err
	}
	pageCount := int64(math.Ceil(float64(count) / float64(size)))

	return &api.ListUsersResp{Users: convertUserList(list), TotalCount: count, PageCount: pageCount}, nil
}

func convertUser(a *user.User) *api.User {
	return &api.User{
		Id:    a.Id,
		Name:  a.Name,
		Age:   a.Age,
		Ctime: a.Ctime.Format("2006-01-02 15:04:05"),
		Mtime: a.Mtime.Format("2006-01-02 15:04:05"),
	}
}

func convertUserList(list []*user.User) []*api.User {
	ret := make([]*api.User, 0, len(list))
	for _, v := range list {
		ret = append(ret, convertUser(v))
	}
	return ret
}

```
> The semi implementation code of the above service only needs to add some parameter verification, or automatically generate the message conversion code from the DB layer model structure to the API layer according to the code of the condition filter, which is convenient and flexible.




> The project is inspired by [facebook/ent](https://github.com/ent/ent) 
