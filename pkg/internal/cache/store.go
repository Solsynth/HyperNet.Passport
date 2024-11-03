package cache

import (
	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/store"
	ristrettoCache "github.com/eko/gocache/store/ristretto/v4"
)

var S store.StoreInterface

func NewStore() error {
	ris, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 27,
		BufferItems: 64,
	})
	if err != nil {
		return err
	}

	S = ristrettoCache.NewRistretto(ris)

	return nil
}
