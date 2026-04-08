package inmemory_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/goeventsource/inmemory/inmemorytest"
	"github.com/goeventsource/goeventsource/goeventsourcetest/goeventsourcetestintegration"
)

func TestStore(t *testing.T) {
	cfg := inmemorytest.NewStoreConfig()
	s := inmemorytest.NewStore[uuid.UUID](cfg)
	goeventsourcetestintegration.TestStore(t, s)
}
