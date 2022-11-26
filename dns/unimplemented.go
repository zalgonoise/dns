package dns

import (
	"context"
	"errors"

	dnsr "github.com/miekg/dns"
	"github.com/zalgonoise/dns/store"
)

var (
	ErrUnimplemented error = errors.New("unimplemented DNS server")
)

type unimplemented struct{}

// Answer implements the dns.Repository interface
func (u unimplemented) Answer(context.Context, *store.Record, *dnsr.Msg) {}

// Fallback implements the dns.Repository interface
func (u unimplemented) Fallback(context.Context, *store.Record, *dnsr.Msg) {}

// Unimplemented returns an unimplemented (and invalid) dns.Repository
func Unimplemented() unimplemented {
	return unimplemented{}
}
