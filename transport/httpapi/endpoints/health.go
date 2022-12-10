package endpoints

import (
	"net/http"
)

func (e *endpoints) Health(w http.ResponseWriter, r *http.Request) {
	ctx, s := e.newCtxAndSpan(r, "http.Health")
	defer s.End()

	report := e.s.Health(ctx)

	res := NewResponse(200, "added record successfully", nil, report)
	res.WriteHTTP(ctx, w)
}
