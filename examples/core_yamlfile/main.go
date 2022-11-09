package main

import (
	"os"

	"github.com/zalgonoise/dns/cmd/config"
	"github.com/zalgonoise/dns/factory"
	"github.com/zalgonoise/zlog/log"
)

func main() {
	cfg := config.New(
		config.StoreType("yaml"),
		config.StorePath("/tmp/dns/dns.list"),
	)

	// create HTTP / UDP server from config
	svr := factory.From(cfg)

	defer func() {
		err := svr.Stop()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}()

	// start HTTP server
	// you need to start DNS server with a GET request to localhost:8080/dns/start
	err := svr.Start()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
