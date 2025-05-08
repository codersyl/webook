package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

// JWT 登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要校验登录态的请求
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 使用JWT来登录校验

		//tokenHeader := ctx.Request.Header.Get("Authorization")
		tokenHeader := ctx.GetHeader("Authorization")
		fmt.Println(tokenHeader)
		if tokenHeader == "" {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 分段长度非法，有人乱传 "Authorization"
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := segs[1]
		key32_ForToken := "iFyeVYqAZPMY2p2Jma6zn22jxbKH6TCI" // 应该与当时加密的key一致
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(key32_ForToken), nil
		})

		if err != nil {
			// token解析错误
			// 按理应该是系统内部错误 500
			// 但可能是攻击者发送过来的字符串，所以按照没登录处理
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		fmt.Println(">> 解析到 token: ", token)
	}
}
