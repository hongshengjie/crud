package main

import (
	"context"
	"example/api"
	"example/crud"
	"example/discovery"
	"example/service"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hongshengjie/crud/xsql"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

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
