package main

import (
	"flag"
	"log"

	"github.com/nahid12105080/cacheDB/config"
	"github.com/nahid12105080/cacheDB/server"
)

func setUpFlags() {
	flag.StringVar(&config.Host, "host", config.Host, "host of the cacheDB server")
	flag.IntVar(&config.Port, "port", config.Port, "port of the cacheDB server")
	flag.Parse()
}

func main() {
	setUpFlags()

	log.Printf("Server starting on %s:%d\n", config.Host, config.Port)

	server.RunSyncTCPServer()
}
