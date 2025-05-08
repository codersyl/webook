package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要校验登录态的请求
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil { // 未登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		now := time.Now().UnixMilli()
		if updateTime == nil { // 登录后，没刷新过
			sess.Set("update_time", now)
			sess.Options(sessions.Options{
				MaxAge: 30 * 60,
			})
			sess.Save()
			return
		}

		// 有 update_time
		updateTimeVal, ok := updateTime.(int64) // 存了int64，应该能转换回去
		if !ok {
			// 非法的刷新时间，说明发过来的请求有问题
			// 按理说其实把请求端的用户状态登出比较合理
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if now-updateTimeVal > 5*1000 { // 超过一分钟，刷新
			sess.Set("update_time", now)
			sess.Options(sessions.Options{
				MaxAge: 30 * 60,
			})
			sess.Save()
			return
		}
	}
}
