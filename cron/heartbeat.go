package cron

import (
	"github.com/coraldane/mymon/g"
	"github.com/toolkits/logger"
	"time"
)

func Heartbeat(server *g.DBServer) {
	time.Sleep(10 * time.Second)
	SleepRandomDuration()
	for {
		heartbeat(server)
		d := time.Duration(g.Config().Interval) * time.Second
		time.Sleep(d)
	}
}

func heartbeat(server *g.DBServer) {
	err := FetchData(server)
	if err != nil {
		logger.Errorln(err)
	}
}
