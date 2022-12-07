package miekgdns

import (
	"context"

	"github.com/google/uuid"
	"github.com/miekg/dns"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/attr"
	"github.com/zalgonoise/x/spanner"
)

func (s *udps) newCtx(service string, attrs ...attr.Attr) context.Context {
	var nsAttr = []attr.Attr{
		attr.String("req_id", uuid.New().String()),
		attr.String("module", service),
	}
	nsAttr = append(nsAttr, attrs...)
	namespace := attr.New("req", nsAttr)

	return logx.InContext(context.Background(), s.logger.With(namespace))
}

func (s *udps) newCtxAndSpan(w dns.ResponseWriter, service string, attrs ...attr.Attr) (context.Context, spanner.Span, func()) {
	var nsAttr = []attr.Attr{
		attr.String("req_id", uuid.New().String()),
		attr.String("module", service),
	}
	nsAttr = append(nsAttr, attrs...)
	namespace := attr.New("req", nsAttr)

	ctx := logx.InContext(context.Background(), s.logger.With(namespace))
	ctx, span := spanner.Start(ctx, service,
		attr.String("remote_addr", w.RemoteAddr().String()),
	)

	return ctx, span, func() {
		span.End()

		spans := spanner.Extract(ctx)
		logx.From(ctx).Trace("trace", attr.Int("span_len", len(spans)), attr.New("spans", spans))
	}
}
