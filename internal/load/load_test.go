package load

import (
	"testing"
	"time"

	"columnstore/internal/col"
	"columnstore/internal/meta"
	"columnstore/internal/part"
)

func TestLoadWritesBlocksAndPartitions(t *testing.T) {
	store := col.NewStore(t.TempDir())
	parts := part.NewRegistry()
	metaRegistry := meta.NewRegistry()
	loader := NewLoader(store, parts, metaRegistry)
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	src := NewSliceSource([]Batch{
		{
			Table:  "events",
			Column: "latency_ms",
			Rows:   []col.Value{col.Int(10), col.Int(20)},
			Start:  now,
			End:    now.Add(time.Hour),
			Seal:   true,
		},
	})
	written, err := loader.Load("events", src)
	if err != nil {
		t.Fatal(err)
	}
	if written != 1 {
		t.Fatalf("expected one batch, got %d", written)
	}
	if store.BlockCount() != 1 || parts.Count() != 1 {
		t.Fatalf("expected one block and one partition, got %d/%d", store.BlockCount(), parts.Count())
	}
}

func TestLoadRejectsInvalidBatch(t *testing.T) {
	store := col.NewStore(t.TempDir())
	parts := part.NewRegistry()
	metaRegistry := meta.NewRegistry()
	loader := NewLoader(store, parts, metaRegistry)
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	src := NewSliceSource([]Batch{
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
		t.Fatal("empty batch must be rejected")
	}
}

func TestBuildPlanSchema(t *testing.T) {
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	src := NewSliceSource([]Batch{
		{
			Table:  "events",
			Column: "latency_ms",
			Rows:   []col.Value{col.Int(1)},
			Start:  now,
			End:    now.Add(time.Hour),
			Seal:   true,
		},
	})
	plan := BuildPlan("events", src)
	if plan.Table != "events" || plan.Batches != 1 {
		t.Fatalf("unexpected plan: %+v", plan)
	}
	if len(plan.Schema.Columns) != 1 || plan.Schema.Columns[0].Name != "latency_ms" {
		t.Fatalf("unexpected schema: %+v", plan.Schema)
	}
}
