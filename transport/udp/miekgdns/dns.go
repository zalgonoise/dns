package miekgdns

import (
	"github.com/miekg/dns"
	"github.com/zalgonoise/dns/service"
	"github.com/zalgonoise/dns/transport/udp"
	"github.com/zalgonoise/logx"
)

// udps implements the udp.Server interface
type udps struct {
	on     bool
	ans    service.Answering
	conf   *udp.DNS
	srv    *dns.Server
	err    error
	logger logx.Logger
}

// NewServer returns a github.com/miekg/dns implementation of udp.Server
//
// It takes in a DNS that serves as configuration, and a service.Answering
// interface to perform the necessary calls alongside the record store, returning
// a udps as a udp.Server
func NewServer(conf *udp.DNS, s service.Answering) udp.Server {
	if conf == nil {
		conf = udp.NewDNS().Build()
	}
	return &udps{
		conf: conf,
		ans:  s,
	}
}

// WithTrace returns the same UDP server configured with a logger
func WithTrace(server udp.Server, logger logx.Logger) udp.Server {
	if server == nil {
		return nil
	}
	if _, ok := server.(*udps); !ok {
		return server
	}
	if logger == nil {
		logger = logx.Default()
	}
	return &udps{
		conf:   server.(*udps).conf,
		ans:    server.(*udps).ans,
		logger: logger,
	}
}
