package service

import (
	"context"

	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/dns/store"
)

// AddRecord uses the store.Repository to create a DNS Record
func (l withLogger) AddRecord(ctx context.Context, r *store.Record) error {
	err := l.s.AddRecord(ctx, r)
	if err != nil {
		l.log.Error("failed to add record",
			attr.String("error", err.Error()),
			attr.New("input", r),
		)
	}

	return err
}

// AddRecords uses the store.Repository to create a set of DNS Records
func (l withLogger) AddRecords(ctx context.Context, records ...*store.Record) error {
	err := l.s.AddRecords(ctx, records...)
	if err != nil {
		l.log.Error("failed to add records",
			attr.String("error", err.Error()),
			attr.New("input", records),
		)
	}

	return err
}

// ListRecord uses the store.Repository to return all DNS Records
func (l withLogger) ListRecords(ctx context.Context) ([]*store.Record, error) {
	records, err := l.s.ListRecords(ctx)
	if err != nil {
		l.log.Error("failed to list records",
			attr.String("error", err.Error()),
		)
	}

	return records, err
}

// GetRecordByDomain uses the store.Repository to return the DNS Record associated with
// the domain name and record type found in store.Record `r`
func (l withLogger) GetRecordByTypeAndDomain(ctx context.Context, rtype, domain string) (*store.Record, error) {
	r, err := l.s.GetRecordByTypeAndDomain(ctx, rtype, domain)
	if err != nil {
		l.log.Error("failed to get record by type and domain",
			attr.String("error", err.Error()),
			attr.String("record_type", rtype),
			attr.String("domain", domain),
		)
	}

	return r, err
}

// GetRecordByDomain uses the store.Repository to return the DNS Records associated with
// the IP address found in store.Record `r`
func (l withLogger) GetRecordByAddress(ctx context.Context, address string) ([]*store.Record, error) {
	records, err := l.s.GetRecordByAddress(ctx, address)
	if err != nil {
		l.log.Error("failed to get records by address",
			attr.String("error", err.Error()),
			attr.String("address", address),
		)
	}

	return records, err
}

// UpdateRecord uses the store.Repository to update the record with domain name `domain`,
// based on the data provided in store.Record `r`
func (l withLogger) UpdateRecord(ctx context.Context, domain string, r *store.Record) error {
	err := l.s.UpdateRecord(ctx, domain, r)
	if err != nil {
		l.log.Error("failed to update record",
			attr.String("error", err.Error()),
			attr.New("input", r),
		)
	}

	return err
}

// DeleteRecord uses the store.Repository to remove the store.Record based on input `r`
func (l withLogger) DeleteRecord(ctx context.Context, r *store.Record) error {
	err := l.s.DeleteRecord(ctx, r)
	if err != nil {
		l.log.Error("failed to delete record",
			attr.String("error", err.Error()),
			attr.New("input", r),
		)
	}

	return err
}
