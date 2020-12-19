package gsc

import (
	"context"
	"math"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	processUpdateInterval = time.Second * 1
)

// Run
// Profile:
// ('webm_vp9_720p_3000', '{"raw":["-vf scale=-2:720","-c:v libvpx-vp9","-level:v 4.0","-b:v 3000k","-pix_fmt yuv420p","-f webm","-y"]}', true, 'webm_vp9_3000_720.webm');
// ('h264_main_1080p_6000', '{"raw":["-vf scale=-2:1080","-c:v libx264","-profile:v main","-level:v 4.2","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 6000k","-maxrate 6000k","-bufsize 6000k","-b:v 6000k","-y"]}');
// ('h264_main_720p_3000', '{"raw":["-vf scale=-2:720","-c:v libx264","-profile:v main","-level:v 4.0","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 3000k","-maxrate 3000k","-bufsize 3000k","-b:v 3000k","-y"]}');
// ('h264_main_480p_1000', '{"raw":["-vf scale=-2:480","-c:v libx264","-profile:v main","-level:v 3.1","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 1000k","-maxrate 1000k","-bufsize 1000k","-b:v 1000k","-y"]}');
// ('h264_baseline_360p_600', '{"raw":["-vf scale=-2:360","-c:v libx264","-profile:v baseline","-level:v 3.0","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 600k","-maxrate 600k","-bufsize 600k","-b:v 600k","-y"]}');

func Run(srcPath, srcFile, destPath, destFile string, cmdOpt string) error {

	// if inURI/outURI not exist.
	_, err := os.Stat(srcPath + srcFile)
	if err != nil {
		log.Error("not a file! ", srcPath+srcFile)
		return err
	}
	// log.Info("fileinfo: ", fileInfo.IsDir(), fileInfo.Name(), fileInfo.Size())

	_, err = os.Stat(destPath)
	if err != nil && os.IsNotExist(err) {
		err = CreateLocalPath(destPath, "")
		if err != nil {
			log.Error("create path error! ", destPath)
			return err
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// probe in file
	ffmpeg := FFmpeg{}
	ffprobe := FFProbe{}
	probeData := ffprobe.Execute(srcPath + srcFile)
	if err != nil {
		log.Error("ffprobe execute err")
		return err
	}
	// log.Info(cmd.JsonFormat(probeData))

	// progress
	go transcodeProgress(ctx, "TranscodeTask", 1, probeData, &ffmpeg)

	// Exec
	// log.Info(ffmpeg.Version())
	err = ffmpeg.Execute(ctx, srcPath+srcFile, destPath+destFile, cmdOpt)
	if err != nil {
		log.Error("ffmpeg execute err")
		return err
	}
	// log.Info(cmd.JsonFormat(ffmpeg.Progress))
	// probe out file
	// probeData = ffprobe.Execute(destPath + destFile)
	// if err != nil {
	// 	log.Error("ffprobe execute err")
	// 	return err
	// }
	// log.Info(cmd.JsonFormat(probeData))

	return nil
}

func transcodeProgress(ctx context.Context, guid string, encodeID int64, p *FFProbeResponse, f *FFmpeg) {

	ticker := time.NewTicker(processUpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("transcode done.")
			return
		case <-ticker.C:
			// log.Info(cmd.JsonFormat(f.Progress))

			// Only track progress if we know the total frames.
			totalFrames, _ := strconv.Atoi(p.Streams[0].NbFrames)
			if totalFrames != 0 {
				currentFrame := f.Progress.Frame
				speed := f.Progress.Speed
				fps := f.Progress.FPS
				pct := (float64(currentFrame) / float64(totalFrames)) * 100

				// Update progress.
				pct = math.Round(pct*100) / 100
				log.Info("progress: %d / %d - %0.2f%%, %d, %s", currentFrame, totalFrames, pct, fps, speed)
			}
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
