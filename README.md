# inmemory

**In-process implementation** of [goeventsource](https://github.com/goeventsource/goeventsource) storage interfaces: maps, mutexes, and no I/O. Use it for fast unit tests, spikes, and CI without standing up a database. It implements the same core contracts as other backends but does **not** mimic database transactions or connection pools—those belong in the [pgx](../pgx/README.md) module.

## Install

```bash
go get github.com/goeventsource/inmemory@latest
```

Requires `github.com/goeventsource/goeventsource` as a direct dependency.

## Why this module exists

The core library defines **what** event storage must do (`Append`, `Stream`, repository `Read` / `Write`, …). **inmemory** is one **how** for that contract: thread-safe in-process structures only. For durable PostgreSQL-backed storage, use [github.com/goeventsource/pgx](../pgx/README.md) instead.

## Packages


| Import                                           | Package name   | Purpose                                                           |
| ------------------------------------------------ | -------------- | ----------------------------------------------------------------- |
| `github.com/goeventsource/inmemory`              | `inmemory`     | `NewStore`, `NewRepository`, `NewSnapshotter`, repository options |
| `github.com/goeventsource/inmemory/inmemorytest` | `inmemorytest` | Shortcuts for tests: pre-wired configs and constructors           |


The root package name is `inmemory`, matching the module path (`inmemory.NewStore`, `inmemory.NewRepository`).

## Quick start: repository in tests

```go
package myapp_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/goeventsource/goeventsource"
	"github.com/goeventsource/inmemory/inmemorytest"
)

func TestWriteRead(t *testing.T) {
	ctx := context.Background()

	factory := func(id uuid.UUID, ver goeventsource.Version) *MyAggregate {
		return &MyAggregate{
			BaseRoot: goeventsource.NewBase(id, MyAggregateName, ver),
		}
	}

	cfg := inmemorytest.NewRepositoryConfig(factory)
	repo, _ := inmemorytest.NewRepository(cfg)

	agg := NewMyAggregate()
	if err := repo.Write(ctx, agg); err != nil {
		t.Fatal(err)
	}
	got, err := repo.Read(ctx, agg.ID())
	if err != nil {
		t.Fatal(err)
	}
	_ = got
}
```

Replace `MyAggregate` with your root type implementing `goeventsource.Root[uuid.UUID]`.

## Manual wiring (no `inmemorytest`)

When you need full control:

```go
import (
	"github.com/goeventsource/goeventsource"
	"github.com/goeventsource/inmemory"
)

store := inmemory.NewStore[goeventsource.ID]() // explicit K when there are no value args to infer from
repo := inmemory.NewRepository(store, factoryFunc,
	inmemory.WithSnapshotterOpt(mySnap),
	inmemory.WithProjectorsOpt(proj),
)
```

Use **goeventsource** codecs (`DomainEventEncodeDecoder`, `RootEncodeDecoder`, …) that match your domain; the in-memory store only stores encoded payloads and versions—it does not interpret SQL or wire protocols.

## Features worth knowing

- **Mutex-backed** `Store` suitable for parallel tests (lock scope per call, not distributed transactions).
- **Snapshotter on the repository path** via `WithSnapshotterOpt` for in-memory snapshot maps.
- `inmemorytest` provides the same kind of small config helpers for this module that `pgxtest` provides for PostgreSQL—each is scoped to its backend.

## Integration with core test suite

If you implement a custom store, run the shared scenarios from:

`github.com/goeventsource/goeventsource/goeventsourcetest/goeventsourcetestintegration`

The **inmemory** module’s own tests demonstrate that pattern.

## Tests

```bash
go test ./...
```

## Unpublished `goeventsource`

Until the core module is on the proxy, add a `replace` in `go.mod` or work inside [new_org/](../README.md) with `go.work`. Satellite CI often uses `go mod edit -replace=github.com/goeventsource/goeventsource@v0.0.0=<path>`.