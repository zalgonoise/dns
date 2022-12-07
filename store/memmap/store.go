package memmap

import (
	"context"

	"github.com/zalgonoise/dns/store"
)

// Create implements the store.Repository interface
//
// It will not perform any lookups before writing the new records, and it will simply
// blindly write the input records to the store.
func (m *MemoryStore) Create(ctx context.Context, rs ...*store.Record) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, r := range rs {
		if _, ok := m.Records[r.Type]; !ok {
			m.Records[r.Type] = map[string]string{}
		}
		m.Records[r.Type][r.Name] = r.Addr
	}
	return nil
}

// List implements the store.Repository interface
//
// It will build a list of pointers to store.Record which is returned alongside
// any errors that are raised (currently there are no scenarios in this implementation)
func (m *MemoryStore) List(ctx context.Context) ([]*store.Record, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	var output []*store.Record

	for rtype, r := range m.Records {
		for domain, addr := range r {
			output = append(
				output,
				store.New().
					Type(rtype).
					Name(domain).
					Addr(addr).
					Build(),
			)
		}
	}
	return output, nil
}

// FindByTypeAndDomain implements the store.Repository interface
//
// It will return a pointer to a store.Record if there is an IP address
// registered to the input store.Record's domain name and record type.
//
// It also returns an error in case the record does not exist
func (m *MemoryStore) FindByTypeAndDomain(ctx context.Context, rtype, domain string) (*store.Record, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, ok := m.Records[rtype]; !ok {
		return nil, store.ErrDoesNotExist
	}
	dest := m.Records[rtype][domain]
	if dest == "" {
		return nil, store.ErrDoesNotExist
	}

	return store.New().Type(rtype).Name(domain).Addr(dest).Build(), nil
}

// FilterByDomain implements the store.Repository interface
//
// It will return a list of pointers to store.Record if there records associated
// with the input domain name, for all record types.
//
// It also returns an error in case the record does not exist
func (m *MemoryStore) FilterByDomain(ctx context.Context, domain string) ([]*store.Record, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	var out = []*store.Record{}
	for rtype, domains := range m.Records {
		for name, addr := range domains {
			if name == domain {
				out = append(
					out,
					store.New().Type(rtype).Name(name).Addr(addr).Build(),
				)
			}
		}
	}

	return out, nil
}

// FilterByDest implements the store.Repository interface
//
// It will return a slice of pointers to store.Records if there are records
// associated with the input store.Record's IP address.
//
// If the call is successful but there are no records associated to that addres,
// returns an empty slice.
//
// It also returns an error in case the operation fails (which is currently not
// a scenario)
func (m *MemoryStore) FilterByDest(ctx context.Context, addr string) ([]*store.Record, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	var output []*store.Record

	for rtype, domains := range m.Records {
		for domain, ipAddr := range domains {
			if ipAddr == addr {
				output = append(
					output,
					store.New().
						Type(rtype).
						Name(domain).
						Addr(addr).
						Build(),
				)
			}
		}
	}

	return output, nil
}

// Update implements the store.Repository interface
//
// It will target a particular domain name, and update its target IP address
// based on the input store.Record.
//
// If it targets a domain which does not exist in the store, or if that domain
// does not have that record type registered, it returns a DoesNotExist error
//
// If the operation is successful, it returns nil
func (m *MemoryStore) Update(ctx context.Context, domain string, r *store.Record) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, ok := m.Records[r.Type]; !ok {
		return store.ErrDoesNotExist
	}
	if _, ok := m.Records[r.Type][domain]; !ok {
		return store.ErrDoesNotExist
	}
	if domain != r.Name {
		delete(m.Records[r.Type], domain)
	}
	m.Records[r.Type][r.Name] = r.Addr

	return nil
}

// DeleteByAddress removes all records with IP address `addr`
func (m *MemoryStore) DeleteByAddress(ctx context.Context, addr string) error {
	for rtype, rmap := range m.Records {
		for domain, address := range rmap {
			if address == addr {
				delete(m.Records[rtype], domain)
			}
		}
	}
	return nil
}

// DeleteByDomain removes all records with domain name `name`
func (m *MemoryStore) DeleteByDomain(ctx context.Context, name string) error {
	for rtype, domains := range m.Records {
		for domain := range domains {
			if domain == name {
				delete(m.Records[rtype], domain)
			}
		}
	}
	return nil
}

// DeleteByTypeAndDomain removes all records with record type `rtype` and domain name `name`
func (m *MemoryStore) DeleteByTypeAndDomain(ctx context.Context, rtype, name string) error {
	for recordtype, domains := range m.Records {
		if recordtype == rtype {
			for domain := range domains {
				if domain == name {
					delete(m.Records[recordtype], domain)
				}
			}
		}
	}
	return nil
}
