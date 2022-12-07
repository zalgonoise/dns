package factory

import (
	"context"
	"os"
	"strings"

	"github.com/zalgonoise/dns/cmd/config"
	"github.com/zalgonoise/dns/service"
	"github.com/zalgonoise/dns/transport/httpapi"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/attr"
)

func From(conf *config.Config) httpapi.Server {
	// initialize logger
	logger := Logger(conf.Logger.Type, conf.Logger.Path)

	// initialize DNS repository
	dnsRepo := DNSRepository(
		conf.DNS.Type,
		strings.Split(conf.DNS.FallbackDNS, ",")...,
	)

	// initialize store repository
	storeRepo := StoreRepository(
		conf.Store.Type,
		conf.Store.Path,
	)

	// initialize health repository
	healthRepo := HealthRepository(
		conf.Health.Type,
	)

	// intialize service
	svc := service.WithTrace(service.New(dnsRepo, storeRepo, healthRepo, conf))

	// initialize HTTP and DNS servers
	https, udps := Server(
		conf.DNS.Type,
		conf.DNS.Address,
		conf.DNS.Prefix,
		conf.DNS.Proto,
		conf.HTTP.Port,
		svc,
		logger,
	)

	// autostart DNS if configured
	if conf.Autostart.DNS {
		go func() {
			ctx := logx.InContext(context.Background(), logger)
			err := udps.Start(ctx)
			if err != nil {
				logx.Fatal("error starting DNS server",
					attr.String("error", err.Error()),
				)
				os.Exit(1)
			}
		}()
	}
	return https
}
