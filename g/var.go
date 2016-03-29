package g

import (
	"log"
	"os"
)

func Hostname(server *DBServer) string {
	hostname := server.Endpoint
	if "" != hostname {
		return hostname
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
	}
	return hostname
}
