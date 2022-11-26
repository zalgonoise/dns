package endpoints

import (
	"context"
	"net/http"

	"github.com/zalgonoise/dns/transport/httpapi"
)

func (e *endpoints) Health(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	out := e.s.Health(ctx)

	w.WriteHeader(200)
	response, _ := e.enc.Encode(httpapi.HealthResponse{
		Message: "status and health report",
		Report:  out,
	})
	_, _ = w.Write(response)
}
