package main

import (
	"gsf/ffmpeg"

	log "github.com/sirupsen/logrus"
)

func main() {

	// str := "ffmpeg -itest.mp4 -c copy -o out.flv"
	// cmd, err := cmd.ParseCommandLine(str)
	// if err != nil {
	// 	fmt.Println("parse cmd err.")
	// }
	// // for k, v := range cmd {
	// // 	fmt.Println(k, v)
	// // }
	// fmt.Println(JsonFormat(cmd))

	res := "./res/video/test.flv"
	dest := "./res/out/test.mp4"
	cmdString := []byte(`{"raw":["-vf scale=-2:960","-c:v libx264","-profile:v main","-level:v 3.1","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 1000k","-maxrate 1000k","-bufsize 1000k","-b:v 1000k","-y"]}`)
	err := ffmpeg.Transcode(res, dest, string(cmdString[:]))
	if err != nil {
		log.Error("transcode err")
	}
	ffmpeg.DelFile(dest)
}
