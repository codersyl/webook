package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	// 静态路由
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 路径参数
	server.GET("/usernom/:name", func(context *gin.Context) {
		name := context.Param("name")
		context.String(http.StatusOK, "你传给我的名字是： %s", name)
	})

	// 查询参数
	server.GET("/query_test", func(context *gin.Context) {
		name := context.Query("name")
		response := ""
		if len(name) == 0 {
			response += fmt.Sprintf("你没给我传递name\n")

		} else {
			response += fmt.Sprintf("你传给我的名字是： %s\n", name)

		}
		response += "接下来显示你传给我的参数列表：\n"
		m := context.Request.URL.Query()
		for k, v := range m {
			response += fmt.Sprintf("%s 的值是 %s\n", k, v)
		}
		context.String(http.StatusOK, response)
	})

	// 通配符路由
	server.GET("pre/*middle", func(context *gin.Context) {
		middle := context.Param("middle")
		context.String(http.StatusOK, "匹配成功，你输入的 middle 段是 %s", middle)
	})

	server.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
