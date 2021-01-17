package main

import (
	cmd "github.com/cgoder/gss/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("gss exec begin ...")

	// cmd.Transcode()
	// cmd.Split()
	cmd.Remux()

	log.Info("gss exec end ...")
}
