package file

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/attr"
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
	logx.From(ctx).Debug("syncing records", attr.String("action", "store:sync"))

	rs, err := f.store.List(context.Background())
	if err != nil {
		logx.From(ctx).Debug("failed to list records", attr.String("error", err.Error()))
		return fmt.Errorf("%w: failed to list store records: %v", store.ErrSync, err)
	}
	b, err := f.enc.Encode(fromEntity(rs...))
	if err != nil {
		logx.From(ctx).Debug("failed to encode records", attr.String("error", err.Error()))
		return fmt.Errorf("%w: failed to marshal store records to JSON: %v", store.ErrSync, err)
	}
	err = os.WriteFile(f.Path, b, fs.FileMode(store.OS_ALL_RW))
	if err != nil {
		logx.From(ctx).Debug("failed to write records to file", attr.String("error", err.Error()))
		return fmt.Errorf("%w: failed to write new reference file: %v", store.ErrSync, err)
	}

	logx.From(ctx).Debug("synced records successfully")
	return nil
}
