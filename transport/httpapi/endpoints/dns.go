package endpoints

import (
	"net/http"

	"github.com/zalgonoise/x/ptr"
)

func (e *endpoints) StartDNS(w http.ResponseWriter, r *http.Request) {
	ctx, _, done := e.newCtxAndSpan(r, "http.StartDNS")
	defer done()

	var err error
	go func() {
		err = e.UDP.Start(ctx)
	}()

	if err != nil {
		res := NewResponse[string](500, "failed to start UDP server", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	res := NewResponse(200, "started UDP server successfully", nil, ptr.To("ok"))
	res.WriteHTTP(ctx, w)
}

func (e *endpoints) StopDNS(w http.ResponseWriter, r *http.Request) {
	ctx, _, done := e.newCtxAndSpan(r, "http.StopDNS")
	defer done()

	err := e.UDP.Stop(ctx)
	if err != nil {
		res := NewResponse[string](500, "failed to stop UDP server", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	res := NewResponse(200, "stopped UDP server successfully", nil, ptr.To("ok"))
	res.WriteHTTP(ctx, w)
}

func (e *endpoints) ReloadDNS(w http.ResponseWriter, r *http.Request) {
	ctx, _, done := e.newCtxAndSpan(r, "http.ReloadDNS")
	defer done()

	var err error
	if e.UDP.Running(ctx) {
		err = e.UDP.Stop(ctx)
		if err != nil {
			res := NewResponse[string](500, "failed to stop UDP server", err, nil)
			res.WriteHTTP(ctx, w)
			return
		}
	}

	go func() {
		err = e.UDP.Start(ctx)
	}()

	if err != nil {
		res := NewResponse[string](500, "failed to start UDP server", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	res := NewResponse(200, "reloaded UDP server successfully", nil, ptr.To("ok"))
	res.WriteHTTP(ctx, w)
}
