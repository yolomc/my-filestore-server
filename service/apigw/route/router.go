package route

import (
	"my-filestore-server/middleware"
	"my-filestore-server/service/apigw/handler"

	"github.com/gin-gonic/gin"
)

//SetUp 配置路由
func SetUp() *gin.Engine {
	router := gin.Default()

	router.Static("/static/", "./static")

	router.GET("/user/signup", handler.SignupHandler)
	router.POST("/user/signup", handler.DoSignupHandler)
	router.GET("/user/signin", handler.SigninHandler)
	router.POST("/user/signin", handler.DoSigninHandler)

	router.Use(middleware.HTTPInterceptor())
	{
		router.POST("/user/info", handler.UserInfoHandler)

		router.POST("/file/query", handler.FileQueryHandler)
	}
	return router

}
