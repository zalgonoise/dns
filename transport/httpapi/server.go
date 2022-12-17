package httpapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/spanner"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type server struct {
	ep   HTTPAPI
	port int
	srv  *http.Server
}

func NewServer(api HTTPAPI, port int) Server {
	mux := http.NewServeMux()
	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}
	srv := &server{
		ep:   api,
		port: port,
		srv:  httpSrv,
	}
	mux.HandleFunc("/dns/start", srv.ep.StartDNS)
	mux.HandleFunc("/dns/stop", srv.ep.StopDNS)
	mux.HandleFunc("/dns/reload", srv.ep.ReloadDNS)
	mux.HandleFunc("/records/add", srv.ep.AddRecord)
	mux.HandleFunc("/records", srv.ep.ListRecords)
	mux.HandleFunc("/records/getAddress", srv.ep.GetRecordByDomain)
	mux.HandleFunc("/records/getDomains", srv.ep.GetRecordByAddress)
	mux.HandleFunc("/records/update", srv.ep.UpdateRecord)
	mux.HandleFunc("/records/delete", srv.ep.DeleteRecord)
	mux.HandleFunc("/health", srv.ep.Health)

	return srv
}

func (s *server) Start(ctx context.Context) error {
	ctx, span := spanner.Start(ctx, "http.Start")
	defer span.End()

	err := s.srv.ListenAndServe()
	if err != nil {
		span.Event("failed to start HTTP server", attr.New("error", err.Error()))
		return err
	}
	return nil
}

func (s *server) Stop(ctx context.Context) error {
	ctx, span := spanner.Start(ctx, "http.Stop")
	defer span.End()

	rw := &responseWriter{}
	s.ep.StopDNS(rw, &http.Request{})

	if rw.header != 200 {
		span.Event("issue stopping UDP server",
			attr.Int("status", rw.header),
			attr.String("response", rw.response),
		)
	}

	err := s.srv.Shutdown(ctx)
	if err != nil {
		span.Event("failed to stop HTTP server", attr.New("error", err.Error()))
		return err
	}
	return nil
}
