package service

import (
	"context"

	dnsr "github.com/miekg/dns"
	"github.com/zalgonoise/dns/store"
)

// AnswerDNS uses the dns.Repository to reply to the dns.Msg `m` with the answer
// in store.Record `r`
func (l withLogger) AnswerDNS(ctx context.Context, r *store.Record, m *dnsr.Msg) {
	l.s.AnswerDNS(ctx, r, m)
}
