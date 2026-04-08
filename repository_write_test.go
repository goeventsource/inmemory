package inmemory_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
	"github.com/goeventsource/goeventsource/goeventsourcetest/goeventsourcetestintegration"
	"github.com/goeventsource/inmemory"
)

type appendFailStore[K goeventsource.ID] struct {
	inner goeventsource.Store[K]
}

func (a *appendFailStore[K]) Append(ctx context.Context, evs ...goeventsource.Event) error {
	return errors.New("append failed")
}

func (a *appendFailStore[K]) Stream(ctx context.Context, id K, f goeventsource.StoreStreamFilter) (goeventsource.Events, error) {
	return a.inner.Stream(ctx, id, f)
}

func TestRepository_Write_appendError_keepsPeekedEvents(t *testing.T) {
	ctx := context.Background()
	mem := inmemory.NewStore[uuid.UUID]()
	store := &appendFailStore[uuid.UUID]{inner: mem}
	factory := func(id uuid.UUID, version goeventsource.Version) *goeventsourcetestintegration.User {
		return &goeventsourcetestintegration.User{
			BaseRoot: goeventsource.NewBase(id, goeventsourcetestintegration.UserAggregateName, version),
		}
	}
	r := inmemory.NewRepository(store, factory)

	u := goeventsourcetestintegration.Register("x", "y@z")
	if err := r.Write(ctx, u); err == nil {
		t.Fatal("expected error")
	}
	if len(goeventsource.PeekEvents(u)) != 1 {
		t.Fatalf("pending events should remain after failed write, got %d", len(goeventsource.PeekEvents(u)))
	}
}
