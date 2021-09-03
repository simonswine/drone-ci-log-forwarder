package main

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/simonswine/drone-ci-log-forwarder/drone"
)

var (
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	token  = os.Getenv("DRONE_TOKEN")
	host   = os.Getenv("DRONE_SERVER")
)

func main() {

	drone := drone.New(host, token).WithLogger(logger)

	if err := drone.Run(); err != nil {
		// log.Fatal(drone)
		level.Error(logger).Log("error", err)
	}
}
