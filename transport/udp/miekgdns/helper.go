package miekgdns

import (
	"context"

	"github.com/google/uuid"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/attr"
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
