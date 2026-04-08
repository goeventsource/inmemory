package inmemorytest_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
	. "github.com/goeventsource/inmemory/inmemorytest"
)

func TestNewStoreConfig(t *testing.T) {
	cfg := NewStoreConfig()
	if len(cfg.StoreAppendOpts) != 0 {
		t.Fatalf("store append options were not empty")
	}
}

func TestNewStore(t *testing.T) {
	opts := []goeventsource.StoreAppendOpt{
		func(ctx context.Context, md goeventsource.Metadata) goeventsource.Metadata {
			if md == nil {
				md = goeventsource.Metadata{}
			}
			md["k"] = "v"
			return md
		},
	}

	cfg := StoreConfig{StoreAppendOpts: opts}
	s := NewStore[uuid.UUID](cfg)

	ctx := context.Background()
	_, err := s.Stream(ctx, uuid.New(), goeventsource.StoreStreamNoFilter())
	if !errors.Is(err, goeventsource.ErrStoreStreamEmpty) {
		t.Fatalf("expected empty stream for unknown id: %v", err)
	}

	streamID := uuid.New()
	ev := goeventsource.Event{
		DomainEventName: "test_ev",
		Version:         1,
		StreamID:        streamID,
		StreamName:      "s",
		OccurredAt:      time.Now(),
	}
	if err := s.Append(ctx, ev); err != nil {
		t.Fatalf("append: %v", err)
	}
	evs, err := s.Stream(ctx, streamID, goeventsource.StoreStreamNoFilter())
	if err != nil {
		t.Fatalf("stream: %v", err)
	}
	if len(evs) != 1 || evs[0].Metadata["k"] != "v" {
		t.Fatalf("append opts not applied: %#v", evs)
	}
}
