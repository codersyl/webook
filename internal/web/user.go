package web

import (
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func (handler *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", u.SignUp) // 注册
	server.POST("/users/login", u.Login)   // 登录
	server.POST("/users/edit", u.Edit)     // 编辑

	server.GET("/users/profile", u.Profile) // 查看个人信息

}
func (handler *UserHandler) SignUp(ctx *gin.Context) {

}

func (handler *UserHandler) Login(ctx *gin.Context) {

}

func (handler *UserHandler) Edit(ctx *gin.Context) {

}

func (handler *UserHandler) Profile(ctx *gin.Context) {

}
