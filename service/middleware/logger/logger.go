package logger

import (
	"context"
	"time"

	dnsr "github.com/miekg/dns"

	"github.com/zalgonoise/dns/health"
	"github.com/zalgonoise/dns/service"
	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/attr"
)

// LoggedService will wrap a service.Service with a logger and
// register the incoming events, as well as their outcome
type LoggedService struct {
	svc    service.Service
	logger logx.Logger
}

// LogService will return a LoggedService in the form of a service.Service,
// by wraping an input service.Service `svc` with log.Logger `logger`
func LogService(svc service.Service, logger logx.Logger) service.Service {
	return &LoggedService{
		svc:    svc,
		logger: logger,
	}
}

// AddRecord uses the store.Repository to create a DNS Record
func (s *LoggedService) AddRecord(ctx context.Context, r *store.Record) error {
	s.logger.Debug("AddRecord request",
		attr.String("module", "service:records:create"),
		attr.New("data", []attr.Attr{
			attr.New("input", r),
		}),
	)

	err := s.svc.AddRecord(ctx, r)
	if err != nil {
		s.logger.Warn("AddRecord error",
			attr.String("module", "service:records:create"),
			attr.New("data", []attr.Attr{
				attr.String("error", err.Error()),
			}),
		)
		return err
	}
	s.logger.Debug("AddRecord success",
		attr.String("module", "service:records:create"),
	)
	return nil
}

// AddRecords uses the store.Repository to create a set of DNS Records
func (s *LoggedService) AddRecords(ctx context.Context, rs ...*store.Record) error {
	s.logger.Debug("AddRecords request",
		attr.String("module", "service:records:create"),
		attr.New("data", []attr.Attr{
			attr.New("input", rs),
			attr.Int("len", len(rs)),
		}),
	)

	err := s.svc.AddRecords(ctx, rs...)
	if err != nil {
		s.logger.Warn("AddRecords error",
			attr.String("module", "service:records:create"),
			attr.New("data", []attr.Attr{
				attr.String("error", err.Error()),
			}),
		)

		return err
	}
	s.logger.Debug("AddRecords success",
		attr.String("module", "service:records:create"),
	)
	return nil
}

// ListRecord uses the store.Repository to return all DNS Records
func (s *LoggedService) ListRecords(ctx context.Context) ([]*store.Record, error) {
	s.logger.Debug("ListRecords request",
		attr.String("module", "service:records:list"),
	)

	rs, err := s.svc.ListRecords(ctx)
	if err != nil {
		s.logger.Warn("ListRecords error",
			attr.String("module", "service:records:list"),
			attr.New("data", []attr.Attr{
				attr.String("error", err.Error()),
			}),
		)
		return rs, err
	}
	s.logger.Debug("ListRecords success",
		attr.String("module", "service:records:list"),
		attr.New("data", []attr.Attr{
			attr.New("output", rs),
			attr.Int("len", len(rs)),
		}),
	)
	return rs, nil
}

// GetRecordByDomain uses the store.Repository to return the DNS Record associated with
// the domain name and record type found in store.Record `r`
func (s *LoggedService) GetRecordByTypeAndDomain(ctx context.Context, rtype, domain string) (*store.Record, error) {
	s.logger.Debug("GetRecordByTypeAndDomain request",
		attr.String("module", "service:records:read"),
		attr.New("data", []attr.Attr{
			attr.New("input", []attr.Attr{
				attr.String("type", rtype),
				attr.String("domain", domain),
			}),
		}),
	)

	out, err := s.svc.GetRecordByTypeAndDomain(ctx, rtype, domain)
	if err != nil {
		s.logger.Warn("GetRecordByTypeAndDomain error",
			attr.String("module", "service:records:read"),
			attr.New("data", []attr.Attr{
				attr.String("error", err.Error()),
			}),
		)
		return out, err
	}

	s.logger.Debug("GetRecordByTypeAndDomain success",
		attr.String("module", "service:records:read"),
		attr.New("data", []attr.Attr{
			attr.New("output", out),
		}),
	)
	return out, nil
}

// GetRecordByDomain uses the store.Repository to return the DNS Records associated with
// the IP address found in store.Record `r`
func (s *LoggedService) GetRecordByAddress(ctx context.Context, address string) ([]*store.Record, error) {
	s.logger.Debug("GetRecordByAddress request",
		attr.String("module", "service:records:list"),
		attr.New("data", []attr.Attr{
			attr.New("input", []attr.Attr{
				attr.String("address", address),
			}),
		}),
	)

	rs, err := s.svc.GetRecordByAddress(ctx, address)
	if err != nil {
		s.logger.Warn("GetRecordByAddress error",
			attr.String("module", "service:records:list"),
			attr.New("data", []attr.Attr{
				attr.String("error", err.Error()),
			}),
		)
		return rs, err
	}

	s.logger.Debug("GetRecordByAddress success",
		attr.String("module", "service:records:list"),
		attr.New("data", []attr.Attr{
			attr.New("output", rs),
			attr.Int("len", len(rs)),
		}),
	)
	return rs, nil
}

