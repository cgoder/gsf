package cmd

import (
	"context"
	"strconv"

	proto "github.com/cgoder/gsc/proto"
	micro "github.com/micro/go-micro/v2"
	log "github.com/sirupsen/logrus"
)

func SnapShot() {
	service := micro.NewService(
		micro.Name("service.gsf.gss"),
	)
	service.Init()
	gsc := proto.NewGscService("service.gsf.gsc", service.Client())

	sc_time := "5"
	pix_type := "jpg"
	frame_num := 10

	srcPath := "./res/video/"
	destPath := "./res/out/"

	srcFile := "test.flv"
	destFile := "test-%5d" + "." + pix_type

	optslice := []string{
		"-y",
		"-f", "image2",
		"-ss", sc_time,
		"-vframes", strconv.Itoa(frame_num),
	}
	gscOpts := proto.GscRequest{Input: srcPath + srcFile, Output: destPath + destFile, OptSlice: optslice}

	resq, err := gsc.Run(context.TODO(), &gscOpts)
	log.Info("snapshot exec done.", resq)
	if err != nil {
		log.Error("gsc run err: ", err)
	}

	// log.Info(common.JsonFormat(gsc.FFProbe))
	// gsc.DelFile(destPath + destFile)

}
