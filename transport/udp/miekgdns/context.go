package miekgdns

import (
	"context"

	"github.com/google/uuid"
	"github.com/miekg/dns"
	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/spanner"
)

func (s *udps) newCtxAndSpan(w dns.ResponseWriter, service string, attrs ...attr.Attr) (context.Context, spanner.Span) {
	ctx, span := spanner.Start(
		context.Background(),
		service,
	)
	span.Add(attr.New("req", []attr.Attr{
		attr.String("module", service),
		attr.String("req_id", uuid.New().String()),
		attr.String("remote_addr", w.RemoteAddr().String()),
	}))
	span.Add(attrs...)
	return ctx, span
}
