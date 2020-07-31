package middleware

import (
	"my-filestore-server/redis"
	"net/http"

	"github.com/gin-gonic/gin"
)

//HTTPInterceptor http请求拦截器
func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		if len(username) < 3 || !isTokenValid(username, token) {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"msg":  "token validate failed",
				"code": -2,
			})
			return
		}
		c.Next()
	}
}

//IsTokenValid 验证token
func isTokenValid(username string, token string) bool {
	t, _ := redis.Get("token_" + username)
	return t == token
}
