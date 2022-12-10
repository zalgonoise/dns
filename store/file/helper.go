package file

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/x/spanner"
)

func toEntity(s *Store) []*store.Record {
	var out []*store.Record

	for _, recordType := range s.Types {
		rtype := recordType.RType
		for _, record := range recordType.Records {
			addr := record.Address
			for _, domain := range record.Domains {
				out = append(out, store.New().
					Name(domain).
					Type(rtype).
					Addr(addr).
					Build())
			}
		}
	}

	return out
}

func fromEntity(rs ...*store.Record) *Store {
	out := &Store{}

inputLoop:
	for _, r := range rs {
		for _, recordType := range out.Types {
			if recordType.RType == r.Type {
				for _, record := range recordType.Records {
					if record.Address == r.Addr {
						for _, domain := range record.Domains {
							if domain == r.Name {
								continue inputLoop
							}
						}
						record.Domains = append(record.Domains, r.Name)
						continue inputLoop
					}
				}
				recordType.Records = append(recordType.Records, &Record{
					Address: r.Addr,
					Domains: []string{r.Name},
				})
				continue inputLoop
			}
		}
		out.Types = append(out.Types, &Type{
			RType: r.Type,
			Records: []*Record{
				{
					Address: r.Addr,
					Domains: []string{r.Name},
				},
			},
		})
		continue inputLoop
	}

	return out
}

func (f *FileStore) sync(ctx context.Context) error {
	ctx, syncSpan := spanner.Start(ctx, "store.file.sync")
	defer syncSpan.End()

	// list records
	_, s := spanner.Start(ctx, "store.List")
	records, err := f.store.List(ctx)
	if err != nil {
		s.Event("failed to list store records", attr.String("error", err.Error()))
		s.End()
		return fmt.Errorf("%w: failed to list store records: %v", store.ErrSync, err)
	} else {
		s.Add(attr.Int("records_length", len(records)), attr.New("records", records))
	}
	s.End()

	// encode records
	_, s = spanner.Start(ctx, "store.file.Encode")
	b, err := f.enc.Encode(fromEntity(records...))
	if err != nil {
		s.Event("failed to encode store records", attr.String("error", err.Error()))
		s.End()
		return fmt.Errorf("%w: failed to marshal store records to JSON: %v", store.ErrSync, err)
	} else {
		s.Add(attr.Int("bytes_length", len(b)))
	}
	s.End()

	// write to file
	_, s = spanner.Start(ctx, "store.file.Write")
	err = os.WriteFile(f.Path, b, fs.FileMode(store.OS_ALL_RW))
	if err != nil {
		s.Event("failed to write store records to file", attr.String("error", err.Error()))
		s.End()
		return fmt.Errorf("%w: failed to write new reference file: %v", store.ErrSync, err)
	}
	s.End()
	return nil
}
