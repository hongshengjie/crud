package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hongshengjie/crud/internal/example/api"
	"github.com/hongshengjie/crud/internal/example/crud"
	"github.com/hongshengjie/crud/internal/example/service"
	"github.com/hongshengjie/crud/xsql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	svr := grpc.NewServer()
	client, err := crud.NewClient(&xsql.Config{
		DSN:          "root:123456@tcp(127.0.0.1:3306)/test?parseTime=true",
		ReadDSN:      []string{"root:123456@tcp(127.0.0.1:3306)/test?parseTime=true"},
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
