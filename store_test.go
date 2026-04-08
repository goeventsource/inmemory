package inmemory_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource/goeventsourcetest/goeventsourcetestintegration"

	"github.com/goeventsource/inmemory/inmemorytest"
)

func TestStore(t *testing.T) {
	cfg := inmemorytest.NewStoreConfig()
	s := inmemorytest.NewStore[uuid.UUID](cfg)
	goeventsourcetestintegration.TestStore(t, s)
}
