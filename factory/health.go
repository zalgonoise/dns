package factory

import (
	"github.com/zalgonoise/dns/health"
	"github.com/zalgonoise/dns/health/simplehealth"
)

func HealthRepository(rtype string) health.Repository {
	var healthRepo health.Repository

	switch rtype {
	case "simple", "simplehealth":
		healthRepo = simplehealth.New()
	default:
		healthRepo = simplehealth.New()
	}

	return healthRepo
}
