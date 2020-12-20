package main

import (
	"encoding/json"
	gsc "gsf/gsc"

	log "github.com/sirupsen/logrus"
)

// ffmpeg cmd option
type ffOpt struct {
	// input file
	Input string
	// output file
	Output string
	// ffmpeg cmmmon option
	CmdOpt []string
}

func main() {
	srcPath := "./res/video/"
	destPath := "./res/out/"

	srcFile := "test.flv"
	destFile := "test.mp4"

	// transcode
	para := []string{
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

	cmdOpt, _ := json.Marshal(ffOpt{srcPath + srcFile, destPath + destFile, para})
	log.Printf("\n inFile-> %v%v\n outFile-> %v%v\n para-> %v\n cmdOpt-> %s\n\n",
		srcPath, srcFile, destPath, destFile, para, cmdOpt)
	err := gsc.Run(srcPath, srcFile, destPath, destFile, string(cmdOpt[:]))
	if err != nil {
		log.Error("gsc run err")
	}
	// gsc.DelFile(dest)
}
