package cron

import (
	"github.com/coraldane/mymon/g"
	"github.com/toolkits/logger"
	"time"
)

func Heartbeat() {
	time.Sleep(10 * time.Second)
	SleepRandomDuration()
	for {
		heartbeat()
		d := time.Duration(g.Config().Interval) * time.Second
		time.Sleep(d)
	}
}

func heartbeat() {
	for _, server := range g.Config().DBServerList {
		err := FetchData(server)
		if err != nil {
			logger.Errorln(err)
		}
	}
}
