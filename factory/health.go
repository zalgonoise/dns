package factory

import (
	"github.com/zalgonoise/dns/health"
	"github.com/zalgonoise/dns/health/simplehealth"
)

func HealthRepository(rtype string) health.Repository {
	switch rtype {
	case "simple", "simplehealth":
		return simplehealth.New()
	default:
		return simplehealth.New()
	}
}
