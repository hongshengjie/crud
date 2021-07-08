# crud 一个goolang服务端数据库增删改查代码生成器

##

这是一个工具，能帮助你生成一些乏味的dao层代码  
支持的功能  
- 支持事务
- 针对每张表生成单独的package
	1. 生成和表结构对应的golang结构体和字段列表
	2. 生成优雅的insert、select、update、delete相关API
	3. 生成表中一些字段where语句的比较方法

## 起源

原先一直使用标准库中的sql包来写dao层代码  
手写SQL语句和Scan方法让我觉得厌倦和烦恼  
然后自测的时候，问题还不少，没有开发效率，也很累。

gorm这些高级的框架,感觉API不优雅 

这让我十分的sad  
于是就有了crud项目  



## 支持的数据库

1. maridb、mysql8、mysql5



## 安装

```
git clone git@github.com:hongshengjie/crud.git
cd crud/crud/
go install 
```

## 如何使用

```shell
action:

cd you_project/dao/

crud -path user.sql

OR

crud -path ./sqlfile/

OR

crud -dsn='user:password@tcp(127.0.0.1:3306)/example?parseTime=true' -table=user

result: 

user                    // 目录名
├── builder.gen.go    // 生成 insert select delete update 构建者
├── custom.go         // 自定义代码，文件，如果你有自定义的方法，写在这里，重新生成的时候不会被覆盖
├── model.gen.go      // 生成和表结构对应的 golang 结构体、 每个字段名常量、字段名Slice
└── where.gen.go      // 生成一些各个字段相对应的 比较方法比如 IDEQ()
```
生成的包是没有状态的;       
业务代码同步传递 *sql.DB 或者 *sql.Tx 来使用生成的方法;       
比如下面的插入方法:   
err := user.Create(db).SetUser(u).Save(ctx);      
直接引用user包下的Create方法,参数db由外部传入;      
没有状态的好处就是包下的方法可以直接引用，而不需要New出一个实例;      

命令参数说明

```
crud 

  -dsn string 数据库连接
        
  -table string 指定表名

  -path string .sql file path or dir generate code from DDL sql file
        
```



## Example

- Init DB 初始化db

```go
db, _ = sql.Open("mysql","user:password@tcp(127.0.0.1:3306)/example?parseTime=true")
```

- Insert 插入方法
```go
	u := &user.User{
		ID: 0,
		Name:  "testA",
		Age:   22,
		Ctime: time.Now(),
	}
	// 如果u的ID是auto-increment的,Save执行成功，会自动设置返回的ID字段到u.ID上
	err := user.Create(db).SetUser(u).Save(ctx)
```

- Update 更新字段

```go 
    // 使用WhereP可以通过 xsql包下的方法，生成比较复杂的自定义where条件，在调用Save()方法的时候在真正执行
	user.Update(db).SetAge(10).WhereP(xsql.EQ(user.ID, 1)).Save(ctx)
    // 使用工具帮你生成的方法 IDEQ() SetName()  SetAge() 等方法
	user.Update(db).SetAge(100).SetName("java").Where(user.IDEQ(1)).Save(ctx)
	// 数字字段可以使用AddAge()方法来生成 x = x + ? 这种表达式
    // update `user` set `age` = `age` + -100 , `name` = 'java' where `id` = 5
	user.Update(db).AddAge(-100).SetName("java").ByID(5).Save(ctx)
```

- Delete 删除记录

```go
	// 在调用Exec方法的时候才真正执行
	user.Delete(db).WhereP(xsql.EQ(user.ID, 32)).Exec(ctx)
	user.Delete(db).Where(user.And(user.IDEQ(3), user.IDIn(1, 3))).Exec(ctx)
	user.Delete(db).ByID(2).Exec(ctx)
```

- Query 查询

```go
	// 查询语句One() 或者All() 支持单条查询和批量查询，可以构建比较复杂的 where语句和 offset limit  orderby 等
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

	// Count() 方法
	c, err := user.Find(db).Count(xsql.Distinct(user.ID)).Where(user.IDGT(0)).Int64(ctx)
	fmt.Println(c, err)

```

- 对事务的支持 

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


## FAQ

1. 表字段修改了怎么办?
	- 表字段修改了需要重新运行生成代码命令；第一次生成时会在costom.go文件的注释中写入运行的命令,但是你后续改动的话，需要手动修改。  
	
2. 是否支持分表?
	-  支持的，默认是的表都是生成代码中的table字段，可以通过Table()方法来修改生成语句中的表名，表名的生成需要业务自定义。   
   
3. where.gen.go 文件为啥不包含所有的字段?
	-  crud默认会对具有索引的字段生成相应的=,<>,>,>=,<,<=,in,not in等查询方法,目的是尽量生成需要字段的方法，减少总的生成代码量。
	-  如果你需要所有字段都生成请指定 -fields=all 
	-  如果你不想所有的字段都生成、但是某个字段字段又不在索引，你可以指定 -fields=xxx,xx 

4. 如何查看生成的sql做一些Debug自测工作
	- 在你项目中某个位置（最好是dao.go文件初始化db的时候）添加一行xsql.Debug(), 在执行到的时候会已日志方式打印出相关的 sql 和 参数列表，记得上线的时候把这段代码删除，否则会有大量的sql日志


## TODO
- [ ] (支持postgresql)
