package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
	"webook_Rouge/internal/domain"
	"webook_Rouge/internal/service"
)

const (
	emailRegexPattern    = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	passwordRegexPattern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*]).{8,16}$`
)

var (
	ErrUserDuplicatedEmail = service.ErrUserDuplicatedEmail
)

type UserHandler struct {
	svc           *service.UserService
	emailRegex    *regexp.Regexp
	passwordRegex *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:           svc,
		emailRegex:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", u.SignUp) // 注册
	server.POST("/users/login", u.Login)   // 登录
	server.POST("/users/edit", u.Edit)     // 编辑

	server.GET("/users/profile", u.Profile) // 查看个人信息
	return

}

func (u *UserHandler) SignUp(ctx *gin.Context) {
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

	// 邮箱校验
	ok, err := u.emailRegex.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "内部错误\n")
	}
	if !ok {
		ctx.String(http.StatusOK, "您的邮箱格式不对\n")
		return
	}

	// 两次密码输入不一致
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致\n")
		return
	}

	// 密码校验
	ok, err = u.passwordRegex.MatchString(req.Password)
	if err != nil {
		fmt.Printf("密码校验正则表达式 错误\n")
		ctx.String(http.StatusInternalServerError, "内部错误\n")
	}
	if !ok {
		ctx.String(http.StatusOK, "密码长度需在8-16位，且包含至少一个小写字母、一个大写字母、一个数字、一个特殊字符（!@#$%^&*）\n")
		return
	}

	// 调用service 存储用户信息
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == ErrUserDuplicatedEmail {
		ctx.String(http.StatusOK, "邮箱冲突\n")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误\n")
		return
	}

	ctx.String(http.StatusOK, "注册成功\n")
	return
}

func (handler *UserHandler) Login(ctx *gin.Context) {
	type LogIn struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LogIn

	// 解析登录请求
	if err := ctx.Bind(&req); err != nil {
		return
	}

	err := handler.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码错误\n")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误\n")
		return
	}

	ctx.String(http.StatusOK, "登录成功\n")
}

func (handler *UserHandler) Edit(ctx *gin.Context) {

}

func (handler *UserHandler) Profile(ctx *gin.Context) {

}
