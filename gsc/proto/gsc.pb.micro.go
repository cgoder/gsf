// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/gsc.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
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
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Gsc service

func NewGscEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Gsc service

type GscService interface {
	Run(ctx context.Context, in *GscRequest, opts ...client.CallOption) (*GscResponse, error)
}

type gscService struct {
	c    client.Client
	name string
}

func NewGscService(name string, c client.Client) GscService {
	return &gscService{
		c:    c,
		name: name,
	}
}

func (c *gscService) Run(ctx context.Context, in *GscRequest, opts ...client.CallOption) (*GscResponse, error) {
	req := c.c.NewRequest(c.name, "Gsc.Run", in)
	out := new(GscResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Gsc service

type GscHandler interface {
	Run(context.Context, *GscRequest, *GscResponse) error
}

func RegisterGscHandler(s server.Server, hdlr GscHandler, opts ...server.HandlerOption) error {
	type gsc interface {
		Run(ctx context.Context, in *GscRequest, out *GscResponse) error
	}
	type Gsc struct {
		gsc
	}
	h := &gscHandler{hdlr}
	return s.Handle(s.NewHandler(&Gsc{h}, opts...))
}

type gscHandler struct {
	GscHandler
}

func (h *gscHandler) Run(ctx context.Context, in *GscRequest, out *GscResponse) error {
	return h.GscHandler.Run(ctx, in, out)
}