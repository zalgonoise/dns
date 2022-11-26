package factory

import (
	"io"
	"os"

	"github.com/zalgonoise/dns/cmd/config"
	"github.com/zalgonoise/dns/dns"
	"github.com/zalgonoise/dns/health"
	"github.com/zalgonoise/dns/service"
	svclog "github.com/zalgonoise/dns/service/middleware/logger"
	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/handlers/jsonh"
	"github.com/zalgonoise/logx/handlers/texth"
)

func writerFromPath(path string) io.Writer {
	confFile, err := os.Open(path)
	if err != nil {
		newFile, err := os.Create(path)
		if err != nil {
			return os.Stderr
		}
		if err := newFile.Sync(); err != nil {
			return os.Stderr
		}
		return newFile
	}
	return confFile
}

func Service(
	dnsRepo dns.Repository,
	storeRepo store.Repository,
	healthRepo health.Repository,
	conf *config.Config,
) service.Service {
	var (
		svc    service.Service = service.New(dnsRepo, storeRepo, healthRepo, conf)
		writer io.Writer       = os.Stderr

		logger logx.Logger
	)

	// short-circuit out
	if conf.Logger.Type == "" {
		return svc
	}

	if conf.Logger.Path != "" {
		writer = writerFromPath(conf.Logger.Path)
	}

	switch conf.Logger.Type {
	case "text":
		logger = logx.New(texth.New(writer))
	case "json":
		logger = logx.New(jsonh.New(writer))
	default:
		logger = nil
	}

	if logger == nil {
		return svc
	}

	return svclog.LogService(svc, logger)
}
