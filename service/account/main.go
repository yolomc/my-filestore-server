package main

import (
	"log"
	"my-filestore-server/service/account/handler"
	"my-filestore-server/service/account/proto"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func main() {
	// go mod edit -require=google.golang.org/grpc@v1.26.0
	// go get -u -x google.golang.org/grpc@v1.26.0

	// 修改consul地址，如果是本机，这段代码和后面的那行使用代码都是可以不用的
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"127.0.0.1:8500"}
	})

	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.user"),
	)
	service.Init()

	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
