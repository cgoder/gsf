package cmd

import (
	"context"

	proto "github.com/cgoder/gsc/proto"
	micro "github.com/micro/go-micro/v2"
	log "github.com/sirupsen/logrus"
)

func Remux() {
	service := micro.NewService(
		micro.Name("service.gsf.gss"),
	)
	service.Init()
	gsc := proto.NewGscService("service.gsf.gsc", service.Client())

	srcPath := "./res/video/"
	destPath := "./res/out/"

	srcFile := "test.flv"
	destFile := "test.mp4"

	optslice := []string{
		"-c", "copy",
		"-y",
	}
	gscOpts := proto.GscRequest{Input: srcPath + srcFile, Output: destPath + destFile, OptSlice: optslice}

	resq, err := gsc.Run(context.TODO(), &gscOpts)
	log.Info("remux exec done.", resq)
	if err != nil {
		log.Error("gsc run err: ", err)
	}

	// log.Info(common.JsonFormat(gsc.FFProbe))
	// gsc.DelFile(destPath + destFile)

}
