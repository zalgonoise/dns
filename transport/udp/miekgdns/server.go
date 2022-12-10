package miekgdns

import (
	"context"

	"github.com/miekg/dns"
	"github.com/zalgonoise/dns/transport/udp"
	"github.com/zalgonoise/logx/attr"
	"github.com/zalgonoise/x/spanner"
)

// Start launches the DNS server, returning an error
func (u *udps) Start(ctx context.Context) error {
	_, s := spanner.Start(ctx, "udp.Start")
	defer s.End()

	if u.on {
		s.Event("UDP server is already running", attr.String("error", udp.ErrAlreadyRunning.Error()))
		return udp.ErrAlreadyRunning
	}

	dns.HandleFunc(u.conf.Prefix, u.handleRequest)
	u.srv = &dns.Server{
		Addr: u.conf.Addr,
		Net:  u.conf.Proto,
	}
	u.on = true

	err := u.srv.ListenAndServe()
	if err != nil {
		s.Event("failed to start UDP server", attr.String("error", err.Error()))
		return err
	}
	return nil
}

// Stop gracefully stops the DNS server, returning an error
func (u *udps) Stop(ctx context.Context) error {
	_, s := spanner.Start(ctx, "udp.Stop")
	defer s.End()

	if !u.on {
		s.Event("UDP server is not running", attr.String("error", udp.ErrNotRunning.Error()))
		return udp.ErrNotRunning
	}

	u.on = false
	err := u.srv.Shutdown()
	if err != nil {
		s.Event("failed to stop UDP server", attr.String("error", err.Error()))
	}
	return nil
}

// Running returns a boolean on whether the UDP server is running or not
func (u *udps) Running(ctx context.Context) bool {
	_, s := spanner.Start(ctx, "udp.Running", attr.New("running", u.on))
	defer s.End()

	return u.on
}
