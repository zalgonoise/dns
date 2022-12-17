package service

import (
	"context"

	"github.com/zalgonoise/dns/health"
)

// StoreHealth uses the health.Repository to generate a health.StoreReport
func (l withLogger) StoreHealth(ctx context.Context) *health.StoreReport {
	return l.s.StoreHealth(ctx)
}

// DNSHealth uses the health.Repository to generate a health.DNSReport
func (l withLogger) DNSHealth(ctx context.Context) *health.DNSReport {
	return l.s.DNSHealth(ctx)
}

// HTTPHealth uses the health.Repository to generate a health.HTTPReport
func (l withLogger) HTTPHealth(ctx context.Context) *health.HTTPReport {
	return l.s.HTTPHealth(ctx)
}

// Health uses the health.Repository to generate a health.Report
func (l withLogger) Health(ctx context.Context) *health.Report {
	return l.s.Health(ctx)
}
