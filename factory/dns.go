package factory

import (
	"github.com/zalgonoise/dns/dns"
	"github.com/zalgonoise/dns/dns/core"
)

func DNSRepository(rtype string, fallbackDNS ...string) dns.Repository {
	var dnsRepo dns.Repository

	switch rtype {
	case "miekgdns":
		dnsRepo = core.New(fallbackDNS...)
	default:
		dnsRepo = core.New(fallbackDNS...)
	}

	return dnsRepo
}
