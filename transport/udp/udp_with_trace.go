package udp

import (
	"context"

	"github.com/zalgonoise/logx/attr"
	"github.com/zalgonoise/x/spanner"
)

type withTrace struct {
	s Server
}

func WithTrace(s Server) Server {
	return withTrace{
		s: s,
	}
}

// Start launches the DNS server, returning an error
func (t withTrace) Start(ctx context.Context) error {
	ctx, s := spanner.Start(ctx, "udp.Start")
	defer s.End()

	err := t.s.Start(ctx)
	if err != nil {
		s.Event("error starting UDP server", attr.New("error", err.Error()))
	}

	return err
}

// Stop gracefully stops the DNS server, returning an error
func (t withTrace) Stop(ctx context.Context) error {
	ctx, s := spanner.Start(ctx, "udp.Stop")
	defer s.End()

	err := t.s.Stop(ctx)
	if err != nil {
		s.Event("error stopping UDP server", attr.New("error", err.Error()))
	}

	return err
}

// Running returns a boolean on whether the UDP server is running or not
func (t withTrace) Running(ctx context.Context) bool {
	ctx, s := spanner.Start(ctx, "udp.Running")
	defer s.End()

	running := t.s.Running(ctx)
	s.Add(attr.New("running", running))

	return running
}
