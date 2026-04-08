package inmemory_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
	"github.com/goeventsource/goeventsource/goeventsourcetest/goeventsourcetestintegration"

	"github.com/goeventsource/inmemory"
	"github.com/goeventsource/inmemory/inmemorytest"
)

func TestRepository(t *testing.T) {
	factoryFunc := func(id uuid.UUID, version goeventsource.Version) *goeventsourcetestintegration.User {
		return &goeventsourcetestintegration.User{BaseRoot: goeventsource.NewBase(id, goeventsourcetestintegration.UserAggregateName, version)}
	}

	t.Run("repository", func(t *testing.T) {
		cfg := inmemorytest.NewRepositoryConfig(factoryFunc)
		r, s := inmemorytest.NewRepository(cfg)
		goeventsourcetestintegration.TestRepository(t, r, s)
	})

	t.Run("repository_with_snapshots", func(t *testing.T) {
		snapStrategy := goeventsource.SnapshotterWriteStrategyAlways[uuid.UUID, *goeventsourcetestintegration.User]()
		snapCfg := inmemorytest.NewSnapshotterConfig(snapStrategy)
		snap := inmemorytest.NewSnapshotter(snapCfg)

		cfg := inmemorytest.NewRepositoryConfig(
			factoryFunc,
			inmemory.WithSnapshotterOpt(snap),
		)

		r, s := inmemorytest.NewRepository(cfg)
		goeventsourcetestintegration.TestRepository(t, r, s)
	})

	t.Run("repository_with_projector", func(t *testing.T) {
		t.Run("repository_with_projector", func(t *testing.T) {
			proj := &goeventsourcetestintegration.Projector{}
			cfg := inmemorytest.NewRepositoryConfig(
				factoryFunc,
				inmemory.WithProjectorsOpt[uuid.UUID, *goeventsourcetestintegration.User](proj),
			)
			r, s := inmemorytest.NewRepository(cfg)
			goeventsourcetestintegration.TestRepositoryWithProjector(t, r, proj, s)
		})
	})
}
