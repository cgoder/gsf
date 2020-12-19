package gsc

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// openencoder
const (
	ffmpegCmd      = "ffmpeg"
	updateInterval = time.Second * 1
)

// FFmpeg struct.
type FFmpeg struct {
	Progress    progress
	cmd         *exec.Cmd
	isCancelled bool
}

type progress struct {
	Frame      int
	FPS        float64
	Bitrate    float64
	TotalSize  int
	OutTimeMS  int
	OutTime    string
	DupFrames  int
	DropFrames int
	Speed      string
	Progress   float64
}

// ffmpegOptions struct passed into Ffmpeg.Run.
type ffmpegOptions struct {
	Input  string
	Output string
	// 优先parse cmdOpt结构中的参数
	CmdOpt []string `json:"CmdOpt"` // Raw flag options.

	// resave
	Container string       `json:"container"`
	Video     videoOptions `json:"video"`
	Audio     audioOptions `json:"audio"`
}

type videoOptions struct {
	Codec                string `json:"codec"`
	Preset               string `json:"preset"`
	HardwareAcceleration string `json:"hardware_acceleration_option"`
	Pass                 string `json:"pass"`
	Crf                  int    `json:"crf"`
	Bitrate              string `json:"bitrate"`
	MinRate              string `json:"minrate"`
	MaxRate              string `json:"maxrate"`
	BufSize              string `json:"bufsize"`
	PixelFormat          string `json:"pixel_format"`
	FrameRate            string `json:"frame_rate"`
	Speed                string `json:"speed"`
	Tune                 string `json:"tune"`
	Profile              string `json:"profile"`
	Level                string `json:"level"`
}

type audioOptions struct {
	Codec string
}

// Execute runs the ffmpeg encoder with options.
func (f *FFmpeg) Execute(ctx context.Context, input, output string, cmdOpt string) error {

	// Parse options and add to args slice.
	args := parseOptions(input, output, cmdOpt)

	// Execute command.
	log.Info("running FFmpeg with options: ", args)
	f.cmd = exec.Command(ffmpegCmd, args...)

	// Capture stderr (if any).
	var stderr bytes.Buffer
	f.cmd.Stderr = &stderr
	stdout, _ := f.cmd.StdoutPipe()

	// Update status goroutine
	go f.execProgress(ctx, stdout)

	// Run
	f.cmd.Start()

	err := f.cmd.Wait()
	if err != nil {
		if f.isCancelled {
			return errors.New("cancelled")
		}
		return err
	}
	return nil
}

// Cancel stops an FFmpeg job from running.
func (f *FFmpeg) Cancel() {
	log.Warn("killing ffmpeg process")
	f.isCancelled = true
	if err := f.cmd.Process.Kill(); err != nil {
		log.Warn("failed to kill process: ", err)
	}
	log.Warn("killed ffmpeg process")
}

// Version gets the ffmpeg version.
func (f *FFmpeg) Version() string {
	out, _ := exec.Command(ffmpegCmd, "-version").Output()
	return string(out)
}

func (f *FFmpeg) setProgressParts(parts []string) {
	for i := 0; i < len(parts); i++ {
		progressSplit := strings.Split(parts[i], "=")
		k := progressSplit[0]
		v := progressSplit[1]

		switch k {
		case "frame":
			frame, _ := strconv.Atoi(v)
			f.Progress.Frame = frame
		case "fps":
			fps, _ := strconv.ParseFloat(v, 64)
			f.Progress.FPS = fps
		case "bitrate":
			v = strings.Replace(v, "kbits/s", "", -1)
			bitrate, _ := strconv.ParseFloat(v, 64)
			f.Progress.Bitrate = bitrate
		case "total_size":
			size, _ := strconv.Atoi(v)
			f.Progress.TotalSize = size
		case "out_time_ms":
			outTimeMS, _ := strconv.Atoi(v)
			f.Progress.OutTimeMS = outTimeMS
		case "out_time":
			f.Progress.OutTime = v
		case "dup_frames":
			frames, _ := strconv.Atoi(v)
			f.Progress.DupFrames = frames
		case "drop_frames":
			frames, _ := strconv.Atoi(v)
			f.Progress.DropFrames = frames
		case "speed":
			f.Progress.Speed = v
		case "progress":
			if v == "end" {
				// end
				f.Progress.Progress = 1
			} else {
				// continue
				f.Progress.Progress = 0
			}
			// progress, _ := strconv.ParseFloat(v, 64)
			// f.Progress.Progress = progress
		}
	}
	// log.Info(cmd.JsonFormat(f.Progress))
}

func (f *FFmpeg) execProgress(ctx context.Context, stdout io.ReadCloser) {
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	scanner := bufio.NewScanner(stdout)
	for {
		select {
		case <-ctx.Done():
			log.Info("ffmpeg process done.")
			return
		case <-ticker.C:
			for scanner.Scan() {
				line := scanner.Text()
				str := strings.Replace(line, " ", "", -1)

				parts := strings.Split(str, " ")

				// log.Info(parts)
				f.setProgressParts(parts)
			}
			// log.Info(cmd.JsonFormat(f.Progress))
		}
	}
}

