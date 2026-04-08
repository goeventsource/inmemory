package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/goeventsource/goeventsource"
)

var (
	// Ensure Store implements goeventsource.Store at compile time
	_ goeventsource.Store[goeventsource.ID] = &Store[goeventsource.ID]{}
)

// Store is an in-memory implementation of the goeventsource.Store interface.
type Store[K goeventsource.ID] struct {
	storage        map[string][]goeventsource.Event
	versionStorage map[string]struct{}
	appendOpts     []goeventsource.StoreAppendOpt

	mu sync.RWMutex
}

// NewStore creates a new instance of the Store.
func NewStore[K goeventsource.ID](opts ...goeventsource.StoreAppendOpt) *Store[K] {
	return &Store[K]{
		storage:        map[string][]goeventsource.Event{},
		versionStorage: map[string]struct{}{},
		appendOpts:     opts,
	}
}

// Append validates and stores the batch under one write lock: on version conflict, none of the events are kept.
// This is mutex-scoped consistency, not a SQL/database transaction.
func (s *Store[K]) Append(ctx context.Context, evs ...goeventsource.Event) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %w", goeventsource.ErrStoreAppend, ctx.Err())
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	prepared := make([]goeventsource.Event, len(evs))
	for i := range evs {
		ev := evs[i]
		for j := range s.appendOpts {
			ev.Metadata = s.appendOpts[j](ctx, ev.Metadata)
		}
		prepared[i] = ev
	}

	seen := make(map[string]struct{}, len(prepared))
	for i := range prepared {
		versionKey := fmt.Sprintf("%s_%d", prepared[i].StreamID.String(), prepared[i].Version)
		if _, ok := s.versionStorage[versionKey]; ok {
			return goeventsource.ErrStoreAppendVersionConflict
		}
		if _, ok := seen[versionKey]; ok {
			return goeventsource.ErrStoreAppendVersionConflict
		}
		seen[versionKey] = struct{}{}
	}

	for i := range prepared {
		ev := prepared[i]
		id := ev.StreamID.String()
		if _, ok := s.storage[id]; !ok {
			s.storage[id] = []goeventsource.Event{}
		}

		versionKey := fmt.Sprintf("%s_%d", id, ev.Version)
		s.storage[id] = append(s.storage[id], ev)
		s.versionStorage[versionKey] = struct{}{}
	}

	return nil
}

// Stream returns a defensive copy of events for id with Version >= filter.From (inclusive).
// From == 0 means no lower bound on version. If there is no stream or no matching events, returns ErrStoreStreamEmpty.
func (s *Store[K]) Stream(ctx context.Context, id K, f goeventsource.StoreStreamFilter) (goeventsource.Events, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%w: %w", goeventsource.ErrStoreStream, ctx.Err())
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	evs, ok := s.storage[id.String()]
	if !ok {
		return nil, goeventsource.ErrStoreStreamEmpty
	}

	var out goeventsource.Events
	for _, ev := range evs {
		if f.From != 0 && ev.Version < f.From {
			continue
		}
		cp := ev
		if ev.Metadata != nil {
			cp.Metadata = cloneMetadata(ev.Metadata)
		}
		out = append(out, cp)
	}
	if len(out) == 0 {
		return nil, goeventsource.ErrStoreStreamEmpty
	}
	return out, nil
}

func cloneMetadata(md goeventsource.Metadata) goeventsource.Metadata {
	if md == nil {
		return nil
	}
	c := make(goeventsource.Metadata, len(md))
	for k, v := range md {
		c[k] = v
	}
	return c
}
