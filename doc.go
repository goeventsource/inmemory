// Package inmemory implements goeventsource Store, Repository, and Snapshotter using
// in-process maps and mutexes. It does not provide SQL, connections, or transaction APIs;
// those are specific to database-backed modules such as github.com/goeventsource/pgx.
package inmemory
