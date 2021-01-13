package gsc

import (
	"context"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/cgoder/gsc/common"
	"github.com/cgoder/gsc/core/ffmpeg"

	log "github.com/sirupsen/logrus"
)

var (
	processUpdateInterval = time.Millisecond * 500
)

func parseOptions(opt GscOptions) ffmpeg.FFOption {
	var paraStr string
	ffopt := ffmpeg.FFOption{}
	ffopt.Input = opt.Input
	ffopt.Output = opt.Output
	// log.Info("GscOptions--->   ", opt)

	// Parse FFOption.CmdString & FFOption.CmdSlice
	if opt.Opts != (Options{}) {
		////// CmdString   string
		paraS, _ := json.Marshal(opt.Opts)
		// log.Printf("paraS--->%s   ", paraS)
		paraStr = string(paraS[:])
		ffopt.CmdString = paraStr

		if opt.OptSlice == nil {
			///// CmdSlice []string
			value := reflect.ValueOf(opt.Opts)
			types := reflect.TypeOf(opt.Opts)
			var argSlice []string
			for i := 0; i < value.NumField(); i++ {
				if value.Field(i).String() != "" {
					argSlice = append(argSlice, types.Field(i).Tag.Get("json"), value.Field(i).String())
					// log.Printf("Field %v: %v\n", types.Field(i).Tag.Get("json"), value.Field(i))
				}
			}
			// log.Printf("args--->%s   ", argSlice)
			ffopt.CmdSlice = argSlice
			// copy(ffopt.CmdSlice, argSlice)
		} else {
			ffopt.CmdSlice = opt.OptSlice
			// log.Printf("OptSlice--->%s   ", opt.OptSlice)
			// copy(ffopt.CmdSlice, opt.OptSlice)
		}
	} else {
		if opt.OptSlice != nil {
			// translate OptSlice to CmdString
			paraStr = strings.Join(opt.OptSlice, " ")
			ffopt.CmdString = paraStr
			// paraStr, _ := json.Marshal(opt.OptSlice)
			// log.Printf("paraStr--->%s   ", paraStr)
			// ffopt.CmdString = string(paraStr[:])

			//TODO
			ffopt.CmdSlice = opt.OptSlice
			// copy(ffopt.CmdSlice, opt.OptSlice)
		}
	}
	log.Printf("CmdString--->%s   ", ffopt.CmdString)
	log.Printf("CmdSlice--->%s   ", ffopt.CmdSlice)

	ffopt.CmdString2Slice()

	// log.Info("ffopt--->   ", ffopt)
	return ffopt
}

// Run
// Profile:
// ('webm_vp9_720p_3000', '{"raw":["-vf scale=-2:720","-c:v libvpx-vp9","-level:v 4.0","-b:v 3000k","-pix_fmt yuv420p","-f webm","-y"]}', true, 'webm_vp9_3000_720.webm');
// ('h264_main_1080p_6000', '{"raw":["-vf scale=-2:1080","-c:v libx264","-profile:v main","-level:v 4.2","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 6000k","-maxrate 6000k","-bufsize 6000k","-b:v 6000k","-y"]}');
// ('h264_main_720p_3000', '{"raw":["-vf scale=-2:720","-c:v libx264","-profile:v main","-level:v 4.0","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 3000k","-maxrate 3000k","-bufsize 3000k","-b:v 3000k","-y"]}');
// ('h264_main_480p_1000', '{"raw":["-vf scale=-2:480","-c:v libx264","-profile:v main","-level:v 3.1","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 1000k","-maxrate 1000k","-bufsize 1000k","-b:v 1000k","-y"]}');
// ('h264_baseline_360p_600', '{"raw":["-vf scale=-2:360","-c:v libx264","-profile:v baseline","-level:v 3.0","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 600k","-maxrate 600k","-bufsize 600k","-b:v 600k","-y"]}');

func Run(opt GscOptions) error {
	ffopt := parseOptions(opt)
	// log.Info("ffopt--->   ", ffopt)

	if exist, err := common.FilePathExists(ffopt.Input); !exist {
		log.Error("src file not exist err: ", err)
		return err
	}

	if ffopt.Output != "" {
		dirname, _ := filepath.Split(ffopt.Output)
		if exist, _ := common.FilePathExists(dirname); !exist {
			log.Info("dst file not exist, create it... ", dirname)
			if err := common.FilePathCreate(dirname, ""); err != nil {
				log.Error("create dst path err: ", err)
				return err
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// probe in file
	// ffprobe := ffmpeg.FFProbe{}
	// if probeData, err := ffprobe.Execute(ffopt.Input, ""); err == nil {
	// 	log.Info("Src file info ---> ", common.JsonFormat(probeData))
	// } else {
	// 	log.Error("ffprobe execute err: ", err)
	// 	return err
	// }

	ffmpeg := ffmpeg.FFmpeg{}

	// progress
	go runProgress(ctx, "TranscodeTask", 1, &ffmpeg)

	// Exec
	// log.Info(ffmpeg.Version())
	if err := ffmpeg.Execute(ctx, ffopt); err != nil {
		log.Error("ffmpeg execute err: ", err)
		return err
	}
	// log.Info(Common.JsonFormat(ffmpeg.Status))

	// probe out file
	// if probeData, err := ffprobe.Execute(ffopt.Input, ""); err == nil {
	// 	log.Info("Dst file info ---> ", common.JsonFormat(probeData))
	// } else {
	// 	log.Error("ffprobe execute err: ", err)
	// 	return err
	// }

	return nil
}

func runProgress(ctx context.Context, guid string, encodeID int64, f *ffmpeg.FFmpeg) {

	ticker := time.NewTicker(processUpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info(common.JsonFormat(f.Status))
			log.Info("transcode done.")
			return
		case <-ticker.C:
		default:
			// log.Info(common.JsonFormat(f.Status))
		}
	}
}
