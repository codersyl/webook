package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"webook_Rouge/internal/web"
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

		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := tokenHeader
		// fmt.Println("鉴权用户一位，token：  ", tokenStr, " ENDtoken")
		claims := &web.UserClaims{}                          // 指针，因为Parse的时候会放数据进来
		key32_ForToken := "iFyeVYqAZPMY2p2Jma6zn22jxbKH6TCI" // 应该与当时加密的key一致
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
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

		if claims.UserAgent != ctx.Request.UserAgent() {
			// 存在安全问题，需要监控，并将该用户登出
			// 此处直接失败
			// ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.String(http.StatusUnauthorized, "你小子换设备了，重新登录\n")
			return
		}

		ctx.Set("claims", claims) // 可以只加Uid，但是后续可能存储更多数据，所以存个总的
	}
}
