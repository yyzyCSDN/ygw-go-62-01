package query

import (
	"testing"
	"time"

	"columnstore/internal/col"
	"columnstore/internal/compress"
	"columnstore/internal/load"
	"columnstore/internal/meta"
	"columnstore/internal/part"
)

func TestExecuteReturnsAggregate(t *testing.T) {
	store := col.NewStore(t.TempDir())
	parts := part.NewRegistry()
	metaRegistry := meta.NewRegistry()
	dicts := compress.NewDictRegistry()
	loader := load.NewLoader(store, parts, metaRegistry)
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	_, err := loader.Load("events", load.NewSliceSource([]load.Batch{
		{
			Table:  "events",
			Column: "latency_ms",
			Rows:   []col.Value{col.Int(4), col.Int(6), col.Int(1), col.Int(9)},
			Start:  now,
			End:    now.Add(time.Hour),
			Seal:   true,
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	engine := NewEngine(store, parts, metaRegistry, dicts)
	plan, err := engine.Build("events", now, now.Add(time.Hour), []string{"latency_ms"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	result, err := engine.Execute(plan)
	if err != nil {
		t.Fatal(err)
	}
	if result.Scanned != 4 || len(result.Rows) != 1 {
		t.Fatalf("unexpected result: scanned=%d rows=%d", result.Scanned, len(result.Rows))
	}
}

func TestEmptyResultForNoData(t *testing.T) {
	store := col.NewStore(t.TempDir())
	parts := part.NewRegistry()
	metaRegistry := meta.NewRegistry()
	dicts := compress.NewDictRegistry()
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	if _, err := metaRegistry.Update(meta.BuildSchema("none", []string{"x"}, nil), now); err != nil {
		t.Fatal(err)
	}
	engine := NewEngine(store, parts, metaRegistry, dicts)
	plan, err := engine.Build("none", now, now.Add(time.Hour), []string{"x"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	result, err := engine.Execute(plan)
	if err != nil {
		t.Fatal(err)
	}
	if result.Scanned != 0 || len(result.Rows) != 0 {
		t.Fatalf("expected empty result, got scanned=%d rows=%d", result.Scanned, len(result.Rows))
	}
}
