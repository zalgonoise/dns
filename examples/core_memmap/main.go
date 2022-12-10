package main

import (
	"context"
	"os"

	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/dns/cmd/config"
	"github.com/zalgonoise/dns/factory"
	"github.com/zalgonoise/x/spanner"
)

func main() {
	ctx, s := spanner.Start(context.Background(), "running in-memory example")

	cfg := config.New()

	// create HTTP / UDP server from config
	svr := factory.From(cfg)

	defer func() {
		err := svr.Stop(ctx)
		if err != nil {
			s.Event("failed to stop HTTP server in the deferred function",
				attr.String("error", err.Error()),
			)
			s.End()
			os.Exit(1)
		}
	}()

	// start HTTP server
	// you need to start DNS server with a GET request to localhost:8080/dns/start
	err := svr.Start(ctx)
	if err != nil {
		s.Event("failed to start HTTP server", attr.String("error", err.Error()))
		s.End()
		os.Exit(1)
	}
	s.End()
}
