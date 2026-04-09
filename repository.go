package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/goeventsource/goeventsource"
)

type (
	repoID   = goeventsource.ID
	repoRoot = goeventsource.Root[repoID]
)

var (
	// Ensure Repository implements goeventsource.Repository at compile time.
	_ goeventsource.Repository[repoID, repoRoot] = (*Repository[repoID, repoRoot])(nil)
)

// RepositoryOpt is a function signature for providing options to configure a Repository.
type RepositoryOpt[K goeventsource.ID, V goeventsource.Root[K]] func(*Repository[K, V])

// WithProjectorsOpt is a RepositoryOpt that sets a slice of goeventsource.Projector for a Repository.
func WithProjectorsOpt[K goeventsource.ID, V goeventsource.Root[K]](ps ...goeventsource.Projector) RepositoryOpt[K, V] {
	return func(r *Repository[K, V]) {
		r.projectors = ps
	}
}

// WithSnapshotterOpt is a RepositoryOpt that sets a Snapshotter for the Repository.
func WithSnapshotterOpt[K goeventsource.ID, V goeventsource.Root[K]](s goeventsource.Snapshotter[K, V]) RepositoryOpt[K, V] {
	return func(r *Repository[K, V]) {
		r.snapshotter = s
	}
}

// Repository is an in-memory goeventsource.Repository implementation.
type Repository[K goeventsource.ID, V goeventsource.Root[K]] struct {
	store       goeventsource.Store[K]
	factoryFunc goeventsource.FactoryFunc[K, V]
	projectors  []goeventsource.Projector
	snapshotter goeventsource.Snapshotter[K, V]
}

// NewRepository creates a new instance of Repository.
func NewRepository[K goeventsource.ID, V goeventsource.Root[K]](
	store goeventsource.Store[K],
	factoryFunc goeventsource.FactoryFunc[K, V],
	opts ...RepositoryOpt[K, V],
) *Repository[K, V] {
	r := &Repository[K, V]{
		store:       store,
		factoryFunc: factoryFunc,
	}

	for i := range opts {
		opts[i](r)
	}

	return r
}

// Read reads the goeventsource.Events from a goeventsource.Store and rebuild the goeventsource.Root state for the given goeventsource.ID.
// It returns the root aggregate and an error if the aggregate rootID is not found or an error occurs.
func (r Repository[K, V]) Read(ctx context.Context, id K) (V, error) {
	var (
		zero        V
		hadSnapshot bool
		root        = r.factoryFunc(id, 0)
		filter      = goeventsource.StoreStreamNoFilter()
	)

	select {
	case <-ctx.Done():
		return zero, fmt.Errorf("%w: %w", goeventsource.ErrRepositoryRead, ctx.Err())
	default:
	}

	if r.snapshotter != nil {
		snap, err := r.snapshotter.ReadSnapshot(ctx, id)
		switch {
		case errors.Is(err, goeventsource.ErrSnapshotterReadNotFound):
			// ignore
		case err != nil:
			return zero, fmt.Errorf("%w: %w", goeventsource.ErrRepositoryRead, err)
		default:
			hadSnapshot = true
			root = snap
			filter.From = goeventsource.RootVersion(root) + 1
		}
	}

	evs, err := r.store.Stream(ctx, id, filter)
	switch {
	case errors.Is(err, goeventsource.ErrStoreStreamEmpty) && hadSnapshot:
		return root, nil
	case errors.Is(err, goeventsource.ErrStoreStreamEmpty):
		return zero, fmt.Errorf("%w: %w", goeventsource.ErrRepositoryReadNotFound, err)
	case err != nil:
		return zero, fmt.Errorf("%w: %w", goeventsource.ErrRepositoryRead, err)
	default:
	}

	goeventsource.PushEvents(root, evs)

	return root, nil
}

// Write appends pending events from root to the store and runs projectors / snapshot side effects.
// It returns an error if an error occurs during the write operation.
func (r Repository[K, V]) Write(ctx context.Context, root V) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %w", goeventsource.ErrRepositoryWrite, ctx.Err())
	default:
	}

	evs := goeventsource.PeekEvents(root)
	if err := r.store.Append(ctx, evs...); err != nil {
		return fmt.Errorf("%w: %w", goeventsource.ErrRepositoryWrite, err)
	}

	for i := range r.projectors {
		if err := r.projectors[i].Project(ctx, evs...); err != nil {
			return fmt.Errorf("%w: %w", goeventsource.ErrRepositoryWrite, err)
		}
	}

	if r.snapshotter != nil {
		if err := r.snapshotter.WriteSnapshot(ctx, root); err != nil {
			return fmt.Errorf("%w: %w", goeventsource.ErrRepositoryWrite, err)
		}
	}

	goeventsource.FlushEvents(root)
	return nil
}
