package cmd

import (
	"context"
	"os"

	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/dns/cmd/config"
	"github.com/zalgonoise/dns/cmd/flags"
	"github.com/zalgonoise/dns/factory"
	"github.com/zalgonoise/spanner"

	"github.com/zalgonoise/dns/transport/httpapi"
)

// Run starts the DNS app based on the input configuration (file configuration,
// CLI flags, and OS environment variables)
//
// Blocking call; will error out (with os.Exit(1)) if failed
func Run() {
	ctx, s := spanner.Start(context.Background(), "starting DNS server")

	var (
		conf *config.Config
		svr  httpapi.Server
	)

	// get config
	conf = flags.ParseFlags()

	// create HTTP / UDP servers
	svr = factory.From(conf)

	// start HTTP server
	// defer graceful closure
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

	err := svr.Start(ctx)
	if err != nil {
		s.Event("failed to start HTTP server", attr.String("error", err.Error()))
		s.End()
		os.Exit(1)
	}
	s.End()
}
