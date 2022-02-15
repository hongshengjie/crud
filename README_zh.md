
# crud is a mysql crud code generate tool



## 概览

crud 是一个非常易学好用的ORM框架，使用crud可以让你快速，优雅，且高性能的实现业务需求。目前支持mariadb、mysql。

- 从SQL DDL表结构设计到对应的Model，Service生成，符合先建表再写代码的流程
- 支持事务,row-level locking 、FOR UPDATE 、LOCK IN SHARE MODE
- 优雅的API，无需丑陋的硬编码，以及sql片段，全静态方法调用，IDE自动提示
- 支持批量插入、Upsert、自增id自动赋值到结构体
- 支持Context
- 高性能，在查询表中所有字段时候，不使用反射创建对象，性能和原生一致
- 查询支持 ForceIndex
- 查询支持 灵活得设置查询条件
- 查询支持 GROUP BY、HAVING 
- 查询支持 查询结果Scan到自定义结构体（使用反射）
- 服务端代码标准化
- 表结构变更可以记录在仓库中
- 支持根据SQL DDL表结构定义文件生成包含GRPC接口定义的proto文件 和 Service半实现代码

## 开始

### 安装

```bash

go install  github.com/hongshengjie/crud@latest

```
### 使用命令行

```bash
crud -h 

Usage of crud:
  -protopkg string
        -protopkg  proto package field value
  -service
        -service  generate GRPC proto message and service implementation
```

```example
在项目下创建crud目录
crud init

在crud目录放入user.sql

# 根据表结构 生成针对该表的增删改查GRPC接口的proto文件以及 Service半实现代码
crud -service -protopkg example

```

## 初始化


### 初始化db
```go
db, _ = sql.Open("mysql","user:pwd@tcp(127.0.0.1:3306)/example?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8")

```

### 或者使用curd包装的client, 拥有读写分离，Context读写超时配置功能
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


### 以user.sql建表文件为例
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
# 在example执行
crud  

# 会生成如下目录
example/
├── crud
│   ├── user
│   │   ├── builder.go
│   │   ├── model.go
│   │   └── where.go
│   └── user.sql

```
> 以上生成user目录，且package 名称为user。

## CRUD 接口使用

### Create

#### 单条插入
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
> 插入单条记录 以上代码插入前需设置ID=0，ID字段为auto_increment，crud会把数据库生成的自增ID赋值给u.ID,插入后u.ID 为db为其生成的ID。

#### 批量插入

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
> 以上会插入2条记录，批量插入的时候无法获取到每条记录返回的LastInsertId, 所以执行插入后 u1 和u2 的ID都为0。

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

> 如果插入的时候遇到唯一键冲突,那么会把所有字段全都更新为传入的新值。

#### 注意点
1. 批量插入的时候结构体不会取数据库返回的LastInsertId
2. 如果数据库的默认值不是其类型的零值，且在插入的操作中相应结构体没有设置该字段的值，那么crud会以其类型的零值插入db
3. 强烈建议:数值类型必须使用：NOT NULL DEFAULT 0 字符类型必须使用：NOT NULL DEFAULT ""

### Query

#### 查询单条记录
```go
u, err = user.
	Find(db).
	Where(user.IDEQ(1)).
	One(ctx)

fmt.Println(u, err)
```
> One(ctx) 会自动设置查询语句limit = 1。


#### 查询多条记录
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
> 查询年龄为18，20，30的所有记录，All(ctx)返回的是[]*user.User。

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
> 丰富的查询条件表达支持

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
> 字符串字段模糊查询和前缀匹配。


#### 查询结果为单列
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
> Count()查询符合条件记录的数量；如果返回结果只包含一列,且只有一行可以使用Int64、String ；如果返回的结果只包含一列，且有多行，可以用Int64s、Strings得到列表。


#### Select()参数说明

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
> 以上两个查询生成的sql语句和结果相同，但是内部有很大的不一样。
> 当Select() 不指定参数的时候，crud会查找model对应的所有字段，返回结果时不使用反射创建对象,如果返回值有NULL值,则会报错。
> 当Select(user.Columns()...) 指定所有列名时，返回结果会使用反射来创建对象，返回值如果有NULL值不会报错，该字段默认零值。

### 事务支持

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

#### 自定义查询结果获取
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
> 以上使用了 Force Index 、 GroupBy 、 Having 、Count 、AS 、 把自定义查询结果扫描到自定义的结构体，其中结构体的json tag 需要和查询结果的返回的列名一致，结构体中的字段需要大写。

> Slice(context,interface{})方法第二个参数需要传入的是：指向某个结构体Slice的指针。


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
> 在调用Exec方法的时候才真正执行,线上数据库账号不一定有删除的权限，可以用update来改为软删除。

### 调试日志

```go
_, err := user.
	Create(xsql.Debug(db)).
	SetUser(u).
	Save(ctx)

fmt.Println(err)
```
> 会打印出生成的sql语句和参数

## 生成GRPC接口定义proto文件和服务实现代码

这个功能帮助我们生成很多需要自己手写的繁琐代码，比如某个项目需要管理后台，增删改查的接口都是需要搭建的，假如在生成的代码的基础上做少许修改就能完成接口编写，那么业务接口实现的会又快又有质量。

### 要提前安装的工具

1. protoc
2. protoc-gen-go
3. protoc-gen-go-grpc
4. make sure /usr/local/include have google/protobuf/empty.proto file

```
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

```

### 用法
```bash
# 执行

crud -service -protopkg example

# 会生成如下目录
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
> 多了 api、proto、service 目录。

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
> 生成和表结构一一对应的 message ,生成的api文件符合google 设计规范。

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
> 以上service的半实现代码只需要自己加一些参数校验，或者根据条件filter的代码，自动生成了db层model结构体的到api层的message转化代码,方便灵活。




> The project is inspired by [facebook/ent](https://github.com/ent/ent) 
