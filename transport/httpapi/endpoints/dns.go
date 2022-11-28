package endpoints

import (
	"net/http"

	"github.com/zalgonoise/dns/transport/httpapi"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/attr"
)

func (e *endpoints) StartDNS(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		ctx = e.newCtx("http:dns:start",
			attr.String("remote_addr", r.RemoteAddr),
			attr.String("user_agent", r.UserAgent()),
		)
	)
	logx.From(ctx).Debug("/dns/start request")

	go func() {
		err = e.UDP.Start(ctx)
	}()

	if err != nil {
		w.WriteHeader(500)
		response, _ := e.enc.Encode(httpapi.DNSResponse{
			Success: false,
			Message: "failed to start DNS server",
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Error("failed to start DNS server", attr.String("error", err.Error()))
		return
	}
	w.WriteHeader(200)
	response, _ := e.enc.Encode(httpapi.DNSResponse{
		Success: true,
		Message: "started DNS server",
	})
	_, _ = w.Write(response)
	logx.From(ctx).Debug("started DNS server")
}

func (e *endpoints) StopDNS(w http.ResponseWriter, r *http.Request) {
	var ctx = e.newCtx("http:dns:stop",
		attr.String("remote_addr", r.RemoteAddr),
		attr.String("user_agent", r.UserAgent()),
	)
	logx.From(ctx).Debug("/dns/stop request")

	err := e.UDP.Stop(ctx)
	if err != nil {
		w.WriteHeader(500)
		response, _ := e.enc.Encode(httpapi.DNSResponse{
			Success: false,
			Message: "failed to stop DNS server",
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Error("failed to stop DNS server", attr.String("error", err.Error()))
		return
	}
	w.WriteHeader(200)
	response, _ := e.enc.Encode(httpapi.DNSResponse{
		Success: true,
		Message: "stopped DNS server",
	})
	_, _ = w.Write(response)
	logx.From(ctx).Debug("stopped DNS server")
}

func (e *endpoints) ReloadDNS(w http.ResponseWriter, r *http.Request) {
	var ctx = e.newCtx("http:dns:reload",
		attr.String("remote_addr", r.RemoteAddr),
		attr.String("user_agent", r.UserAgent()),
	)
	logx.From(ctx).Debug("/dns/reload request")

	err := e.UDP.Stop(ctx)
	if err != nil {
		w.WriteHeader(500)
		response, _ := e.enc.Encode(httpapi.DNSResponse{
			Success: false,
			Message: "failed to stop DNS server",
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Error("failed to stop DNS server", attr.String("error", err.Error()))
		return
	}

	go func() {
		err = e.UDP.Start(ctx)
	}()

	if err != nil {
		w.WriteHeader(500)
		response, _ := e.enc.Encode(httpapi.DNSResponse{
			Success: false,
			Message: "failed to start DNS server",
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Error("failed to start DNS server", attr.String("error", err.Error()))
		return
	}
	w.WriteHeader(200)
	response, _ := e.enc.Encode(httpapi.DNSResponse{
		Success: true,
		Message: "reloaded DNS server",
	})
	_, _ = w.Write(response)
	logx.From(ctx).Debug("started DNS server")
}
