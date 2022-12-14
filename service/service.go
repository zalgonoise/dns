package service

import (
	"context"
	"errors"

	dnsr "github.com/miekg/dns"
	"github.com/zalgonoise/dns/cmd/config"
	"github.com/zalgonoise/dns/dns"
	"github.com/zalgonoise/dns/health"
	"github.com/zalgonoise/dns/health/simplehealth"
	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/logx"
)

var (
	ErrNoAddr      = errors.New("no IP address provided")
	ErrNoName      = errors.New("no domain name provided")
	ErrNoType      = errors.New("no DNS record type provided")
	ErrEmtpyRecord = errors.New("record cannot be empty")
)

// Service interface packs the different sets of operations supported by this
// DNS service.
//
// While the repositories the service contains will work independently one from the other,
// the Service interface is there to *glue them* together and to allow inter-repository
// operations. It is one level below the transport layer.
type Service interface {
	StoreService
	DNSService
	HealthService
}

// StoreService interface joins the set of methods leveraging the store.Repository
type StoreService interface {
	// AddRecord uses the store.Repository to create a DNS Record
	AddRecord(ctx context.Context, r *store.Record) error
	// AddRecords uses the store.Repository to create a set of DNS Records
	AddRecords(ctx context.Context, rs ...*store.Record) error
	// ListRecord uses the store.Repository to return all DNS Records
	ListRecords(ctx context.Context) ([]*store.Record, error)
	// GetRecordByDomain uses the store.Repository to return the DNS Record associated with
	// the domain name and record type found in store.Record `r`
	GetRecordByTypeAndDomain(ctx context.Context, rtype, domain string) (*store.Record, error)
	// GetRecordByDomain uses the store.Repository to return the DNS Records associated with
	// the IP address found in store.Record `r`
	GetRecordByAddress(ctx context.Context, address string) ([]*store.Record, error)
	// UpdateRecord uses the store.Repository to update the record with domain name `domain`,
	// based on the data provided in store.Record `r`
	UpdateRecord(ctx context.Context, domain string, r *store.Record) error
	// DeleteRecord uses the store.Repository to remove the store.Record based on input `r`
	DeleteRecord(ctx context.Context, r *store.Record) error
}

// DNSService interface joins the set of methods leveraging the dns.Repository
type DNSService interface {
	// AnswerDNS uses the dns.Repository to reply to the dns.Msg `m` with the answer
	// in store.Record `r`
	AnswerDNS(ctx context.Context, r *store.Record, m *dnsr.Msg)
}

// HealthService interface joins the set of methods leveraging the health.Repository
type HealthService interface {
	// StoreHealth uses the health.Repository to generate a health.StoreReport
	StoreHealth(ctx context.Context) *health.StoreReport
	// DNSHealth uses the health.Repository to generate a health.DNSReport
	DNSHealth(ctx context.Context) *health.DNSReport
	// HTTPHealth uses the health.Repository to generate a health.HTTPReport
	HTTPHealth(ctx context.Context) *health.HTTPReport
	// Health uses the health.Repository to generate a health.Report
	Health(ctx context.Context) *health.Report
}

// StoreService interface joins the store.Repository-derived methods
// with the health.Repository ones. This is done to provide more granular scope of which
// operations can a certain module access
type StoreWithHealth interface {
	StoreService
	HealthService
}

// Answering interface joins the needed store.Repository-derived method
// with the dns.Repository one for answering DNS queries. This is done to provide more
//
//	granular scope of which operations can a certain module access
type Answering interface {
	GetRecordByTypeAndDomain(context.Context, string, string) (*store.Record, error)
	AnswerDNS(context.Context, *store.Record, *dnsr.Msg)
}

type service struct {
	dns    dns.Repository
	store  store.Repository
	health health.Repository
	conf   *config.Config
}

type withTrace struct {
	s Service
}

type withLogger struct {
	s   Service
	log logx.Logger
}

// New will create a Service based on the input dns.Repository, store.Repository,
// health.Repository and configuration. These elements are required in order to initialize
// the Service with all elements and settings needed (such as HTTP port / UDP address)
func New(
	dnsR dns.Repository,
	storeR store.Repository,
	healthR health.Repository,
	conf *config.Config,
) Service {
	if dnsR == nil {
		dnsR = dns.Unimplemented()
	}
	if storeR == nil {
		storeR = store.Unimplemented()
	}
	if healthR == nil {
		healthR = simplehealth.New()
	}

	return &service{
		dns:    dnsR,
		store:  storeR,
		health: healthR,
		conf:   conf,
	}
}

// WithTrace wraps the input Service `s` with tracing instrumentation
func WithTrace(s Service) Service {
	return withTrace{
		s: s,
	}
}

// WithLogger wraps the input Service `s` with the input logx.Logger `logger`
func WithLogger(s Service, logger logx.Logger) Service {
	return withLogger{
		s:   s,
		log: logger,
	}
}
