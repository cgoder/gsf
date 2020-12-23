package ffmpeg

import (
	"encoding/json"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

const ffprobeCmd = "ffprobe"

// FFProbe struct.
type FFProbe struct{}

// Execute runs an FFProbe command.
func (f FFProbe) Execute(input string, cmdOpt string) *Metadata {
	args := []string{
		"-i", input,
		"-show_format",
		"-show_streams",
		"-print_format", "json",
		"-v", "quiet",
	}

	if cmdOpt != "" {
		args = append(args, cmdOpt)
	}

	// Execute command.
	cmd := exec.Command(ffprobeCmd, args...)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(err.Error())
	}
	// log.Info(string(stdout))

	dat := &Metadata{}
	if err := json.Unmarshal([]byte(stdout), &dat); err != nil {
		panic(err)
	}
	return dat
}
