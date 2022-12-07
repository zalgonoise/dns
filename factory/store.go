package factory

import (
	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/dns/store/file"
	"github.com/zalgonoise/dns/store/memmap"
)

func StoreRepository(rtype string, path string) store.Repository {
	var storeRepo store.Repository

	switch rtype {
	case "memmap", "memory", "in-memory":
		storeRepo = memmap.New()
	case "jsonfile", "json":
		storeRepo = file.New("json", path)
	case "yamlfile", "yaml":
		storeRepo = file.New("yaml", path)
	default:
		storeRepo = memmap.New()
	}

	return store.WithTrace(storeRepo)
}
