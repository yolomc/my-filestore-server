package main

import (
	"my-filestore-server/route"
)

func main() {
	route.SetUp().Run(":8080")

	//用户
	//http.HandleFunc("/user/signup", handler.SignupHandler)
	//http.HandleFunc("/user/signin", handler.SigninHandler)
	//http.HandleFunc("/user/info", middleware.HTTPInterceptor(handler.UserInfoHandler))

	//文件上传
	// http.HandleFunc("/file/upload", middleware.HTTPInterceptor(handler.UploadHandler))
	// http.HandleFunc("/file/fastupload", middleware.HTTPInterceptor(handler.TryFastUploadHandler))
	// http.HandleFunc("/file/upload/suc", middleware.HTTPInterceptor(handler.UploadSucHandler))
	// http.HandleFunc("/file/meta", middleware.HTTPInterceptor(handler.GetFileMetaHandler))
	// http.HandleFunc("/file/query", middleware.HTTPInterceptor(handler.FileQueryHandler))
	// http.HandleFunc("/file/download", middleware.HTTPInterceptor(handler.DownloadHandler))
	// http.HandleFunc("/file/update", middleware.HTTPInterceptor(handler.FileMetaUpdateHandler))
	// http.HandleFunc("/file/delete", middleware.HTTPInterceptor(handler.FileDeleteHandler))

	//文件分块上传
	// http.HandleFunc("/file/mpupload/init",
	// 	middleware.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	// http.HandleFunc("/file/mpupload/uppart",
	// 	middleware.HTTPInterceptor(handler.UploadPartHandler))
	// http.HandleFunc("/file/mpupload/complete",
	// 	middleware.HTTPInterceptor(handler.CompleteUploadHandler))

	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	fmt.Printf("Failed to start server: %s", err.Error())
	// }

}
