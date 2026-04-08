package inmemory_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
	"github.com/goeventsource/inmemory/inmemorytest"
	"github.com/goeventsource/goeventsource/goeventsourcetest/goeventsourcetestintegration"
)

func TestSnapshotter(t *testing.T) {
	t.Run("always", func(t *testing.T) {
		cfg := inmemorytest.NewSnapshotterConfig(goeventsource.SnapshotterWriteStrategyAlways[uuid.UUID, *goeventsourcetestintegration.User]())
		s := inmemorytest.NewSnapshotter(cfg)
		goeventsourcetestintegration.TestSnapshotterWithAlwaysStrategy(t, s)
	})

	t.Run("never", func(t *testing.T) {
		cfg := inmemorytest.NewSnapshotterConfig(goeventsource.SnapshotterWriteStrategyNever[uuid.UUID, *goeventsourcetestintegration.User]())
		s := inmemorytest.NewSnapshotter(cfg)
		goeventsourcetestintegration.TestSnapshotterWithNeverStrategy(t, s)
	})
}
