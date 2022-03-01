// Code generated by protoc-gen-go-gin. DO NOT EDIT.
// versions:
// - protoc-gen-go-gin v1.0.0
// - protoc             v3.19.3
// source: proto/user.api.proto

package api

import (
	context "context"
	gin "github.com/gin-gonic/gin"
	binding "github.com/gin-gonic/gin/binding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
)

var userServiceGin UserServiceServerGin

// UserServiceServer is the server API for UserService service.
type UserServiceServerGin interface {
	CreateUser(context.Context, *User) (*User, error)
	DeleteUser(context.Context, *UserId) (*emptypb.Empty, error)
	UpdateUser(context.Context, *UpdateUserReq) (*User, error)
	GetUser(context.Context, *UserId) (*User, error)
	ListUsers(context.Context, *ListUsersReq) (*ListUsersResp, error)
}

func createUserG(c *gin.Context) {
	a := &User{}
	err := c.BindWith(a, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type")))
	type Status struct {
		Code    int32       `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	status := &Status{}
	if err != nil {
		status.Code = http.StatusBadRequest
		status.Message = err.Error()
		c.JSON(http.StatusBadRequest, status)
		return
	}
	resp, err := userServiceGin.CreateUser(c, a)
	if err != nil {
		status.Code = http.StatusInternalServerError
		status.Message = err.Error()
		c.JSON(http.StatusInternalServerError, status)
		return
	}
	status.Data = resp
	c.JSON(http.StatusOK, status)
}
func deleteUserG(c *gin.Context) {
	a := &UserId{}
	err := c.BindWith(a, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type")))
	type Status struct {
		Code    int32       `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	status := &Status{}
	if err != nil {
		status.Code = http.StatusBadRequest
		status.Message = err.Error()
		c.JSON(http.StatusBadRequest, status)
		return
	}
	resp, err := userServiceGin.DeleteUser(c, a)
	if err != nil {
		status.Code = http.StatusInternalServerError
		status.Message = err.Error()
		c.JSON(http.StatusInternalServerError, status)
		return
	}
	status.Data = resp
	c.JSON(http.StatusOK, status)
}
func updateUserG(c *gin.Context) {
	a := &UpdateUserReq{}
	err := c.BindWith(a, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type")))
	type Status struct {
		Code    int32       `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	status := &Status{}
	if err != nil {
		status.Code = http.StatusBadRequest
		status.Message = err.Error()
		c.JSON(http.StatusBadRequest, status)
		return
	}
	resp, err := userServiceGin.UpdateUser(c, a)
	if err != nil {
		status.Code = http.StatusInternalServerError
		status.Message = err.Error()
		c.JSON(http.StatusInternalServerError, status)
		return
	}
	status.Data = resp
	c.JSON(http.StatusOK, status)
}
func getUserG(c *gin.Context) {
	a := &UserId{}
	err := c.BindWith(a, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type")))
	type Status struct {
		Code    int32       `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	status := &Status{}
	if err != nil {
		status.Code = http.StatusBadRequest
		status.Message = err.Error()
		c.JSON(http.StatusBadRequest, status)
		return
	}
	resp, err := userServiceGin.GetUser(c, a)
	if err != nil {
		status.Code = http.StatusInternalServerError
		status.Message = err.Error()
		c.JSON(http.StatusInternalServerError, status)
		return
	}
	status.Data = resp
	c.JSON(http.StatusOK, status)
}
func listUsersG(c *gin.Context) {
	a := &ListUsersReq{}
	err := c.BindWith(a, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type")))
	type Status struct {
		Code    int32       `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	status := &Status{}
	if err != nil {
		status.Code = http.StatusBadRequest
		status.Message = err.Error()
		c.JSON(http.StatusBadRequest, status)
		return
	}
	resp, err := userServiceGin.ListUsers(c, a)
	if err != nil {
		status.Code = http.StatusInternalServerError
		status.Message = err.Error()
		c.JSON(http.StatusInternalServerError, status)
		return
	}
	status.Data = resp
	c.JSON(http.StatusOK, status)
}
func RegisterUserServiceGin(e *gin.Engine, svr UserServiceServerGin, middleware map[string][]gin.HandlerFunc) {
	userServiceGin = svr
	e.POST("/example.UserService/CreateUser", append(middleware["/example.UserService/CreateUser"], createUserG)...)
	e.POST("/example.UserService/DeleteUser", append(middleware["/example.UserService/DeleteUser"], deleteUserG)...)
	e.POST("/example.UserService/UpdateUser", append(middleware["/example.UserService/UpdateUser"], updateUserG)...)
	e.POST("/example.UserService/GetUser", append(middleware["/example.UserService/GetUser"], getUserG)...)
	e.POST("/example.UserService/ListUsers", append(middleware["/example.UserService/ListUsers"], listUsersG)...)
}
