package cmd

import (
	"context"

	proto "github.com/cgoder/gsc/proto"
	micro "github.com/micro/go-micro/v2"
	log "github.com/sirupsen/logrus"
)

func Split() {
	service := micro.NewService(
		micro.Name("service.gsf.gss"),
	)
	service.Init()
	gsc := proto.NewGscService("service.gsf.gsc", service.Client())

	srcPath := "./res/video/"
	destPath := "./res/out/"

	srcFile := "test.flv"
	destFile := "%d.mp4"

	optslice := []string{
		"-c", "copy",
		"-f", "segment",
		"-segment_time", "5",
		"-reset_timestamps", "1",
		"-map", "0:0",
		"-map", "0:1",
		"-y",
	}
	gscOpts := proto.GscRequest{Input: srcPath + srcFile, Output: destPath + destFile, OptSlice: optslice}

	resq, err := gsc.Run(context.TODO(), &gscOpts)
	log.Info("split exec done.", resq)
	if err != nil {
		log.Error("gsc run err: ", err)
	}

	// log.Info(common.JsonFormat(gsc.FFProbe))
	// gsc.DelFile(destPath + destFile)

}
