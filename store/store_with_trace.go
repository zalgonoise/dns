package store

import (
	"context"

	"github.com/zalgonoise/logx/attr"
	"github.com/zalgonoise/x/spanner"
)

type withTrace struct {
	r Repository
}

func WithTrace(r Repository) Repository {
	return withTrace{
		r: r,
	}
}

// Create will add a new entry in they key-value store to include a
// new Record, returning an error
func (t withTrace) Create(ctx context.Context, records ...*Record) error {
	ctx, s := spanner.Start(ctx, "store.Create",
		attr.Int("records_length", len(records)),
		attr.New("records", records),
	)
	defer s.End()

	err := t.r.Create(ctx, records...)
	if err != nil {
		s.Event("error creating records", attr.New("error", err.Error()))
	}

	return err
}

// List will fetch all records in the key-value store
func (t withTrace) List(ctx context.Context) ([]*Record, error) {
	ctx, s := spanner.Start(ctx, "store.List")
	defer s.End()

	records, err := t.r.List(ctx)
	if err != nil {
		s.Event("error listing record", attr.New("error", err.Error()))
	} else {
		s.Add(
			attr.Int("records_length", len(records)),
			attr.New("records", records),
		)
	}

	return records, err
}

// FindByTypeAndDomain will fetch an address based on its address and type strings
//
// FindByTypeAndDomain(ctx, "A", "service.mydomain") -> { "127.0.0.1", nil }
func (t withTrace) FindByTypeAndDomain(ctx context.Context, rtype string, domain string) (*Record, error) {
	ctx, s := spanner.Start(ctx, "store.FindByTypeAndDomain",
		attr.String("record_type", rtype),
		attr.String("domain", domain),
	)
	defer s.End()

	r, err := t.r.FindByTypeAndDomain(ctx, rtype, domain)
	if err != nil {
		s.Event("error finding records", attr.New("error", err.Error()))
	} else {
		s.Add(attr.New("record", r))
	}

	return r, err
}

// FilterByDomain will fetch an address based on its address for all types
//
// FilterByDomain(ctx, "service.mydomain") -> { ["127.0.0.1"], nil }
func (t withTrace) FilterByDomain(ctx context.Context, domain string) ([]*Record, error) {
	ctx, s := spanner.Start(ctx, "store.FilterByDomain",
		attr.String("domain", domain),
	)
	defer s.End()

	records, err := t.r.FilterByDomain(ctx, domain)
	if err != nil {
		s.Event("error filtering records", attr.New("error", err.Error()))
	} else {
		s.Add(
			attr.Int("records_length", len(records)),
			attr.New("records", records),
		)
	}

	return records, err
}

// FilterByDest will fetch an address based on its target IP
//
// FilterByDest(ctx, "127.0.0.1") -> { ["service.mydomain"], nil }
func (t withTrace) FilterByDest(ctx context.Context, address string) ([]*Record, error) {
	ctx, s := spanner.Start(ctx, "store.FilterByDest",
		attr.String("address", address),
	)
	defer s.End()

	records, err := t.r.FilterByDest(ctx, address)
	if err != nil {
		s.Event("error filtering records", attr.New("error", err.Error()))
	} else {
		s.Add(
			attr.Int("records_length", len(records)),
			attr.New("records", records),
		)
	}

	return records, err
}

// Update will modify an existing record by targetting its domain string,
// and by supplying a new version of the Record to update. Returns an error
func (t withTrace) Update(ctx context.Context, domain string, r *Record) error {
	ctx, s := spanner.Start(ctx, "store.Update",
		attr.String("target_domain", domain),
		attr.New("record", r),
	)
	defer s.End()

	if domain != r.Name {
		s.Event("different domain name -- removing former entry")
	}

	err := t.r.Update(ctx, domain, r)
	if err != nil {
		s.Event("error filtering records", attr.New("error", err.Error()))
	}

	return err
}

// DeleteByAddress removes all records with IP address `addr`
func (t withTrace) DeleteByAddress(ctx context.Context, addr string) error {
	ctx, s := spanner.Start(ctx, "store.DeleteByAddress",
		attr.String("address", addr),
	)
	defer s.End()

	err := t.r.DeleteByAddress(ctx, addr)
	if err != nil {
		s.Event("error deleting records", attr.New("error", err.Error()))
	}

	return err
}

// DeleteByDomain removes all records with domain name `name`
func (t withTrace) DeleteByDomain(ctx context.Context, name string) error {
	ctx, s := spanner.Start(ctx, "store.DeleteByDomain",
		attr.String("domain", name),
	)
	defer s.End()

	err := t.r.DeleteByDomain(ctx, name)
	if err != nil {
		s.Event("error deleting records", attr.New("error", err.Error()))
	}

	return err
}

// DeleteByTypeAndDomain removes all records with record type `rtype` and domain name `name`
func (t withTrace) DeleteByTypeAndDomain(ctx context.Context, rtype, name string) error {
	ctx, s := spanner.Start(ctx, "store.DeleteByTypeAndDomain",
		attr.String("record_type", rtype),
		attr.String("domain", name),
	)
	defer s.End()

	err := t.r.DeleteByTypeAndDomain(ctx, rtype, name)
	if err != nil {
		s.Event("error deleting records", attr.New("error", err.Error()))
	}

	return err
}
