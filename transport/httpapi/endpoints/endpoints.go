package endpoints

import (
	"github.com/zalgonoise/dns/service"
	"github.com/zalgonoise/dns/store/encoder"
	"github.com/zalgonoise/dns/transport/httpapi"
	"github.com/zalgonoise/dns/transport/udp"
	"github.com/zalgonoise/logx"
)

type endpoints struct {
	s      service.StoreWithHealth
	UDP    udp.Server
	enc    encoder.EncodeDecoder
	logger logx.Logger
}

func NewAPI(s service.Service, udps udp.Server) httpapi.HTTPAPI {
	return &endpoints{
		s:   s,
		UDP: udps,
		enc: encoder.New("json"),
	}
}

func NewAPIWithTrace(api httpapi.HTTPAPI, logger logx.Logger) httpapi.HTTPAPI {
	if _, ok := api.(*endpoints); !ok {
		return api
	}
	if logger == nil {
		logger = logx.Default()
	}
	return &endpoints{
		s:      api.(*endpoints).s,
		UDP:    api.(*endpoints).UDP,
		enc:    api.(*endpoints).enc,
		logger: logger,
	}
}
