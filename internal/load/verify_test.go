package load_test

import (
	"testing"
	"time"

	"columnstore/internal/col"
	"columnstore/internal/load"
	"columnstore/internal/meta"
	"columnstore/internal/part"
)

func TestLoadRollbackCoversAllPartitions(t *testing.T) {
	store := col.NewStore(t.TempDir())
	parts := part.NewRegistry()
	metaRegistry := meta.NewRegistry()
	loader := load.NewLoader(store, parts, metaRegistry)
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	src := load.NewSliceSource([]load.Batch{
		{
			Table:  "events",
			Column: "latency_ms",
			Rows:   []col.Value{col.Int(7), col.Int(8)},
			Start:  now,
			End:    now.Add(time.Hour),
			Seal:   false,
		},
		{
			Table:  "events",
			Column: "latency_ms",
			Rows:   []col.Value{},
			Start:  now,
			End:    now.Add(time.Hour),
			Seal:   true,
		},
	})
	if _, err := loader.Load("events", src); err == nil {
		t.Fatal("load must fail on the invalid batch")
	}
	if store.BlockCount() != 0 {
		t.Fatalf("rollback must remove in-flight blocks, got %d blocks", store.BlockCount())
	}
	if parts.Count() != 0 {
		t.Fatalf("rollback must remove partitions, got %d", parts.Count())
	}
}
