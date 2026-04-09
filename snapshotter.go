package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/goeventsource/goeventsource"
)

// Snapshotter holds aggregate snapshots in a process-local map guarded by a mutex.
// It stores the root value directly (often a pointer): mutating that aggregate after WriteSnapshot corrupts the snapshot.
// ReadSnapshot returns the same instance each time—no isolation like persistent snapshotters that serialize/deserialize.
// This is intentional for an in-memory implementation intended for testing and development, not production.
type Snapshotter[K goeventsource.ID, V goeventsource.Root[K]] struct {
	snapshots     map[string]V
	WriteStrategy goeventsource.SnapshotterWriteStrategy[K, V]

	mu sync.RWMutex
}

// NewSnapshotter returns an instance of a Snapshotter
func NewSnapshotter[K goeventsource.ID, V goeventsource.Root[K]](strategy goeventsource.SnapshotterWriteStrategy[K, V]) *Snapshotter[K, V] {
	return &Snapshotter[K, V]{
		snapshots:     map[string]V{},
		WriteStrategy: strategy,
	}
}

// WriteSnapshot stores root under its ID when the write strategy allows it.
func (s *Snapshotter[K, V]) WriteSnapshot(ctx context.Context, root V) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %w", goeventsource.ErrSnapshotterWrite, ctx.Err())
	default:
	}

	if !s.WriteStrategy(root) {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.snapshots[goeventsource.RootID(root).String()] = root
	return nil
}

// ReadSnapshot returns the stored root for id, or a not-found error from the core module.
func (s *Snapshotter[K, V]) ReadSnapshot(ctx context.Context, k K) (V, error) {
	var zero V
	select {
	case <-ctx.Done():
		return zero, fmt.Errorf("%w: %w", goeventsource.ErrSnapshotterRead, ctx.Err())
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	root, ok := s.snapshots[k.String()]
	if !ok {
		return zero, goeventsource.ErrSnapshotterReadNotFound
	}

	return root, nil
}
