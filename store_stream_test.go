package inmemory

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
)

type stubEv string

func (stubEv) DomainEventName() goeventsource.DomainEventName { return "stub" }

func TestStore_Stream_versionFilter_matchesPgxSemantics(t *testing.T) {
	ctx := context.Background()
	s := NewStore[uuid.UUID]()
	id := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	for _, v := range []goeventsource.Version{1, 2, 3} {
		ev := goeventsource.Event{
			DomainEvent:     stubEv("x"),
			DomainEventName: "stub",
			Version:         v,
			StreamID:        id,
			StreamName:      "agg",
			OccurredAt:      time.Now(),
		}
		if err := s.Append(ctx, ev); err != nil {
			t.Fatal(err)
		}
	}

	evs, err := s.Stream(ctx, id, goeventsource.StoreStreamFilter{From: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(evs) != 1 || evs[0].Version != 3 {
		t.Fatalf("want single event v3, got %#v", evs)
	}

	_, err = s.Stream(ctx, id, goeventsource.StoreStreamFilter{From: 4})
	if !errors.Is(err, goeventsource.ErrStoreStreamEmpty) {
		t.Fatalf("want empty for From past tail, got %v", err)
	}
}
