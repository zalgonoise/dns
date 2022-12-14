package miekgdns

import (
	"github.com/miekg/dns"
	"github.com/zalgonoise/dns/service"
	"github.com/zalgonoise/dns/transport/udp"
)

// udps implements the udp.Server interface
type udps struct {
	on   bool
	ans  service.Answering
	conf *udp.DNS
	srv  *dns.Server
	err  error
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
