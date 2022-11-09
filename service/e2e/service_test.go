package e2e

import (
	"github.com/zalgonoise/dns/cmd/config"
	"github.com/zalgonoise/dns/dns/core"
	"github.com/zalgonoise/dns/health/simplehealth"
	"github.com/zalgonoise/dns/service"
	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/dns/store/memmap"
)

var (
	record1 *store.Record = store.New().Type("A").Name("not.a.dom.ain").Addr("192.168.0.10").Build()
	record2 *store.Record = store.New().Type("A").Name("also.not.a.dom.ain").Addr("192.168.0.15").Build()
	record3 *store.Record = store.New().Type("A").Name("really.not.a.dom.ain").Addr("192.168.0.10").Build()
)

func initializeService() service.Service {
	return service.New(
		core.New(),
		memmap.New(),
		simplehealth.New(),
		config.Default(),
	)
}
