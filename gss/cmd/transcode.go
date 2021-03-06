package cmd

import (
	"context"

	proto "github.com/cgoder/gsc/proto"
	micro "github.com/micro/go-micro/v2"
	log "github.com/sirupsen/logrus"
)

func Transcode() {
	service := micro.NewService(
		micro.Name("service.gsf.gss"),
	)
	service.Init()
	gsc := proto.NewGscService("service.gsf.gsc", service.Client())

	srcPath := "./res/video/"
	destPath := "./res/out/"

	srcFile := "test.flv"
	destFile := "test.mp4"

	// opts := gsc.Options{
	// 	VideoFilter:     "scale=-2:960",
	// 	VideoCodec:      "libx264",
	// 	VideoProfile:    "main",
	// 	VideoMinBitrate: "1000k",
	// 	VideoMaxBitRate: "1000k",
	// 	BufferSize:      "1000k",
	// 	VideoBitRate:    "1000k",
	// 	Overwrite:       "-y",
	// }
	// gscOpts := gsc.GscOptions{Input: srcPath + srcFile, Output: destPath + destFile, Opts: opts}

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
	gscOpts := proto.GscRequest{Input: srcPath + srcFile, Output: destPath + destFile, OptSlice: optslice}

	resq, err := gsc.Run(context.TODO(), &gscOpts)
	log.Info("transcode exec done.", resq)
	if err != nil {
		log.Error("gsc run err: ", err)
	}

	// log.Info(common.JsonFormat(gsc.FFProbe))
	// gsc.DelFile(destPath + destFile)

}
