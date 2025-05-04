package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
}

func (handler *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", handler.SignUp) // 注册
	server.POST("/users/login", handler.Login)   // 登录
	server.POST("/users/edit", handler.Edit)     // 编辑

	server.GET("/users/profile", handler.Profile) // 查看个人信息

}

func (handler *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq

	// Bind 方法会根据 Content-Type 来解析数据进 req 中
	// 若解析错误，会写入 400 的HTTP状态码
	if err := ctx.Bind(&req); err != nil {
		return
	}

	fmt.Printf("%v\n", req)
	ctx.String(http.StatusOK, "已收到您的注册请求\n")
}

func (handler *UserHandler) Login(ctx *gin.Context) {

}

func (handler *UserHandler) Edit(ctx *gin.Context) {

}

func (handler *UserHandler) Profile(ctx *gin.Context) {

}
