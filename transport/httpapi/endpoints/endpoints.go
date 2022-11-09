package endpoints

import (
	"github.com/zalgonoise/dns/service"
	"github.com/zalgonoise/dns/store/encoder"
	"github.com/zalgonoise/dns/transport/httpapi"
	"github.com/zalgonoise/dns/transport/udp"
)

type endpoints struct {
	s   service.StoreWithHealth
	UDP udp.Server
	enc encoder.EncodeDecoder
}

func NewAPI(s service.Service, udps udp.Server) httpapi.HTTPAPI {
	return &endpoints{
		s:   s,
		UDP: udps,
		enc: encoder.New("json"),
	}
}
