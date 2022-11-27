package miekgdns

import (
	"github.com/miekg/dns"
	"github.com/zalgonoise/dns/transport/udp"
	"github.com/zalgonoise/logx/attr"
)

// Start launches the DNS server, returning an error
func (u *udps) Start() error {
	if u.on {
		return udp.ErrAlreadyRunning
	}
	dns.HandleFunc(u.conf.Prefix, u.handleRequest)
	u.srv = &dns.Server{
		Addr: u.conf.Addr,
		Net:  u.conf.Proto,
	}
	u.on = true

	if u.logger != nil {
		u.logger.Debug("starting UDP server", attr.New("config", []attr.Attr{
			attr.String("address", u.conf.Addr),
			attr.String("protocol", u.conf.Proto),
		}))
	}

	return u.srv.ListenAndServe()
}

// Stop gracefully stops the DNS server, returning an error
func (u *udps) Stop() error {
	if !u.on {
		return udp.ErrNotRunning
	}
	if u.logger != nil {
		u.logger.Debug("stopping UDP server", attr.New("config", []attr.Attr{
			attr.String("address", u.conf.Addr),
			attr.String("protocol", u.conf.Proto),
		}))
	}
	u.on = false
	return u.srv.Shutdown()
}

// Running returns a boolean on whether the UDP server is running or not
func (u *udps) Running() bool {
	return u.on
}
