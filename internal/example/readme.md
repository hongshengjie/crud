## 用crud快速创建GRPC微服务

- [x] 用crud生成持久层代码、proto接口定义、grpc服务实现
- [x] 用etcd作服务发现
- [x] 使用grpc-web生成前端js代码
- [x] 使用envoy作为最外层网关代理、均衡负载、健康检查等功能
- [x] 使用protoc-gen-go-gin 生成对应的http服务
- [ ] (使envoy可以根据后端节点变更，自动更新代理配置)

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

const (
	appID = "example.user"
)

func init() {
	flag.IntVar(&port, "port", 9000, "server listen on port")
	flag.StringVar(&dsn, "dsn", "root:123456@tcp(127.0.0.1:3306)/test?parseTime=true", "mysql dsn example(root:123456@tcp(127.0.0.1:3306)/example?parseTime=true)")
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

	m := cmux.New(l)

	e := gin.Default()
	api.RegisterUserServiceGin(e, u, nil)
	api.RegisterAllTypeTableServiceGin(e, al, nil)
	// Match connections in order:
	// First grpc, then HTTP, and otherwise Go RPC/TCP.
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())
	hsvr := &http.Server{
		Handler: e,
	}
	go func() {
		go svr.Serve(grpcL)
		go hsvr.Serve(httpL)
		go m.Serve()
	}()

	instanceID := appID + "/" + uuid.New().String()
	err = discovery.Register(context.Background(), appID, instanceID, fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		panic(err)
	}
	fmt.Println(instanceID)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			discovery.DeleteRegister(context.Background(), instanceID)
			hsvr.Shutdown(context.Background())
			svr.GracefulStop()
			m.Close()
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




