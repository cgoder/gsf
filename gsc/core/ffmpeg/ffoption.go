package ffmpeg

import (
	"encoding/json"

	"github.com/cgoder/gsc/common"

	log "github.com/sirupsen/logrus"
)

type FfOption struct {
	// ffmpeg cmmon option slice. like [-maxrate 1000k]
	CmdSlice []string
	// ffmpeg common option string.like [-hide_banner -progress /dev/stdout -i ./gsc/res/video/test.flv -b:v 1000k -maxrate 1000k -minrate 1000k -c:v libx264 -bufsize 1000k -profile:v main -vf scale=-2:960 -y ./gsc/res/out/test.mp4]
	CmdString string
	// all parsed args from CmdString. like []
	arguments map[string]string
}

// GetArgument
func (opts *FfOption) GetArgument(arg string) (string, bool) {
	if opts.CmdSlice == nil {
		return "", false
	}
	// cmd := strings.Join(opts.CmdSlice, " ")
	// log.Info(common.JsonFormat(cmd))
	if opts.arguments == nil {
		if err := json.Unmarshal([]byte(opts.CmdString), &opts.arguments); err != nil {
			log.Error("cmd args encode to map fail")
			return "", false
		}

		// if opts.CmdString == "" {
		// 	log.Error("parse args first!")
		// 	return "", false
		// } else {
		// 	opts.SetStrArguments(opts.CmdSlice)
		// }
	}
	log.Info(common.JsonFormat(opts.arguments))

	if v, ok := opts.arguments[arg]; ok {
		return v, true
	}

	log.Error("can not find arg-> .", arg)
	return "", false
}

// GetStrArguments ...
func (opts *FfOption) GetStrArguments() map[string]string {
	return opts.arguments
}
