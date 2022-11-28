package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/attr"
)

// AddRecord uses the store.Repository to create a DNS Record
func (s *service) AddRecord(ctx context.Context, r *store.Record) error {
	logx.From(ctx).Debug("creating record",
		attr.String("action", "service:store:add"),
		attr.New("record", r),
	)

	err := s.store.Create(ctx, r)
	if err != nil {
		logx.From(ctx).Error("failed to add record", attr.String("error", err.Error()))
		return fmt.Errorf("couldn't add target record: %w", err)
	}

	logx.From(ctx).Debug("added record")
	return nil
}

// AddRecords uses the store.Repository to create a set of DNS Records
func (s *service) AddRecords(ctx context.Context, rs ...*store.Record) error {
	logx.From(ctx).Debug("creating records",
		attr.String("action", "service:store:add"),
		attr.New("records", rs),
		attr.Int("len", len(rs)),
	)

	err := s.store.Create(ctx, rs...)
	if err != nil {
		logx.From(ctx).Error("failed to add records", attr.String("error", err.Error()))
		return fmt.Errorf("couldn't add target records: %w", err)
	}

	logx.From(ctx).Debug("added records")
	return nil
}

// ListRecord uses the store.Repository to return all DNS Records
func (s *service) ListRecords(ctx context.Context) ([]*store.Record, error) {
	logx.From(ctx).Debug("listing record",
		attr.String("action", "service:store:list"),
	)

	rs, err := s.store.List(ctx)
	if err != nil {
		logx.From(ctx).Error("failed to list records", attr.String("error", err.Error()))
		return nil, fmt.Errorf("couldn't list any records: %w", err)
	}

	logx.From(ctx).Debug("listed records", attr.Int("len", len(rs)))
	return rs, nil
}

// GetRecordByDomain uses the store.Repository to return the DNS Record associated with
// the domain name and record type found in store.Record `r`
//
// Returns a NoName error if no domain name is provided
// Returns a NoType if no record type is provided
func (s *service) GetRecordByTypeAndDomain(ctx context.Context, rtype, domain string) (*store.Record, error) {
	logx.From(ctx).Debug("getting record",
		attr.String("action", "service:store:get"),
		attr.String("type", rtype),
		attr.String("domain", domain),
	)

	if domain == "" {
		logx.From(ctx).Error("failed to get record", attr.String("error", ErrNoName.Error()))
		return nil, ErrNoName
	}
	if rtype == "" {
		logx.From(ctx).Error("failed to get record", attr.String("error", ErrNoType.Error()))
		return nil, ErrNoType
	}

	r, err := s.store.FilterByTypeAndDomain(ctx, rtype, domain)
	if err != nil {
		logx.From(ctx).Error("failed to get record", attr.String("error", err.Error()))
		return nil, fmt.Errorf("couldn't fetch target record: %w", err)
	}

	logx.From(ctx).Debug("got record", attr.String("address", r.Addr))
	return r, nil
}

// GetRecordByDomain uses the store.Repository to return the DNS Records associated with
// the IP address found in store.Record `r`
//
// Returns a NoAddr error if no IP address is provided
func (s *service) GetRecordByAddress(ctx context.Context, ip string) ([]*store.Record, error) {
	logx.From(ctx).Debug("listing records",
		attr.String("action", "service:store:list"),
		attr.String("address", ip),
	)
	if ip == "" {
		logx.From(ctx).Error("failed to list records", attr.String("error", ErrNoAddr.Error()))
		return nil, ErrNoAddr
	}

	rs, err := s.store.FilterByDest(ctx, ip)
	if err != nil {
		logx.From(ctx).Error("failed to list records", attr.String("error", err.Error()))
		return nil, fmt.Errorf("couldn't fetch target records: %w", err)
	}
	if len(rs) == 0 {
		return rs, fmt.Errorf("no records in the store: %w", store.ErrZeroRecords)
	}

	logx.From(ctx).Debug("got record", attr.Int("len", len(rs)), attr.New("records", rs))
	return rs, nil
}

// UpdateRecord uses the store.Repository to update the record with domain name `domain`,
// based on the data provided in store.Record `r`
//
// Returns a NoAddr error if no IP address is provided
func (s *service) UpdateRecord(ctx context.Context, domain string, r *store.Record) error {
	logx.From(ctx).Debug("updating record",
		attr.String("action", "service:store:update"),
		attr.String("domain", domain),
		attr.New("record", r),
	)

	if domain == "" {
		logx.From(ctx).Error("failed to list records", attr.String("error", ErrNoName.Error()))
		return store.ErrNoName
	}
	if r.Type == "" {
		logx.From(ctx).Error("failed to list records", attr.String("error", ErrNoType.Error()))
		return store.ErrNoType
	}
	// if record is nil or empty, request is parsed as a delete operation
	if r == nil || reflect.DeepEqual(r, &store.Record{}) {
		logx.From(ctx).Debug("target record is empty -- handling as delete action")
		return s.store.Delete(ctx, &store.Record{Name: domain})
	}

	_, err := s.store.FilterByTypeAndDomain(ctx, r.Type, domain)
	if err != nil && errors.Is(err, store.ErrDoesNotExist) {
		logx.From(ctx).Error("failed to find target record", attr.String("error", err.Error()))
		return fmt.Errorf("couldn't find target record: %w", err)
	}

	err = s.store.Update(ctx, domain, r)
	if err != nil {
		logx.From(ctx).Error("failed to update target record", attr.String("error", err.Error()))
		return fmt.Errorf("couldn't update target record: %w", err)
	}

	logx.From(ctx).Debug("updated record")
	return nil
}

// DeleteRecord uses the store.Repository to remove the store.Record based on input `r`
func (s *service) DeleteRecord(ctx context.Context, r *store.Record) error {
	logx.From(ctx).Debug("deleting record",
		attr.String("action", "service:store:delete"),
		attr.String("domain", r.Name),
	)
	err := s.store.Delete(ctx, r)
	if err != nil {
		logx.From(ctx).Error("failed to delete record", attr.String("error", err.Error()))
		return fmt.Errorf("couldn't delete target record: %w", err)
	}
	logx.From(ctx).Debug("deleted record")
	return nil
}
