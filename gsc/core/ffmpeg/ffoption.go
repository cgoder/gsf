package ffmpeg

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type FFOption struct {
	Input  string
	Output string
	// ffmpeg cmmon option slice. like [-maxrate 1000k]
	CmdSlice []string
	// ffmpeg common option string.like [-hide_banner -progress /dev/stdout -i ./gsc/res/video/test.flv -b:v 1000k -maxrate 1000k -minrate 1000k -c:v libx264 -bufsize 1000k -profile:v main -vf scale=-2:960 -y ./gsc/res/out/test.mp4]
	CmdString string
	// all parsed args from CmdString. like []
	arguments map[string]string
}

// GetArgument
func (opts *FFOption) GetArgument(arg string) (string, bool) {
	// if opts.CmdSlice == nil {
	// 	log.Error("FFOption.CmdSlice is nil!")
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

func (opts *FFOption) CmdString2Slice() bool {

	if opts.arguments == nil {
		if opts.CmdString != "" {
			// log.Printf("CmdString--- %s", opts.CmdString)
			if err := json.Unmarshal([]byte(opts.CmdString), &opts.arguments); err != nil {
				log.Error("cmd args encode to map fail. err: ", err)

				return false
			}
		} else {
			log.Error("FFOption.CmdString is null!")
			return false
		}

		// log.Info(common.JsonFormat(opts.arguments))
	}
	return true
}

// GetStrArguments ...
func (opts *FFOption) GetStrArguments() map[string]string {
	return opts.arguments
}
