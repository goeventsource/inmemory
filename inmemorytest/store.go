package inmemorytest

import (
	"github.com/goeventsource/goeventsource"
	"github.com/goeventsource/inmemory"
)

// StoreConfig represents the configuration to create an inmemory.Store via NewStore
type StoreConfig struct {
	StoreAppendOpts []goeventsource.StoreAppendOpt
}

// NewStoreConfig creates a new StoreConfig with default values for testing purposes.
func NewStoreConfig() StoreConfig {
	return StoreConfig{}
}

// NewStore creates a new instance of an inmemory.Store based on the provided StoreConfig.
func NewStore[K goeventsource.ID](cfg StoreConfig) *inmemory.Store[K] {
	return inmemory.NewStore[K](cfg.StoreAppendOpts...)
}
