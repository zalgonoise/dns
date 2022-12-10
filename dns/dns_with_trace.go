package dns

import (
	"context"

	"github.com/miekg/dns"
	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/x/spanner"
)

type withTrace struct {
	r Repository
}

func WithTrace(r Repository) Repository {
	return withTrace{
		r: r,
	}
}

// Answer will write the IP address present in the store.Record in the dns.Msg
// slice of Answers
func (t withTrace) Answer(ctx context.Context, r *store.Record, m *dns.Msg) {
	ctx, s := spanner.Start(ctx, "dns.Answer",
		attr.New("record", r),
	)
	defer s.End()

	t.r.Answer(ctx, r, m)
	s.Add(attr.New("answer", m.Answer))
}

// Fallback is called when the DNS store does not hold a record for the requested
// domain, so the DNS service spawns a DNS client that will query the fallback server
// and write that answer to the dns.Msg
func (t withTrace) Fallback(ctx context.Context, r *store.Record, m *dns.Msg) {
	ctx, s := spanner.Start(ctx, "dns.Fallback",
		attr.New("record", r),
	)
	defer s.End()

	t.r.Fallback(ctx, r, m)
	s.Add(attr.New("answer", m.Answer))
}
