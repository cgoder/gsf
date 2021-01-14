package main

import (
	"context"
	"fmt"

	pb "github.com/cgoder/gsc/proto"
	"github.com/micro/go-micro/v2"
)

type Gsc struct{}

func (gsc *Gsc) Run(ctx context.Context, gscReq *pb.GscRequest, gscRes *pb.GscResponse) error {
	gscRes.Result = "it is ok!"
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
