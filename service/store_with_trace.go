package service

import (
	"context"

	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/spanner"
)

// AddRecord uses the store.Repository to create a DNS Record
func (t withTrace) AddRecord(ctx context.Context, r *store.Record) error {
	ctx, s := spanner.Start(ctx, "service.AddRecord")
	defer s.End()
	s.Add(attr.New("record", r))

	err := t.s.AddRecord(ctx, r)
	if err != nil {
		s.Event("error adding record", attr.New("error", err.Error()))
	}

	return err
}

// AddRecords uses the store.Repository to create a set of DNS Records
func (t withTrace) AddRecords(ctx context.Context, records ...*store.Record) error {
	ctx, s := spanner.Start(ctx, "service.AddRecords")
	defer s.End()
	s.Add(attr.New("records", records), attr.Int("len", len(records)))

	err := t.s.AddRecords(ctx, records...)
	if err != nil {
		s.Event("error adding records", attr.New("error", err.Error()))
	}

	return err
}

// ListRecord uses the store.Repository to return all DNS Records
func (t withTrace) ListRecords(ctx context.Context) ([]*store.Record, error) {
	ctx, s := spanner.Start(ctx, "service.ListRecords")
	defer s.End()

	records, err := t.s.ListRecords(ctx)
	if err != nil {
		s.Event("error listing records", attr.New("error", err.Error()))
	} else {
		s.Add(attr.Int("len", len(records)), attr.New("records", records))
	}

	return records, err
}

// GetRecordByDomain uses the store.Repository to return the DNS Record associated with
// the domain name and record type found in store.Record `r`
func (t withTrace) GetRecordByTypeAndDomain(ctx context.Context, rtype, domain string) (*store.Record, error) {
	ctx, s := spanner.Start(ctx, "service.GetRecordByTypeAndDomain")
	defer s.End()
	s.Add(
		attr.String("record_type", rtype),
		attr.String("domain", domain),
	)

	r, err := t.s.GetRecordByTypeAndDomain(ctx, rtype, domain)
	if err != nil {
		s.Event("error fetching record", attr.New("error", err.Error()))
	} else {
		s.Add(attr.New("record", r))
	}

	return r, err
}

// GetRecordByDomain uses the store.Repository to return the DNS Records associated with
// the IP address found in store.Record `r`
func (t withTrace) GetRecordByAddress(ctx context.Context, address string) ([]*store.Record, error) {
	ctx, s := spanner.Start(ctx, "service.GetRecordByAddress")
	defer s.End()
	s.Add(attr.String("address", address))

	records, err := t.s.GetRecordByAddress(ctx, address)
	if err != nil {
		s.Event("error fetching records", attr.New("error", err.Error()))
	} else {
		s.Add(attr.Int("len", len(records)), attr.New("records", records))
	}

	return records, err
}

// UpdateRecord uses the store.Repository to update the record with domain name `domain`,
// based on the data provided in store.Record `r`
func (t withTrace) UpdateRecord(ctx context.Context, domain string, r *store.Record) error {
	ctx, s := spanner.Start(ctx, "service.UpdateRecord")
	defer s.End()
	s.Add(
		attr.String("target_domain", domain),
		attr.New("record", r),
	)

	err := t.s.UpdateRecord(ctx, domain, r)
	if err != nil {
		s.Event("error updating records", attr.New("error", err.Error()))
	}

	return err
}

// DeleteRecord uses the store.Repository to remove the store.Record based on input `r`
func (t withTrace) DeleteRecord(ctx context.Context, r *store.Record) error {
	ctx, s := spanner.Start(ctx, "service.DeleteRecord")
	defer s.End()
	s.Add(attr.New("record", r))

	err := t.s.DeleteRecord(ctx, r)
	if err != nil {
		s.Event("error deleting records", attr.New("error", err.Error()))
	}

	return err
}
