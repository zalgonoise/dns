package service

import (
	"context"

	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/dns/health"
	"github.com/zalgonoise/x/spanner"
)

// StoreHealth uses the health.Repository to generate a health.StoreReport
func (t withTrace) StoreHealth(ctx context.Context) *health.StoreReport {
	ctx, s := spanner.Start(ctx, "service.StoreHealth")
	defer s.End()

	report := t.s.StoreHealth(ctx)
	s.Add(attr.New("report", report))

	return report
}

// DNSHealth uses the health.Repository to generate a health.DNSReport
func (t withTrace) DNSHealth(ctx context.Context) *health.DNSReport {
	ctx, s := spanner.Start(ctx, "service.DNSHealth")
	defer s.End()

	report := t.s.DNSHealth(ctx)
	s.Add(attr.New("report", report))

	return report
}

// HTTPHealth uses the health.Repository to generate a health.HTTPReport
func (t withTrace) HTTPHealth(ctx context.Context) *health.HTTPReport {
	ctx, s := spanner.Start(ctx, "service.HTTPHealth")
	defer s.End()

	report := t.s.HTTPHealth(ctx)
	s.Add(attr.New("report", report))

	return report
}

// Health uses the health.Repository to generate a health.Report
func (t withTrace) Health(ctx context.Context) *health.Report {
	ctx, s := spanner.Start(ctx, "service.Health")
	defer s.End()

	report := t.s.Health(ctx)
	s.Add(attr.New("report", report))

	return report
}
