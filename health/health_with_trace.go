package health

import (
	"context"
	"time"

	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/spanner"
)

type withTrace struct {
	r Repository
}

func WithTrace(r Repository) Repository {
	return withTrace{
		r: r,
	}
}

// Store will take in the number of records in the store and the time.Duration for a
// store.List operation, and return a StoreReport based off of this information
func (t withTrace) Store(ctx context.Context, length int, d time.Duration) *StoreReport {
	ctx, s := spanner.Start(ctx, "health.Store")
	defer s.End()
	s.Add(
		attr.Int("record_length", length),
		attr.New("duration", d.Milliseconds()),
	)

	report := t.r.Store(ctx, length, d)
	s.Add(attr.New("report", report))

	return report
}

// DNS will take in the address of the UDP server, the fallback DNS address (if set),
// and a store.Record, which are used to answer internal and external DNS queries as part
// of a health check; returning a DNSReport based off of this information
func (t withTrace) DNS(ctx context.Context, address string, fallback string, r *store.Record) *DNSReport {
	ctx, s := spanner.Start(ctx, "health.DNS")
	defer s.End()
	s.Add(
		attr.String("address", address),
		attr.String("fallback", fallback),
		attr.New("record", r),
	)

	report := t.r.DNS(ctx, address, fallback, r)
	s.Add(attr.New("report", report))

	return report
}

// HTTP will take the HTTP server's port so it can perform a HTTP request against one
// of its endpoints, and returning a HTTPReport based off of this information
func (t withTrace) HTTP(ctx context.Context, port int) *HTTPReport {
	ctx, s := spanner.Start(ctx, "health.HTTP")
	defer s.End()
	s.Add(attr.Int("port", port))

	report := t.r.HTTP(ctx, port)
	s.Add(attr.New("report", report))

	return report
}

// Merge will unite a StoreReport, DNSReport and HTTPReport, returning a Report which
// encapsulates these as well as an overall status for the service
func (t withTrace) Merge(ctx context.Context, sr *StoreReport, dr *DNSReport, hr *HTTPReport) *Report {
	ctx, s := spanner.Start(ctx, "health.Merge")
	defer s.End()
	s.Add(
		attr.New("store_report", sr),
		attr.New("dns_report", dr),
		attr.New("http_report", hr),
	)

	report := t.r.Merge(ctx, sr, dr, hr)
	s.Add(attr.New("report", report))

	return report
}
