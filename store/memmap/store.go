package memmap

import (
	"context"

	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/attr"
)

// Create implements the store.Repository interface
//
// It will not perform any lookups before writing the new records, and it will simply
// blindly write the input records to the store.
func (m *MemoryStore) Create(ctx context.Context, rs ...*store.Record) error {
	logx.From(ctx).Debug("creating record", attr.String("action", "store:create"))

	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, r := range rs {
		if _, ok := m.Records[r.Type]; !ok {
			m.Records[r.Type] = map[string]string{}
		}
		m.Records[r.Type][r.Name] = r.Addr
	}
	logx.From(ctx).Debug("added record(s) successfully")
	return nil
}

// List implements the store.Repository interface
//
// It will build a list of pointers to store.Record which is returned alongside
// any errors that are raised (currently there are no scenarios in this implementation)
func (m *MemoryStore) List(ctx context.Context) ([]*store.Record, error) {
	logx.From(ctx).Debug("listing records", attr.String("action", "store:list"))

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
	logx.From(ctx).Debug("listed records successfully")
	return output, nil
}

// FilterByTypeAndDomain implements the store.Repository interface
//
// It will return a pointer to a store.Record if there is an IP address
// registered to the input store.Record's domain name and record type.
//
// It also returns an error in case the record does not exist
func (m *MemoryStore) FilterByTypeAndDomain(ctx context.Context, rtype, domain string) (*store.Record, error) {
	logx.From(ctx).Debug("reading records", attr.String("action", "store:read"))

	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, ok := m.Records[rtype]; !ok {
		logx.From(ctx).Error("failed to find record by type", attr.String("error", store.ErrDoesNotExist.Error()))
		return nil, store.ErrDoesNotExist
	}
	dest := m.Records[rtype][domain]
	if dest == "" {
		logx.From(ctx).Error("failed to find record by domain", attr.String("error", store.ErrDoesNotExist.Error()))
		return nil, store.ErrDoesNotExist
	}

	logx.From(ctx).Debug("read records successfully")
	return store.New().Type(rtype).Name(domain).Addr(dest).Build(), nil
}

// FilterByDomain implements the store.Repository interface
//
// It will return a list of pointers to store.Record if there records associated
// with the input domain name, for all record types.
//
// It also returns an error in case the record does not exist
func (m *MemoryStore) FilterByDomain(ctx context.Context, domain string) ([]*store.Record, error) {
	logx.From(ctx).Debug("listing records", attr.String("action", "store:list"))

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

	logx.From(ctx).Debug("listed records successfully")
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
	logx.From(ctx).Debug("reading records", attr.String("action", "store:read"))

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

	logx.From(ctx).Debug("read records successfully")
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
	logx.From(ctx).Debug("reading records", attr.String("action", "store:update"))

	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, ok := m.Records[r.Type]; !ok {
		logx.From(ctx).Error("failed to find record by type", attr.String("error", store.ErrDoesNotExist.Error()))
		return store.ErrDoesNotExist
	}
	if _, ok := m.Records[r.Type][domain]; !ok {
		logx.From(ctx).Error("failed to find record by domain", attr.String("error", store.ErrDoesNotExist.Error()))
		return store.ErrDoesNotExist
	}
	if domain != r.Name {
		logx.From(ctx).Info("different domain name -- removing former entry",
			attr.String("type", r.Type),
			attr.String("domain", domain),
			attr.String("addr_to_delete", m.Records[r.Type][domain]),
		)
		delete(m.Records[r.Type], domain)
	}
	m.Records[r.Type][r.Name] = r.Addr

	logx.From(ctx).Debug("updated records successfully")
	return nil
}

// Delete implements the store.Repository interface
//
// It will leverage the input store.Record to commit to a record deletion.
//
// If:
//   - a domain name is provided: its target IP address and record types are removed
//   - a domain name and record type are provided: remove the target IP address associated
//   - IP address is populated: delete all records associated with that address
//
// It returns an error if the operation is unsuccessful (which is not a scenario as of now)
func (m *MemoryStore) Delete(ctx context.Context, r *store.Record) error {
	logx.From(ctx).Debug("deleting records", attr.String("action", "store:delete"))

	m.mtx.Lock()
	defer m.mtx.Unlock()

	if r.Name != "" && r.Type == "" && r.Addr == "" {
		logx.From(ctx).Debug("deleting all entries with domain filter", attr.String("domain", r.Name))
		deleteDomain(m, r.Name)
	}
	if r.Name != "" && r.Type != "" && r.Addr == "" {
		logx.From(ctx).Debug("deleting all entries with domain and type filters", attr.String("domain", r.Name), attr.String("type", r.Type))
		deleteDomainByType(m, r.Name, r.Type)
	}
	if r.Addr != "" {
		logx.From(ctx).Debug("deleting all entries with address filter", attr.String("address", r.Addr))
		deleteAddress(m, r.Addr)
	}

	logx.From(ctx).Debug("deleted records successfully")
	return nil
}
