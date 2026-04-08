package inmemorytest

import (
	"github.com/goeventsource/goeventsource"

	"github.com/goeventsource/inmemory"
)

// RepositoryConfig represents the configuration to create an inmemory.Repository via NewRepository
type RepositoryConfig[K goeventsource.ID, V goeventsource.Root[K]] struct {
	StoreConfig
	FactoryFunc goeventsource.FactoryFunc[K, V]
	Opts        []inmemory.RepositoryOpt[K, V]
}

// NewRepositoryConfig creates a new RepositoryConfig with default values for testing purposes.
func NewRepositoryConfig[K goeventsource.ID, V goeventsource.Root[K]](
	factoryFunc goeventsource.FactoryFunc[K, V],
	opts ...inmemory.RepositoryOpt[K, V],
) RepositoryConfig[K, V] {
	return RepositoryConfig[K, V]{
		StoreConfig: NewStoreConfig(),
		FactoryFunc: factoryFunc,
		Opts:        opts,
	}
}

// NewRepository creates a new instance of an inmemory.Repository based on the provided StoreConfig.
func NewRepository[K goeventsource.ID, V goeventsource.Root[K]](
	cfg RepositoryConfig[K, V],
) (*inmemory.Repository[K, V], *inmemory.Store[K]) {
	s := inmemory.NewStore[K](cfg.StoreAppendOpts...)
	return inmemory.NewRepository(s, cfg.FactoryFunc, cfg.Opts...), s
}
