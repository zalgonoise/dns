package store

import (
	"context"
	"errors"
)

var (
	ErrUnimplemented error = errors.New("unimplemented DNS Record Store")
)

type unimplemented struct{}

// Create implements the store.Repository interface
func (u unimplemented) Create(ctx context.Context, rs ...*Record) error {
	return ErrUnimplemented
}

// List implements the store.Repository interface
func (u unimplemented) List(ctx context.Context) ([]*Record, error) {
	return nil, ErrUnimplemented
}

// FindByTypeAndDomain implements the store.Repository interface
func (u unimplemented) FindByTypeAndDomain(ctx context.Context, rtype, domain string) (*Record, error) {
	return nil, ErrUnimplemented
}

// FilterByDomain implements the store.Repository interface
func (u unimplemented) FilterByDomain(ctx context.Context, domain string) ([]*Record, error) {
	return nil, ErrUnimplemented
}

// FilterByDest implements the store.Repository interface
func (u unimplemented) FilterByDest(ctx context.Context, addr string) ([]*Record, error) {
	return nil, ErrUnimplemented
}

// Update implements the store.Repository interface
func (u unimplemented) Update(ctx context.Context, domain string, r *Record) error {
	return ErrUnimplemented
}

// DeleteByAddress removes all records with IP address `addr`
func (u unimplemented) DeleteByAddress(ctx context.Context, addr string) error {
	return ErrUnimplemented
}

// DeleteByDomain removes all records with domain name `name`
func (u unimplemented) DeleteByDomain(ctx context.Context, name string) error {
	return ErrUnimplemented
}

// DeleteByTypeAndDomain removes all records with record type `rtype` and domain name `name`
func (u unimplemented) DeleteByTypeAndDomain(ctx context.Context, rtype, name string) error {
	return ErrUnimplemented
}

// Unimplemented returns an unimplemented (and invalid) store.Repository
func Unimplemented() unimplemented {
	return unimplemented{}
}
