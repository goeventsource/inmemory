package inmemorytest_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
	"github.com/goeventsource/goeventsource/goeventsourcetest/goeventsourcetestintegration"

	. "github.com/goeventsource/inmemory/inmemorytest"
)

func TestNewSnapshotterConfig(t *testing.T) {
	strategy := goeventsource.SnapshotterWriteStrategyAlways[uuid.UUID, *goeventsourcetestintegration.User]()
	cfg := NewSnapshotterConfig(strategy)

	if reflect.ValueOf(cfg.SnapshotterWriteStrategy).Pointer() != reflect.ValueOf(strategy).Pointer() {
		t.Fatalf("snapshotter storage was not the given one: %v", cfg.SnapshotterWriteStrategy)
	}
}

func TestNewSnapshotter(t *testing.T) {
	strategy := goeventsource.SnapshotterWriteStrategyAlways[uuid.UUID, *goeventsourcetestintegration.User]()
	cfg := SnapshotterConfig[uuid.UUID, *goeventsourcetestintegration.User]{
		SnapshotterWriteStrategy: strategy,
	}

	snap := NewSnapshotter(cfg)

	if reflect.ValueOf(snap.WriteStrategy).Pointer() != reflect.ValueOf(strategy).Pointer() {
		t.Fatalf("snapshotter write strategy was not the given one: %v", snap.WriteStrategy)
	}

	_, err := snap.ReadSnapshot(context.Background(), uuid.New())
	if !errors.Is(err, goeventsource.ErrSnapshotterReadNotFound) {
		t.Fatalf("expected empty snapshotter: %v", err)
	}
}
