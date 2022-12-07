package file

import (
	"context"
	"fmt"

	"github.com/zalgonoise/dns/store"
)

// Create implements the store.Repository interface
//
// It will call the in-memory store's method of the same signature, while deferring
// a `Sync()` call to ensure the records file is up-to-date
func (f *FileStore) Create(ctx context.Context, rs ...*store.Record) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	err := f.store.Create(ctx, rs...)
	if err != nil {
		return fmt.Errorf("failed to create record: %w", err)
	}
	err = f.sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync store to file: %w", err)
	}
	return nil
}

// List implements the store.Repository interface
//
// It will call the in-memory store's method of the same signature
func (f *FileStore) List(ctx context.Context) ([]*store.Record, error) {
	rs, err := f.store.List(ctx)
	if err != nil {
		return rs, fmt.Errorf("failed to list records: %w", err)
	}
	return rs, nil
}

// FindByTypeAndDomain implements the store.Repository interface
//
// It will call the in-memory store's method of the same signature
func (f *FileStore) FindByTypeAndDomain(ctx context.Context, rtype, domain string) (*store.Record, error) {
	r, err := f.store.FindByTypeAndDomain(ctx, rtype, domain)
	if err != nil {
		return r, fmt.Errorf("failed to get record by domain: %w", err)
	}
	return r, nil
}

// FilterByDomain implements the store.Repository interface
//
// It will call the in-memory store's method of the same signature
func (f *FileStore) FilterByDomain(ctx context.Context, domain string) ([]*store.Record, error) {
	rs, err := f.store.FilterByDomain(ctx, domain)
	if err != nil {
		return rs, fmt.Errorf("failed to get record by domain name: %w", err)
	}
	return rs, nil
}

// FilterByDest implements the store.Repository interface
//
// It will call the in-memory store's method of the same signature
func (f *FileStore) FilterByDest(ctx context.Context, addr string) ([]*store.Record, error) {
	rs, err := f.store.FilterByDest(ctx, addr)
	if err != nil {
		return rs, fmt.Errorf("failed to get record by address: %w", err)
	}
	return rs, nil
}

// Update implements the store.Repository interface
//
// It will call the in-memory store's method of the same signature, while deferring
// a `Sync()` call to ensure the records file is up-to-date
func (f *FileStore) Update(ctx context.Context, domain string, r *store.Record) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	err := f.store.Update(ctx, domain, r)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}
	err = f.sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync store to file: %w", err)
	}
	return nil
}

// DeleteByAddress removes all records with IP address `addr`
func (f *FileStore) DeleteByAddress(ctx context.Context, addr string) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	err := f.store.DeleteByAddress(ctx, addr)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}
	err = f.sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync store to file: %w", err)
	}
	return nil
}

// DeleteByDomain removes all records with domain name `name`
func (f *FileStore) DeleteByDomain(ctx context.Context, name string) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	err := f.store.DeleteByDomain(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}
	err = f.sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync store to file: %w", err)
	}
	return nil
}

// DeleteByTypeAndDomain removes all records with record type `rtype` and domain name `name`
func (f *FileStore) DeleteByTypeAndDomain(ctx context.Context, rtype, name string) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	err := f.store.DeleteByTypeAndDomain(ctx, rtype, name)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}
	err = f.sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync store to file: %w", err)
	}
	return nil
}
