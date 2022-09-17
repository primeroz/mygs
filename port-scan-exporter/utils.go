package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debugf("%s took %s", name, elapsed)
}
