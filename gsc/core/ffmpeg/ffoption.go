package ffmpeg

import (
	"encoding/json"

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
	// if opts.CmdSlice == nil {
	// 	log.Error("FfOption.CmdSlice is nil!")
	// 	return "", false
	// }
	// // cmd := strings.Join(opts.CmdSlice, " ")

	opts.CmdString2Slice()

	if v, ok := opts.arguments[arg]; ok {
		return v, true
	}

	log.Error("can not find arg-> .", arg)
	return "", false
}

func (opts *FfOption) CmdString2Slice() bool {

	if opts.arguments == nil {
		if opts.CmdString != "" {
			// log.Printf("CmdString--- %s", opts.CmdString)
			if err := json.Unmarshal([]byte(opts.CmdString), &opts.arguments); err != nil {
				log.Error("cmd args encode to map fail")
				return false
			}
		} else {
			log.Error("FfOption.CmdString is null!")
			return false
		}

		// log.Info(common.JsonFormat(opts.arguments))
	}
	return true
}

// GetStrArguments ...
func (opts *FfOption) GetStrArguments() map[string]string {
	return opts.arguments
}
