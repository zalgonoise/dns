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

func (s *udps) newCtxAndSpan(w dns.ResponseWriter, service string, attrs ...attr.Attr) (context.Context, spanner.Span) {
	var nsAttr = []attr.Attr{
		attr.String("module", service),
		attr.String("req_id", uuid.New().String()),
		attr.String("remote_addr", w.RemoteAddr().String()),
	}

	nsAttr = append(nsAttr, attrs...)
	namespace := attr.New("req", nsAttr)

	ctx, span := spanner.Start(context.Background(), service, namespace)

	return ctx, span
}
