package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"webook_Rouge/internal/web"
)

func main() {
	server := gin.Default()

	// 解决跨域问题
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"}, //
		// AllowMethods: []string{"POST", "GET"}, // 不写的话默认的几个简单方法都ok
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// ExposeHeaders:    []string{"Content-Length"},
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

	u := &web.UserHandler{}

	u.RegisterRoutes(server)

	server.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
