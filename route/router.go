package route

import (
	"my-filestore-server/handler"
	"my-filestore-server/middleware"

	"github.com/gin-gonic/gin"
)

//SetUp 配置路由
func SetUp() *gin.Engine {
	router := gin.Default()

	//静态资源
	router.Static("/static/", "./static")

	//无需验证就能访问的路由
	router.GET("/user/signup", handler.SignupHandler)
	router.POST("/user/signup", handler.DoSignupHandler)
	router.GET("/user/signin", handler.SigninHandler)
	router.POST("/user/signin", handler.DoSigninHandler)

	//加入中间件，use之后的所有handler都会经过中间件
	router.Use(middleware.HTTPInterceptor())
	{
		router.POST("/user/info", handler.UserInfoHandler)

		//文件上传
		router.GET("/file/upload", handler.UploadHandler)
		router.POST("/file/upload", handler.DoUploadHandler)
		router.POST("/file/fastupload", handler.TryFastUploadHandler)
		router.POST("/file/meta", handler.GetFileMetaHandler)
		router.POST("/file/query", handler.FileQueryHandler)
		// router.GET("/file/download", handler.DownloadHandler)
		// router.POST("/file/update", handler.FileMetaUpdateHandler)
		// router.POST("/file/delete", handler.FileDeleteHandler)

		//文件分块上传
		router.POST("/file/mpupload/init", handler.InitialMultipartUploadHandler)
		router.POST("/file/mpupload/uppart", handler.UploadPartHandler)
		router.POST("/file/mpupload/complete", handler.CompleteUploadHandler)
	}
	return router
}
