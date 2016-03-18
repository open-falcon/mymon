// MySQL Performance Monitor(For open-falcon)
// Write by coraldane<coraldane@163.com>
package main

import (
	"flag"
	"fmt"
	"github.com/coraldane/mymon/cron"
	"github.com/coraldane/mymon/g"
	"github.com/toolkits/logger"
	"log"
	"os"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if err := g.ParseConfig(*cfg); err != nil {
		log.Fatalln(err)
	}

	logger.SetLevelWithDefault(g.Config().LogLevel, "info")

	go cron.Heartbeat()

	select {}
}
