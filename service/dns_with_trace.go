package service

import (
	"context"

	dnsr "github.com/miekg/dns"
	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/x/spanner"
)

// AnswerDNS uses the dns.Repository to reply to the dns.Msg `m` with the answer
// in store.Record `r`
func (t withTrace) AnswerDNS(ctx context.Context, r *store.Record, m *dnsr.Msg) {
	ctx, s := spanner.Start(ctx, "service.AnswerDNS", attr.New("record", r))
	defer s.End()

	t.s.AnswerDNS(ctx, r, m)
}
