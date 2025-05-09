package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
	"webook_Rouge/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"webook_Rouge/internal/domain"
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
	//server.POST("/users/login", u.Login)   // 登录
	server.POST("/users/login", u.LoginJWT) // 登录
	server.POST("/users/edit", u.Edit)      // 编辑

	server.GET("/users/profile", u.ProfileJWT) // 查看个人信息
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

	u, err := handler.svc.Login(ctx, domain.User{
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

	// 设置session
	sess := sessions.Default(ctx)
	sess.Set("userId", u.ID)
	sess.Options(sessions.Options{ // Gin用这些配置来初始化Cookie
		//Secure:   true, // 开发环境建议默认开启 Secure 与 HttpOnly
		// HttpOnly: true,
		MaxAge: 30 * 60, // 登录状态 30min会过期
	})
	sess.Save()
	ctx.String(http.StatusOK, "登录成功\n")
}

func (handler *UserHandler) LoginJWT(ctx *gin.Context) {
	type LogIn struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LogIn

	// 解析登录请求
	if err := ctx.Bind(&req); err != nil {
		return
	}

	u, err := handler.svc.Login(ctx, domain.User{
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

	// 使用JWT设置登录态
	// 生成一个JWT token

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)), // 2天过期
		},
		Uid: u.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	key32_ForToken := "iFyeVYqAZPMY2p2Jma6zn22jxbKH6TCI" // 随机生成的
	tokenStr, err := token.SignedString([]byte(key32_ForToken))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统内部错误\n")
		return
	}

	// 控制台输出
	fmt.Println("token Start")
	fmt.Printf(tokenStr)
	fmt.Println("")
	fmt.Println("token End")

	ctx.Header(("x-jwt-token"), tokenStr)

	sess := sessions.Default(ctx)
	sess.Set("userId", u.ID)
	sess.Options(sessions.Options{ // Gin用这些配置来初始化Cookie
		//Secure:   true, // 开发环境建议默认开启 Secure 与 HttpOnly
		// HttpOnly: true,
		MaxAge: 30 * 60, // 登录状态 30min会过期
	})
	sess.Save()
	ctx.String(http.StatusOK, "登录成功\n")
}

func (handler *UserHandler) LogOut(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{ // Gin用这些配置来初始化Cookie
		//Secure:   true, // 开发环境建议默认开启 Secure 与 HttpOnly
		// HttpOnly: true,
		MaxAge: -1, // 删除Cookie
	})
	sess.Save()
	ctx.String(http.StatusOK, "登出成功\n")
}

func (handler *UserHandler) Edit(ctx *gin.Context) {

}

func (handler *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, ok := ctx.Get("claims")
	//if !ok { // 按理来说一条逻辑写下来必然是有uid的，如果为了保险，也需要判断一下
	//	ctx.AbortWithStatus(http.StatusInternalServerError)
	//	return
	//}
	claims, ok := c.(*UserClaims) // 类型断言
	if !ok {
		// 断言失败
		// 有这个判断就不需要上一个判断了
		// if ctx.Get("Uid") false, uidRaw will be nil, nil cast into int64 will fail
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Println(">> Uid : ", claims.Uid)
	ctx.String(http.StatusOK, "Profile页面（迫真\n")
}

func (handler *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Profile页面（迫真\n")
}

type UserClaims struct {
	jwt.RegisteredClaims // 实现了Claims的接口，直接组合这个类型，可免去自己实现Claims接口的麻烦
	// 另外加上自己需要的字段，但不要放pwd、用户个人隐私数据之类的敏感数据
	Uid int64
}
