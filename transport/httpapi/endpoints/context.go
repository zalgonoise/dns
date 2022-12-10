package endpoints

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/attr"
	"github.com/zalgonoise/x/spanner"
)

// newCtx creates a new context for service `service` with attributes `attrs`, scoped to
// a "req" namespace that includes a UUID for the request and the service string `service`
//
// The context is also wrapped with the endpoints encoder (or a new JSON encoder) so
// the response writer can use it
func (e *endpoints) newCtx(service string, attrs ...attr.Attr) context.Context {
	var nsAttr = []attr.Attr{
		attr.String("req_id", uuid.New().String()),
		attr.String("module", service),
	}
	nsAttr = append(nsAttr, attrs...)
	namespace := attr.New("req", nsAttr)

	ctx := context.WithValue(context.Background(), ResponseEncoderKey, e.enc)
	return logx.InContext(ctx, e.logger.With(namespace))
}

// newCtxAndSpan creates a new context for service `service` with attributes `attrs`, scoped to
// a "req" namespace that includes a UUID for the request and the service string `service`,
// and creates a Span for the given request
//
// The context is also wrapped with the endpoints encoder (or a new JSON encoder) so
// the response writer can use it
//
// The input *http.Request `r` is used to registed the remote address and user agent in the root span
//
// The resulting context is returned alongside the created Span and a closure function for it
func (e *endpoints) newCtxAndSpan(r *http.Request, service string, attrs ...attr.Attr) (context.Context, spanner.Span) {
	var nsAttr = []attr.Attr{
		attr.String("module", service),
		attr.String("req_id", uuid.New().String()),
		attr.String("remote_addr", r.RemoteAddr),
		attr.String("user_agent", r.UserAgent()),
	}
	nsAttr = append(nsAttr, attrs...)
	namespace := attr.New("req", nsAttr)

	ctx := context.WithValue(r.Context(), ResponseEncoderKey, e.enc)
	ctx, span := spanner.Start(ctx, service, namespace)

	return ctx, span
}
