package gsc

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"time"

	"github.com/cgoder/gsc/common"
	"github.com/cgoder/gsc/core/ffmpeg"

	log "github.com/sirupsen/logrus"
)

var (
	processUpdateInterval = time.Millisecond * 500
)

func parseArgs(cmdArgs ffmpeg.CmdArgs) ffmpeg.Options {
	///// CmdSlice []string
	value := reflect.ValueOf(cmdArgs)
	types := reflect.TypeOf(cmdArgs)
	var args []string
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).String() != "" {
			args = append(args, types.Field(i).Tag.Get("json"), value.Field(i).String())
			log.Printf("Field %v: %v\n", types.Field(i).Tag.Get("json"), value.Field(i))
		}
	}
	log.Printf("args--->%s   ", args)

	////// CmdOpt   string
	paraS, _ := json.Marshal(cmdArgs)
	// log.Printf("paraS--->%s   ", paraS)
	paraStr := string(paraS[:])
	log.Info("paraStr--->   ", paraStr)

	ffopt := ffmpeg.Options{CmdOpt: paraStr, CmdSlice: args}
	log.Info("ffopt--->   ", ffopt)
	return ffopt
}

// Run
// Profile:
// ('webm_vp9_720p_3000', '{"raw":["-vf scale=-2:720","-c:v libvpx-vp9","-level:v 4.0","-b:v 3000k","-pix_fmt yuv420p","-f webm","-y"]}', true, 'webm_vp9_3000_720.webm');
// ('h264_main_1080p_6000', '{"raw":["-vf scale=-2:1080","-c:v libx264","-profile:v main","-level:v 4.2","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 6000k","-maxrate 6000k","-bufsize 6000k","-b:v 6000k","-y"]}');
// ('h264_main_720p_3000', '{"raw":["-vf scale=-2:720","-c:v libx264","-profile:v main","-level:v 4.0","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 3000k","-maxrate 3000k","-bufsize 3000k","-b:v 3000k","-y"]}');
// ('h264_main_480p_1000', '{"raw":["-vf scale=-2:480","-c:v libx264","-profile:v main","-level:v 3.1","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 1000k","-maxrate 1000k","-bufsize 1000k","-b:v 1000k","-y"]}');
// ('h264_baseline_360p_600', '{"raw":["-vf scale=-2:360","-c:v libx264","-profile:v baseline","-level:v 3.0","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 600k","-maxrate 600k","-bufsize 600k","-b:v 600k","-y"]}');

func Run(args ffmpeg.CmdArgs) error {
	ffopt := parseArgs(args)

	// if inURI/outURI not exist.
	src, ok := ffopt.GetArgument("-i")
	if !ok {
		err := errors.New("ffmpeg get args error")
		return err
	}
	_, err := os.Stat(src)
	if err != nil {
		log.Error(err)
		return err
	}
	// log.Info(src)
	// log.Info("fileinfo: ", fileInfo.IsDir(), fileInfo.Name(), fileInfo.Size())

	// _, err = os.Stat(destPath)
	// if err != nil && os.IsNotExist(err) {
	// 	err = CreateLocalPath(destPath, "")
	// 	if err != nil {
	// 		log.Error("create path error! ", destPath)
	// 		return err
	// 	}
	// }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// probe in file
	ffprobe := ffmpeg.FFProbe{}
	ffmpeg := ffmpeg.FFmpeg{}
	probeData := ffprobe.Execute(src, "")
	if err != nil {
		log.Error("ffprobe execute err")
		return err
	}
	// log.Info("Source file info ---> ", common.JsonFormat(probeData))

	// progress
	go runProgress(ctx, "TranscodeTask", 1, probeData, &ffmpeg)

	// Exec
	// log.Info(ffmpeg.Version())
	err = ffmpeg.Execute(ctx, ffopt)
	if err != nil {
		log.Error("ffmpeg execute err")
		return err
	}
	// log.Info(Common.JsonFormat(ffmpeg.Status))
	// probe out file
	// probeData = ffprobe.Execute(destPath + destFile)
	// if err != nil {
	// 	log.Error("ffprobe execute err")
	// 	return err
	// }
	// log.Info(Common.JsonFormat(probeData))

	return nil
}

func runProgress(ctx context.Context, guid string, encodeID int64, p *ffmpeg.Metadata, f *ffmpeg.FFmpeg) {

	ticker := time.NewTicker(processUpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info(common.JsonFormat(f.Status))
			log.Info("transcode done.")
			return
		case <-ticker.C:
			// default:
			log.Info(common.JsonFormat(f.Status))
		}
	}
}

func CreateLocalPath(dirPath string, GUID string) error {
	// Get local destination path.
	var tmpDir string
	if GUID == "" {
		tmpDir = dirPath
	} else {
		tmpDir = dirPath + "/" + GUID
	}

	err := os.MkdirAll(tmpDir, 0700)
	if err != nil {
		return err
	}
	return nil
}

func DelFile(filePath string) error {
	log.Info("del file: ", filePath)

	err := os.RemoveAll(filePath)
	if err != nil {
		return err
	}
	return nil
}