// Utilities for parsing ffmpeg options.
func parseOptions(input, output string, cmdOpt string) []string {
	if input == "" || output == "" || cmdOpt == "" {
		return nil
	}

	args := []string{
		"-hide_banner",
		"-loglevel", "error", // Set loglevel to fail job on errors.
		"-progress", "pipe:1",
		"-i", input,
	}

	// Decode JSON get options list from cmdOpt.
	options := &ffmpegOptions{}
	log.Info("parse ffmpeg option -> ", cmdOpt)
	if err := json.Unmarshal([]byte(cmdOpt), &options); err != nil {
		panic(err)
	}
	// If raw options provided, add the list of raw options from ffmpeg presets.
	if len(options.CmdOpt) > 0 {
		for _, v := range options.CmdOpt {
			args = append(args, strings.Split(v, " ")...)
		}
		args = append(args, output)
		return args
	}

	// Set options from struct.
	args = append(args, transformOptions(options)...)

	// Add output arg last.
	args = append(args, output)
	return args
}

// transformOptions converts the ffmpegOptions{} struct and converts into
// a slice of ffmpeg options to be passed to exec.Command arguments.
//
// NOTE: There is probably a better way of iterating the struct fields and values
// using reflect, but there are some tricky ffmpeg options here, such as video filters.
// TODO: Look into refactor using reflect. Example:
//   fields := reflect.TypeOf(opt)
//   values := reflect.ValueOf(opt)
func transformOptions(opt *ffmpegOptions) []string {
	args := []string{}

	// Video codec.
	if opt.Video.Codec != "" {
		arg := []string{"-c:v", opt.Video.Codec}
		args = append(args, arg...)
	}

	// Audio codec.
	if opt.Audio.Codec != "" {
		arg := []string{"-c:a", opt.Audio.Codec}
		args = append(args, arg...)
	}

	// Video preset.
	if opt.Video.Preset != "" && opt.Video.Preset != "none" {
		arg := []string{"-preset", opt.Video.Preset}
		args = append(args, arg...)
	}

	// Hardware Acceleration.
	if opt.Video.HardwareAcceleration == "nvenc" {
		// Replace encoder with NVidia hardware accelerated encoder.
		for i := 0; i < len(args); i++ {
			if args[i] == "libx264" {
				args[i] = "h264_nvenc"
			} else if args[i] == "libx265" {
				args[i] = "hevc_nvenc"
			}
		}
	} else if opt.Video.HardwareAcceleration != "off" {
		arg := []string{"-hwaccel", opt.Video.HardwareAcceleration}
		args = append(args, arg...)
	}

	// CRF.
	if opt.Video.Crf != 0 && opt.Video.Pass == "crf" {
		crf := strconv.Itoa(opt.Video.Crf)
		arg := []string{"-crf", crf}
		args = append(args, arg...)
	}

	// Bitrate.
	if opt.Video.Bitrate != "" && opt.Video.Bitrate != "0" {
		arg := []string{"-b:v", opt.Video.Bitrate}
		args = append(args, arg...)
	}

	// Minrate.
	if opt.Video.MinRate != "" && opt.Video.MinRate != "0" {
		arg := []string{"-minrate", opt.Video.MinRate}
		args = append(args, arg...)
	}

	// Maxrate.
	if opt.Video.MaxRate != "" && opt.Video.MaxRate != "0" {
		arg := []string{"-maxrate", opt.Video.MaxRate}
		args = append(args, arg...)
	}

	// Buffer Size.
	if opt.Video.BufSize != "" && opt.Video.BufSize != "0" {
		arg := []string{"-bufsize", opt.Video.BufSize}
		args = append(args, arg...)
	}

	// Pixel Format.
	if opt.Video.PixelFormat != "" && opt.Video.PixelFormat != "auto" {
		arg := []string{"-pix_fmt", opt.Video.PixelFormat}
		args = append(args, arg...)
	}

	// Frame Rate.
	if opt.Video.FrameRate != "" && opt.Video.PixelFormat != "auto" {
		arg := []string{"-r", opt.Video.FrameRate}
		args = append(args, arg...)
	}

	// Tune.
	if opt.Video.Tune != "" && opt.Video.Tune != "none" {
		arg := []string{"-tune", opt.Video.Tune}
		args = append(args, arg...)
	}

	// Profile.
	if opt.Video.Profile != "" && opt.Video.Profile != "none" {
		arg := []string{"-profile:v", opt.Video.Profile}
		args = append(args, arg...)
	}

	// Level.
	if opt.Video.Level != "" && opt.Video.Level != "none" {
		arg := []string{"-level", opt.Video.Level}
		args = append(args, arg...)
	}

	// Video Filters.
	vf := []string{"-vf", "\""}

	// Speed.
	if opt.Video.Speed != "" && opt.Video.Speed != "auto" {
		arg := "setpts=" + opt.Video.Speed
		vf = append(vf, arg)
	}

	vf = append(vf, "\"") // End of video filters.

	// Only push -vf flag if there are video filter arguments.
	if len(vf) > 3 {
		args = append(args, vf...)
	}

	extra := []string{
		"-y",
	}
	args = append(args, extra...)
	return args
}
