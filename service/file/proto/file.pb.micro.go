// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: file.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for FileService service

type FileService interface {
	FileQuery(ctx context.Context, in *ReqFileQuery, opts ...client.CallOption) (*RespFileQuery, error)
}

type fileService struct {
	c    client.Client
	name string
}

func NewFileService(name string, c client.Client) FileService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "proto"
	}
	return &fileService{
		c:    c,
		name: name,
	}
}

func (c *fileService) FileQuery(ctx context.Context, in *ReqFileQuery, opts ...client.CallOption) (*RespFileQuery, error) {
	req := c.c.NewRequest(c.name, "FileService.FileQuery", in)
	out := new(RespFileQuery)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for FileService service

type FileServiceHandler interface {
	FileQuery(context.Context, *ReqFileQuery, *RespFileQuery) error
}

func RegisterFileServiceHandler(s server.Server, hdlr FileServiceHandler, opts ...server.HandlerOption) error {
	type fileService interface {
		FileQuery(ctx context.Context, in *ReqFileQuery, out *RespFileQuery) error
	}
	type FileService struct {
		fileService
	}
	h := &fileServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&FileService{h}, opts...))
}

type fileServiceHandler struct {
	FileServiceHandler
}

func (h *fileServiceHandler) FileQuery(ctx context.Context, in *ReqFileQuery, out *RespFileQuery) error {
	return h.FileServiceHandler.FileQuery(ctx, in, out)
}
