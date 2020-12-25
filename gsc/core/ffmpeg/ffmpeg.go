package ffmpeg

import (
	"bufio"
	"context"
	"errors"
	"io"
	"math"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// openencoder
const (
	ffmpegCmd = "ffmpeg"
)

// regex extracting the duration value of an input stream.
// the duration has the format HOURS:MINUTES:SECONDS.MILLISECONDS
var ffmpegStreamDurationRegex = regexp.MustCompile(`Duration: ([0-9]*):([0-9]*):([0-9]*).([0-9]*)`)

// regex extracting the out_time_ms value from an ffmpeg progress line.
var ffmpegOutTimeMSRegex = regexp.MustCompile(`out_time_ms=([0-9]*)`)

// FFmpeg struct.
type FFmpeg struct {
	Status      ProcStatus
	cmd         *exec.Cmd
	isCancelled bool
}

type ProcStatus struct {
	Frame      int
	FPS        float64
	Bitrate    float64
	TotalSize  int
	OutTimeMS  int
	OutTime    string
	DupFrames  int
	DropFrames int
	Speed      string
	Progress   string
	Percent    float64
}

// ffmpegOptions struct passed into Ffmpeg.Run.
type ffmpegOptions struct {
	Input  string
	Output string
	// 优先parse cmdOpt结构中的参数
	CmdString []string `json:"CmdString"` // Raw flag options.

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
func (f *FFmpeg) Execute(ctx context.Context, ffopt FFOption) error {

	// log.Info("\t ffopt.CmdString--- ", ffopt.CmdString)
	// log.Info("\t ffopt.CmdSlice--- ", ffopt.CmdSlice)
	// log.Info("\t ffopt.arguments--- ", ffopt.arguments)
	// Parse options and add to args slice.
	ok, cmds := getFFmpegCmdString(ffopt)
	if !ok {
		err := errors.New("ffmpeg get cmds error")
		return err
	}

	// Execute command.
	log.Printf("running FFmpeg with options: %s", cmds)
	// for _, v := range cmds {
	// 	log.Info("vvv--- ", v)
	// }
	// f.cmd = exec.Command(ffmpegCmd, cmds...)
	f.cmd = exec.CommandContext(ctx, ffmpegCmd, cmds...)

	// ffmpeg writes its output to stderr
	stderr, err := f.cmd.StderrPipe()
	if err != nil {
		return err
	}

	// the progress data is written to stdout
	stdout, err := f.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrChan := make(chan string)
	go readLinesToChannel(stderr, stderrChan)
	stdoutChan := make(chan string)
	go readLinesToChannel(stdout, stdoutChan)

	// Update status goroutine
	go f.execProgress(ctx, stdoutChan, stderrChan)

	// Run
	f.cmd.Start()

	errs := f.cmd.Wait()
	if errs != nil {
		if f.isCancelled {
			return errors.New("cancelled")
		}
		return errs
	}
	return nil
}

// Utilities for parsing ffmpeg options.
func getFFmpegCmdString(opt FFOption) (bool, []string) {
	if opt.CmdSlice == nil {
		return false, nil
	}

	var stdoutName string
	if runtime.GOOS == "windows" {
		// pipe:1 is the windows equivalent of /dev/stdout
		stdoutName = "pipe:1"
	} else {
		stdoutName = os.Stdout.Name()
	}

	args := []string{
		"-hide_banner",
		// "-loglevel", "error", // Set loglevel to fail job on errors.
		"-progress", stdoutName,
	}

	// Add src
	if opt.Input != "" {
		args = append(args, "-i", opt.Input)
	}

	// Add option
	for _, v := range opt.CmdSlice {
		args = append(args, v)
	}

	// Add dst
	if opt.Output != "" {
		args = append(args, opt.Output)
	}

	return true, args
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
		if progressSplit[0] == "" {
			return
		}
		// log.Info(progressSplit)
		k := progressSplit[0]
		v := progressSplit[1]

		switch k {
		case "frame":
			frame, _ := strconv.Atoi(v)
			f.Status.Frame = frame
		case "fps":
			fps, _ := strconv.ParseFloat(v, 64)
			f.Status.FPS = fps
		case "bitrate":
			v = strings.Replace(v, "kbits/s", "", -1)
			bitrate, _ := strconv.ParseFloat(v, 64)
			f.Status.Bitrate = bitrate
		case "total_size":
			size, _ := strconv.Atoi(v)
			f.Status.TotalSize = size
		case "out_time_ms":
			outTimeMS, _ := strconv.Atoi(v)
			f.Status.OutTimeMS = outTimeMS
		case "out_time":
			f.Status.OutTime = v
		case "dup_frames":
			frames, _ := strconv.Atoi(v)
			f.Status.DupFrames = frames
		case "drop_frames":
			frames, _ := strconv.Atoi(v)
			f.Status.DropFrames = frames
		case "speed":
			f.Status.Speed = v
		case "progress":
			// if v == "end" {
			// 	// end
			// 	f.Status.Progress = "1"
			// } else {
			// 	// continue
			// 	f.Status.Progress = "0"
			// }
			f.Status.Progress = v
		}
	}
}

func (f *FFmpeg) execProgress(ctx context.Context, stdoutChan <-chan string, stderrChan <-chan string) {

	var duration time.Duration

	for {
		select {
		case <-ctx.Done():
			log.Info("ffmpeg process done.")
			return
		case line, ok := <-stderrChan:
			if !ok {
				stderrChan = nil
			}
			if duration == 0 {
				// look for video duration line
				if m := ffmpegStreamDurationRegex.FindStringSubmatch(line); m != nil {
					hours, err := strconv.ParseInt(m[1], 10, 64)
					if err != nil {
						log.Errorf("error parsing duration value: %s", err.Error())
						// return
					}
					duration += time.Duration(hours) * time.Hour

					minutes, err := strconv.ParseInt(m[2], 10, 64)
					if err != nil {
						log.Errorf("error parsing duration value: %s", err.Error())
						// return
					}
					duration += time.Duration(minutes) * time.Minute

					seconds, err := strconv.ParseInt(m[3], 10, 64)
					if err != nil {
						log.Errorf("error parsing duration value: %s", err.Error())
						// return
					}
					duration += time.Duration(seconds) * time.Second

					millis, err := strconv.ParseInt(m[4], 10, 64)
					if err != nil {
						log.Errorf("error parsing duration value: %s", err.Error())
						// return
					}
					duration += time.Duration(millis) * time.Millisecond
					// log.Info("execProgress input file Duration: ", duration)
				}
			}
		case line, ok := <-stdoutChan:
			if !ok {
				stdoutChan = nil
			}

			str := strings.Replace(line, " ", "", -1)
			parts := strings.Split(str, " ")
			f.setProgressParts(parts)

			if m := ffmpegOutTimeMSRegex.FindStringSubmatch(line); m != nil {
				if duration == 0 {
					// we haven't found the input duration value,
					// which should always occur before the -progress output
					log.Errorf("could not find duration of input file")
					// return
				}

				millis, err := strconv.ParseInt(m[1], 10, 64)
				if err != nil {
					log.Errorf("error parsing output time value: %s", err.Error())
					// return
				}

				p := time.Duration(millis) * time.Microsecond
				// limit the progress to 1.0, as the duration printed
				// by ffmpeg may be a bit inaccurate.
				// we could use ffprobe to get the precise duration of the
				// period, but it's really not worth the hassle.
				progress := math.Min(1, float64(p)/float64(duration))
				f.Status.Percent = progress * 100

				// log.Info("execProgress:progress: ", f.Status.Percent, "[", f.Status.Progress, "]")
			}

			// log.Info(JsonFormat(f.Status))
		}
	}
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

func readLinesToChannel(reader io.Reader, lineChan chan<- string) {
	r := bufio.NewScanner(reader)
	for r.Scan() {
		lineChan <- r.Text()
	}
	close(lineChan)
}
