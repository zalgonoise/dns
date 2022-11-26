package service

import (
	"context"

	dnsr "github.com/miekg/dns"
	"github.com/zalgonoise/dns/store"
)

// AnswerDNS uses the dns.Repository to reply to the dns.Msg `m` with the answer
// in store.Record `r`
func (s *service) AnswerDNS(ctx context.Context, r *store.Record, m *dnsr.Msg) {
	switch r.Type {
	case "", "ANY":
		answers, err := s.store.FilterByDomain(ctx, r.Name)
		if err != nil || len(answers) == 0 {
			r.Type = "ANY"
			s.dns.Fallback(ctx, r, m)
			return
		}

		for _, ans := range answers {
			s.dns.Answer(ctx, ans, m)
		}
	default:
		answer, err := s.store.FilterByTypeAndDomain(ctx, r.Type, r.Name)
		if err != nil || answer.Addr == "" {
			s.dns.Fallback(ctx, r, m)
			return
		}

		s.dns.Answer(ctx, answer, m)
	}
}
