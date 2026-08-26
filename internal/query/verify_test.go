package query_test

import (
	"testing"
	"time"

	"columnstore/internal/col"
	"columnstore/internal/compress"
	"columnstore/internal/load"
	"columnstore/internal/meta"
	"columnstore/internal/part"
	"columnstore/internal/query"
)

func TestQueryUsesLatestColumnMeta(t *testing.T) {
	store := col.NewStore(t.TempDir())
	parts := part.NewRegistry()
	metaRegistry := meta.NewRegistry()
	dicts := compress.NewDictRegistry()
	loader := load.NewLoader(store, parts, metaRegistry)
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	if _, err := loader.Load("events", load.NewSliceSource([]load.Batch{
		{
			Table:  "events",
			Column: "latency_ms",
			Rows:   []col.Value{col.Int(10), col.Int(20)},
			Start:  now,
			End:    now.Add(time.Hour),
			Seal:   true,
		},
	})); err != nil {
		t.Fatal(err)
	}
	engine := query.NewEngine(store, parts, metaRegistry, dicts)
	if _, err := engine.Build("events", now, now.Add(time.Hour), []string{"latency_ms"}, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := loader.Load("events", load.NewSliceSource([]load.Batch{
		{
			Table:  "events",
			Column: "region",
			Rows:   []col.Value{col.Str("east"), col.Str("west")},
			Start:  now,
			End:    now.Add(time.Hour),
			Seal:   true,
		},
		{
			Table:  "events",
			Column: "latency_ms",
			Rows:   []col.Value{col.Int(30), col.Int(40)},
			Start:  now,
			End:    now.Add(time.Hour),
			Seal:   true,
		},
	})); err != nil {
		t.Fatal(err)
	}
	plan, err := engine.Build("events", now, now.Add(time.Hour), []string{"latency_ms"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	result, err := engine.Execute(plan)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Rows) != 1 || result.Rows[0].Values[0].Num != 4 {
		t.Fatalf("query must parse new blocks with the latest schema, got %+v", result.Rows)
	}
}
