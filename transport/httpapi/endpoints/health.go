package endpoints

import (
	"net/http"

	"github.com/zalgonoise/dns/transport/httpapi"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/attr"
)

func (e *endpoints) Health(w http.ResponseWriter, r *http.Request) {
	ctx := e.newCtx("http:health:get",
		attr.String("remote_addr", r.RemoteAddr),
		attr.String("user_agent", r.UserAgent()),
	)
	out := e.s.Health(ctx)

	w.WriteHeader(200)
	response, _ := e.enc.Encode(httpapi.HealthResponse{
		Message: "status and health report",
		Report:  out,
	})
	_, _ = w.Write(response)
	logx.From(ctx).Debug("successful status report response")
}