// UpdateRecord uses the store.Repository to update the record with domain name `domain`,
// based on the data provided in store.Record `r`
func (s *LoggedService) UpdateRecord(ctx context.Context, domain string, r *store.Record) error {
	s.logger.Debug("UpdateRecord request",
		attr.String("module", "service:records:update"),
		attr.New("data", []attr.Attr{
			attr.New("input", []attr.Attr{
				attr.String("target", domain),
				attr.New("record", r),
			}),
		}),
	)

	err := s.svc.UpdateRecord(ctx, domain, r)
	if err != nil {
		s.logger.Warn("UpdateRecord error",
			attr.String("module", "service:records:update"),
			attr.New("data", []attr.Attr{
				attr.String("error", err.Error()),
			}),
		)
		return err
	}

	s.logger.Debug("UpdateRecord success",
		attr.String("module", "service:records:update"),
	)
	return nil
}

// DeleteRecord uses the store.Repository to remove the store.Record based on input `r`
func (s *LoggedService) DeleteRecord(ctx context.Context, r *store.Record) error {
	s.logger.Debug("DeleteRecord request",
		attr.String("module", "service:records:delete"),
		attr.New("data", []attr.Attr{
			attr.New("input", r),
		}),
	)

	err := s.svc.DeleteRecord(ctx, r)
	if err != nil {
		s.logger.Warn("DeleteRecord error",
			attr.String("module", "service:records:delete"),
			attr.New("data", []attr.Attr{
				attr.String("error", err.Error()),
			}),
		)
		return err
	}
	s.logger.Debug("DeleteRecord success",
		attr.String("module", "service:records:delete"),
	)
	return nil
}

// AnswerDNS uses the dns.Repository to reply to the dns.Msg `m` with the answer
// in store.Record `r`
func (s *LoggedService) AnswerDNS(ctx context.Context, r *store.Record, m *dnsr.Msg) {
	s.logger.Debug("AnswerDNS request",
		attr.String("module", "service:dns:answer"),
	)

	s.svc.AnswerDNS(ctx, r, m)
	go func() {
		time.Sleep(5 * time.Millisecond)
		s.logger.Debug("AnswerDNS response",
			attr.String("module", "service:dns:answer"),
			attr.New("data", []attr.Attr{
				attr.New("output", r),
			}),
		)
	}()
}

// StoreHealth uses the health.Repository to generate a health.StoreReport
func (s *LoggedService) StoreHealth(ctx context.Context) *health.StoreReport {
	s.logger.Debug("StoreHealth request",
		attr.String("module", "service:health:store"),
	)

	r := s.svc.StoreHealth(ctx)
	go func() {
		time.Sleep(5 * time.Millisecond)
		s.logger.Debug("StoreHealth response",
			attr.String("module", "service:health:store"),
			attr.New("data", []attr.Attr{
				attr.New("output", r),
			}),
		)
	}()
	return r
}

// DNSHealth uses the health.Repository to generate a health.DNSReport
func (s *LoggedService) DNSHealth(ctx context.Context) *health.DNSReport {
	s.logger.Debug("DNSHealth request",
		attr.String("module", "service:health:dns"),
	)

	r := s.svc.DNSHealth(ctx)
	go func() {
		time.Sleep(5 * time.Millisecond)
		s.logger.Debug("DNSHealth response",
			attr.String("module", "service:health:dns"),
			attr.New("data", []attr.Attr{
				attr.New("output", r),
			}),
		)
	}()
	return r
}

// HTTPHealth uses the health.Repository to generate a health.HTTPReport
func (s *LoggedService) HTTPHealth(ctx context.Context) *health.HTTPReport {
	s.logger.Debug("HTTPHealth request",
		attr.String("module", "service:health:http"),
	)

	r := s.svc.HTTPHealth(ctx)
	go func() {
		time.Sleep(5 * time.Millisecond)
		s.logger.Debug("HTTPHealth response",
			attr.String("module", "service:health:http"),
			attr.New("data", []attr.Attr{
				attr.New("output", r),
			}),
		)
	}()
	return r
}

// Health uses the health.Repository to generate a health.Report
func (s *LoggedService) Health(ctx context.Context) *health.Report {
	s.logger.Debug("MergeHealth request",
		attr.String("module", "service:health"),
	)

	r := s.svc.Health(ctx)
	go func() {
		time.Sleep(5 * time.Millisecond)
		s.logger.Debug("MergeHealth response",
			attr.String("module", "service:health"),
			attr.New("data", []attr.Attr{
				attr.New("output", r),
			}),
		)
	}()
	return r
}
