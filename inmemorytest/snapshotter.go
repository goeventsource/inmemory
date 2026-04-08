package inmemorytest

import (
	"github.com/goeventsource/goeventsource"

	"github.com/goeventsource/inmemory"
)

// SnapshotterConfig represents the configuration to create an inmemory.Snapshotter via NewSnapshotter
type SnapshotterConfig[K goeventsource.ID, V goeventsource.Root[K]] struct {
	SnapshotterWriteStrategy goeventsource.SnapshotterWriteStrategy[K, V]
}

// NewSnapshotterConfig creates a new SnapshotterConfig with default values for testing purposes.
func NewSnapshotterConfig[K goeventsource.ID, V goeventsource.Root[K]](
	strategy goeventsource.SnapshotterWriteStrategy[K, V],
) SnapshotterConfig[K, V] {
	return SnapshotterConfig[K, V]{
		SnapshotterWriteStrategy: strategy,
	}
}

// NewSnapshotter creates a new instance of an inmemory.Snapshotter based on the provided SnapshotterConfig.
func NewSnapshotter[K goeventsource.ID, V goeventsource.Root[K]](cfg SnapshotterConfig[K, V]) *inmemory.Snapshotter[K, V] {
	return inmemory.NewSnapshotter(cfg.SnapshotterWriteStrategy)
}
