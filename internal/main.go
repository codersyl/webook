package main

import (
	"github.com/gin-gonic/gin"
	"webook_Rouge/internal/web"
)

func main() {
	server := gin.Default()
	// 静态路由
	u := &web.UserHandler{}

	u.RegisterRoutes(server)

	server.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
