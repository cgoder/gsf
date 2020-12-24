package main

import (
	"github.com/cgoder/gsc"
	"github.com/cgoder/gsc/core"

	log "github.com/sirupsen/logrus"
)

// ffmpeg cmd option

func main() {
	srcPath := "./gsc/res/video/"
	destPath := "./gsc/res/out/"

	srcFile := "test.flv"
	destFile := "test.mp4"

	// transcode
	// para := []string{
	// 	"-vf", "scale=-2:960",
	// 	"-c:v", "libx264",
	// 	"-profile:v", "main",
	// 	"-level:v", "3.1",
	// 	"-x264opts", "scenecut=0:open_gop=0:min-keyint=72:keyint=72",
	// 	"-minrate", "1000k",
	// 	"-maxrate", "1000k",
	// 	"-bufsize", "1000k",
	// 	"-b:v", "1000k",
	// 	"-y",
	// 	"-i", srcPath + srcFile,
	// 	destPath + destFile,
	// }

	// para := []string{
	// 	"-vf scale=-2:960",
	// 	"-c:v libx264",
	// 	"-profile:v main",
	// 	"-level:v 3.1",
	// 	"-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72",
	// 	"-minrate 1000k",
	// 	"-maxrate 1000k",
	// 	"-bufsize 1000k",
	// 	"-b:v 1000k",
	// 	"-y",
	// }

	para := core.Option{
		Input:           srcPath + srcFile,
		VideoFilter:     "scale=-2:960",
		VideoCodec:      "libx264",
		VideoProfile:    "main",
		VideoMinBitrate: "1000k",
		VideoMaxBitRate: "1000k",
		BufferSize:      "1000k",
		VideoBitRate:    "1000k",
		Overwrite:       destPath + destFile,
	}

	// remux
	// para := []string{
	// 	"-c", "copy",
	// }

	// split
	// para := []string{
	// 	"-c", "copy",
	// 	"-f", "segment",
	// 	"-segment_time", "5",
	// 	"-reset_timestamps", "1",
	// 	"-map", "0:0",
	// 	"-map", "0:1",
	// 	"-y",
	// }
	// destFile = "%d.mp4"

	err := gsc.Run(para)

	if err != nil {
		log.Error("gsc run err")
	}
	// gsc.DelFile(dest)
}
