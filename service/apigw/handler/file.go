package handler

import (
	"context"
	"log"
	"my-filestore-server/service/file/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

var (
	fileCli proto.FileService
)

func init() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(micro.Registry(reg))
	//初始化service，解析命令行参数等
	service.Init()
	//初始化一个rpcClient
	fileCli = proto.NewFileService("go.micro.service.file", service.Client())
}

//FileQueryHandler 批量获取用户文件信息
func FileQueryHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	limit, _ := strconv.Atoi(c.Request.FormValue("limit"))

	resp, err := fileCli.FileQuery(context.TODO(), &proto.ReqFileQuery{
		Username: username,
		Limit:    int32(limit),
	})
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "Application/json", resp.Data)
}
