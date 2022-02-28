package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var userHttpService UserServiceServer

func CreateUserG(c *gin.Context) {
	a := &User{}
	err := c.BindWith(a, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type")))
	if err != nil {
		return
	}
	resp, err := userHttpService.CreateUser(c, a)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, resp)
}

func RegisterHttp(e *gin.Engine, svr UserServiceServer) {
	userHttpService = svr
	e.GET("/xxx/xxx", CreateUserG)

}
