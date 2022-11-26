package miekgdns

import (
	"context"

	"github.com/miekg/dns"
	"github.com/zalgonoise/dns/store"
)

func (u *udps) answer(ctx context.Context, r *store.Record, m *dns.Msg) {
	name := r.Name
	if r.Name[len(r.Name)-1] == u.conf.Prefix[0] {
		name = r.Name[:len(r.Name)-1]
	}

	u.ans.AnswerDNS(
		ctx,
		store.New().Name(name).Type(r.Type).Build(),
		m,
	)
}

func (u *udps) handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	ctx := context.Background()
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		u.parseQuery(ctx, m)
		if u.err != nil {
			return
		}
	}

	err := w.WriteMsg(m)
	if err != nil {
		u.err = err
	}
}

func (u *udps) parseQuery(ctx context.Context, m *dns.Msg) {
	for _, question := range m.Question {
		r := store.New().Name(question.Name)
		switch question.Qtype {
		case dns.TypeA:
			u.answer(
				ctx,
				r.Type(store.TypeA.String()).Build(),
				m,
			)
		case dns.TypeAAAA:
			u.answer(
				ctx,
				r.Type(store.TypeAAAA.String()).Build(),
				m,
			)
		case dns.TypeCNAME:
			u.answer(
				ctx,
				r.Type(store.TypeCNAME.String()).Build(),
				m,
			)
		case dns.TypeANY:
			u.answer(
				ctx,
				r.Build(),
				m,
			)
		}
	}
}
