package inmemorytest_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
	"github.com/goeventsource/inmemory"
	. "github.com/goeventsource/inmemory/inmemorytest"
	"github.com/goeventsource/goeventsource/goeventsourcetest/goeventsourcetestintegration"
)

func TestNewRepositoryConfig(t *testing.T) {
	factoryFn := goeventsource.FactoryFunc[uuid.UUID, *goeventsourcetestintegration.User](func(id uuid.UUID, version goeventsource.Version) *goeventsourcetestintegration.User {
		return &goeventsourcetestintegration.User{BaseRoot: goeventsource.NewBase(id, goeventsourcetestintegration.UserAggregateName, version)}
	})

	opts := []inmemory.RepositoryOpt[uuid.UUID, *goeventsourcetestintegration.User]{
		inmemory.WithProjectorsOpt[uuid.UUID, *goeventsourcetestintegration.User](nil),
		inmemory.WithSnapshotterOpt[uuid.UUID, *goeventsourcetestintegration.User](nil),
	}
	cfg := NewRepositoryConfig(factoryFn, opts...)

	if !reflect.DeepEqual(cfg.StoreConfig, StoreConfig{}) {
		t.Fatalf("repository store config was not the expected one: %v", cfg.StoreConfig)
	}

	if reflect.ValueOf(cfg.FactoryFunc).Pointer() != reflect.ValueOf(factoryFn).Pointer() {
		t.Fatalf("repository factory func was not the given one: %v", cfg.FactoryFunc)
	}

	if !reflect.DeepEqual(cfg.Opts, opts) {
		t.Fatalf("repository options were not the given one: %v", cfg.Opts)
	}
}

func TestNewRepository(t *testing.T) {
	factoryFn := goeventsource.FactoryFunc[uuid.UUID, *goeventsourcetestintegration.User](func(id uuid.UUID, version goeventsource.Version) *goeventsourcetestintegration.User {
		return &goeventsourcetestintegration.User{BaseRoot: goeventsource.NewBase(id, goeventsourcetestintegration.UserAggregateName, version)}
	})

	projectors := []goeventsource.Projector{}
	snapshotter := NewSnapshotter(NewSnapshotterConfig[uuid.UUID, *goeventsourcetestintegration.User](nil))

	opts := []inmemory.RepositoryOpt[uuid.UUID, *goeventsourcetestintegration.User]{
		inmemory.WithProjectorsOpt[uuid.UUID, *goeventsourcetestintegration.User](projectors...),
		inmemory.WithSnapshotterOpt[uuid.UUID, *goeventsourcetestintegration.User](snapshotter),
	}

	cfg := NewRepositoryConfig(factoryFn, opts...)
	repo, store := NewRepository(cfg)

	if repo == nil {
		t.Fatal("repository should not be nil")
	}

	if store == nil {
		t.Fatal("store should not be nil")
	}

	if !reflect.DeepEqual(store, NewStore[uuid.UUID](NewStoreConfig())) {
		t.Fatalf("store was not the default one: %v", store)
	}
}
