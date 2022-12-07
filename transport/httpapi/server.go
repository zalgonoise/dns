package httpapi

import (
	"context"
	"fmt"
	"net/http"
)

type Server interface {
	Start() error
	Stop() error
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

func (s *server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *server) Stop() error {
	s.ep.StopDNS(&responseWriter{}, &http.Request{})
	return s.srv.Shutdown(context.Background())

}
