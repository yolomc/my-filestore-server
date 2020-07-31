package main

import (
	"log"
	"my-filestore-server/service/file/handler"
	"my-filestore-server/service/file/proto"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func main() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.file"),
	)
	service.Init()

	proto.RegisterFileServiceHandler(service.Server(), new(handler.File))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
