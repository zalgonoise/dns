package factory

import (
	"github.com/zalgonoise/dns/service"
	"github.com/zalgonoise/dns/transport/httpapi"
	"github.com/zalgonoise/dns/transport/httpapi/endpoints"
	"github.com/zalgonoise/dns/transport/udp"
	"github.com/zalgonoise/dns/transport/udp/miekgdns"
	"github.com/zalgonoise/logx"
)

func UDPServer(stype, address, prefix, proto string, svc service.Service, logger logx.Logger) udp.Server {
	var udps udp.Server

	switch stype {
	case "miekgdns":
		udps = miekgdns.WithTrace(
			miekgdns.NewServer(
				udp.NewDNS().
					Addr(address).
					Prefix(prefix).
					Proto(proto).
					Build(),
				svc,
			),
			logger,
		)
	default:
		udps = miekgdns.WithTrace(
			miekgdns.NewServer(
				udp.NewDNS().
					Addr(address).
					Prefix(prefix).
					Proto(proto).
					Build(),
				svc,
			),
			logger,
		)
	}

	return udps
}

func Server(
	dnstype, dnsAddress, dnsPrefix, dnsProto string,
	httpPort int,
	svc service.Service,
	logger logx.Logger,
) (httpapi.Server, udp.Server) {
	udps := UDPServer(dnstype, dnsAddress, dnsPrefix, dnsProto, svc, logger)
	apis := endpoints.WithTrace(
		endpoints.NewAPI(svc, udps),
		logger,
	)
	https := httpapi.NewServer(apis, httpPort)

	return https, udps
}
