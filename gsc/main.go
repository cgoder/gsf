package main

import (
	"context"
	"errors"
	"fmt"

	core "github.com/cgoder/gsc/core"
	pb "github.com/cgoder/gsc/proto"
	"github.com/micro/go-micro/v2"
)

type Gsc struct{}

func (gsc *Gsc) Run(ctx context.Context, gscReq *pb.GscRequest, gscRes *pb.GscResponse) error {
	if gscReq == nil {
		return errors.New("request is nil!")
	}

	gscOpts := core.GscOptions{Input: gscReq.Input, Output: gscReq.Output, OptSlice: gscReq.OptSlice}
	err := core.Run(gscOpts)
	if err != nil {
		fmt.Println("gsc run err")
		gscRes.Result = "gsc run err!"
	}
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("service.gsf.gsc"),
	)
	service.Init()

	pb.RegisterGscHandler(service.Server(), new(Gsc))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
