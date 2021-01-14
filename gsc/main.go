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

	optslice := []string{
		"-vf", "scale=-2:960",
		"-c:v", "libx264",
		"-profile:v", "main",
		"-level:v", "3.1",
		"-x264opts", "scenecut=0:open_gop=0:min-keyint=72:keyint=72",
		"-minrate", "1000k",
		"-maxrate", "1000k",
		"-bufsize", "1000k",
		"-b:v", "1000k",
		"-y",
	}
	opts := gscReq.OptSlice
	if opts == nil {
		opts = optslice
	}

	gscOpts := core.GscOptions{Input: gscReq.Input, Output: gscReq.Output, OptSlice: opts}
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
