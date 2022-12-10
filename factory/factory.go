package factory

import (
	"context"
	"os"
	"strings"

	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/dns/cmd/config"
	"github.com/zalgonoise/dns/dns"
	"github.com/zalgonoise/dns/health"
	"github.com/zalgonoise/dns/service"
	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/dns/transport/httpapi"
	"github.com/zalgonoise/x/spanner"
	"github.com/zalgonoise/x/spanner/export"
)

func From(conf *config.Config) httpapi.Server {
	_, s := spanner.Start(context.Background(), "factory.From")
	defer s.End()

	// initialize logger and register it in the spanner
	logger := Logger(conf.Logger.Type, conf.Logger.Path)
	spanner.To(export.Logger(logger))
	s.Event("initialized logger")

	// initialize DNS repository
	dnsRepo := dns.WithTrace(DNSRepository(
		conf.DNS.Type,
		strings.Split(conf.DNS.FallbackDNS, ",")...,
	))
	s.Event("initialized DNS repository")

	// initialize store repository
	storeRepo := store.WithTrace(StoreRepository(
		conf.Store.Type,
		conf.Store.Path,
	))
	s.Event("initialized Store repository")

	// initialize health repository
	healthRepo := health.WithTrace(HealthRepository(
		conf.Health.Type,
	))
	s.Event("initialized Health repository")

	// intialize service
	svc := service.WithTrace(service.New(dnsRepo, storeRepo, healthRepo, conf))
	s.Event("initialized service")

	// initialize HTTP and DNS servers
	https, udps := Server(
		conf.DNS.Type,
		conf.DNS.Address,
		conf.DNS.Prefix,
		conf.DNS.Proto,
		conf.HTTP.Port,
		svc,
	)
	s.Event("initialized HTTP and UDP servers")

	// autostart DNS if configured
	if conf.Autostart.DNS {
		go func() {
			ctx, s := spanner.Start(context.Background(), "factory.From",
				attr.String("action", "autostart DNS server"),
			)

			err := udps.Start(ctx)
			if err != nil {
				s.Event("error starting DNS server", attr.String("error", err.Error()))
				s.End()
				os.Exit(1)
			}
			s.End()
		}()
	}

	s.Event("initialized all modules successfully")
	return https
}
