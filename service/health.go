package service

import (
	"context"
	"strings"
	"time"

	"github.com/zalgonoise/dns/health"
	"github.com/zalgonoise/dns/store"
)

// StoreHealth uses the health.Repository to generate a health.StoreReport
func (s *service) StoreHealth(ctx context.Context) *health.StoreReport {
	before := time.Now()
	r, err := s.store.List(ctx)
	after := time.Since(before)
	if err != nil {
		return s.health.Store(ctx, 0, 0)
	}
	return s.health.Store(ctx, len(r), after)
}

// DNSHealth uses the health.Repository to generate a health.DNSReport
func (s *service) DNSHealth(ctx context.Context) *health.DNSReport {
	var addr string

	r, err := s.store.List(ctx)
	if err != nil || len(r) == 0 {
		r = []*store.Record{nil}
	}

	if s.conf.DNS.FallbackDNS != "" {
		addr = strings.Split(s.conf.DNS.FallbackDNS, ",")[0]
	}

	return s.health.DNS(
		ctx,
		s.conf.DNS.Address,
		addr,
		r[0],
	)

}

// HTTPHealth uses the health.Repository to generate a health.HTTPReport
func (s *service) HTTPHealth(ctx context.Context) *health.HTTPReport {
	return s.health.HTTP(ctx, s.conf.HTTP.Port)
}

// Health uses the health.Repository to generate a health.Report
func (s *service) Health(ctx context.Context) *health.Report {
	return s.health.Merge(
		ctx,
		s.StoreHealth(ctx),
		s.DNSHealth(ctx),
		s.HTTPHealth(ctx),
	)
}
