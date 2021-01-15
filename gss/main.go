package main

import (
	cmd "github.com/cgoder/gss/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("transcode begin ...")

	cmd.Transcode()

	log.Info("transcode end ...")
}
