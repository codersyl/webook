package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	// "github.com/gin-contrib/sessions/cookie"
	// "github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook_Rouge/internal/repository"
	"webook_Rouge/internal/repository/dao"
	"webook_Rouge/internal/service"
	"webook_Rouge/internal/web"
	"webook_Rouge/internal/web/middleware"
)

func main() {
	db := initDB()
	server := InitWebServer()

	u := initUser(db)
	u.RegisterRoutes(server)

	server.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func InitWebServer() *gin.Engine {
	server := gin.Default()

	// 解决跨域问题
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"}, //
		// AllowMethods: []string{"POST", "GET"}, // 不写的话默认的几个简单方法都ok
		AllowHeaders: []string{"Content-Type", "Authorization"}, // 允许你过来的请求带的东西
		// ExposeHeaders 是我发给前端的东西中，允许前端读取的东西
		ExposeHeaders:    []string{"x-jwt-token"}, // 不加这个，前端拿不到token
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// 开发环境判断
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			// 如果需要，加上生产环境判断，例如公司域名
			return false
		},
		MaxAge: 12 * time.Hour,
	}))

	//store := cookie.NewStore([]byte("secret")) // 此处存在 cookie中，其实不太安全
	key1_32 := "iFyeVYqAZPMY2p2Jma6zn22jxbKH6TCI"
	key2_32 := "x07p2PVBF9eE7zZTXRzkuS26sztgQuCh"
	//store := memstore.NewStore([]byte(key1_32), []byte(key2_32))
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", "", []byte(key1_32), []byte(key2_32))
	if err != nil {
		panic(err)
	}

	server.Use(sessions.Sessions("webook_session", store))

	// login := middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/login").IgnorePaths("/users/signup")
	login := middleware.NewLoginJWTMiddlewareBuilder().IgnorePaths("/users/login").IgnorePaths("/users/signup")
	server.Use(login.CheckLogin())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		// 只在初始化过程panic
		panic("failed to connect database")
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}
