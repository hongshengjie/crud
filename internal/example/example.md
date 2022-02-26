## 用crud自动创建GRPC服务，构建GRPC-Web前端代码，etcd作为服务发现,并用envoy负载均衡


### install
```bash 

makdir example
cd example
crud init 
touch crud/user.sql
```

paste sql content to crud/user.sql

```sql
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
crud -service -protopkg example
```

the api/ proto/ service/ will be generated

```bash
touch main.go
```

paste content to main.go

```go

var port int
var dsn string

func init() {
	flag.IntVar(&port, "port", 9000, "server listen on port")
	flag.StringVar(&dsn, "dsn", "root:123456@tcp(127.0.0.1:3306)/example?parseTime=true", "mysql dsn example(root:123456@tcp(127.0.0.1:3306)/example?parseTime=true)")
}
func main() {
	flag.Parse()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	svr := grpc.NewServer()
	client, err := crud.NewClient(&xsql.Config{
		DSN:          dsn,
		ReadDSN:      []string{dsn},
		Active:       20,
		Idle:         10,
		IdleTimeout:  time.Hour * 4,
		QueryTimeout: time.Second,
		ExecTimeout:  time.Second,
	})
	if err != nil {
		panic(err)
	}
	u := &service.UserServiceImpl{Client: client}
	al := &service.AllTypeTableServiceImpl{Client: client}

	api.RegisterAllTypeTableServiceServer(svr, al)
	api.RegisterUserServiceServer(svr, u)
	grpc_health_v1.RegisterHealthServer(svr, health.NewServer())
	reflection.Register(svr)
	go func() {
		svr.Serve(l)
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svr.GracefulStop()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

```

start 2 grpc services on different port  

```bash
go run main -port 9000 -dsn 'root:123456@tcp(127.0.0.1:3306)/example?parseTime=true'
go run main -port 9001 -dsn 'root:123456@tcp(127.0.0.1:3306)/example?parseTime=true'
```

start envoy proxy 

```bash
envoy -c config/envoy.yaml
```


start front web nodejs server 

```
cd web/
npm i 
npm start
```



open http://localhost:4321/, you will see hello xxx




