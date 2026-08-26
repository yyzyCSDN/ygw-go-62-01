package meta

import (
	"testing"
	"time"

	"columnstore/internal/col"
)

func TestRegistryUpdateAndLookup(t *testing.T) {
	registry := NewRegistry()
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	schema := BuildSchema("events", []string{"ts", "latency_ms"}, []col.Kind{col.KindFloat, col.KindInt})
	first, err := registry.Update(schema, now)
	if err != nil {
		t.Fatal(err)
	}
	second, err := registry.Update(schema, now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if first.SchemaVersion.Number != 1 || second.SchemaVersion.Number != 2 {
		t.Fatalf("versions must advance: %d -> %d", first.SchemaVersion.Number, second.SchemaVersion.Number)
	}
	latest, err := registry.Lookup("events")
	if err != nil {
		t.Fatal(err)
	}
	if !latest.AtVersion(2) {
		t.Fatal("latest meta must be version 2")
	}
}

func TestSchemaPosition(t *testing.T) {
	schema := BuildSchema("events", []string{"a", "b"}, []col.Kind{col.KindInt, col.KindString})
	if schema.Position("b") != 1 || schema.Position("missing") != -1 {
		t.Fatal("position lookup is wrong")
	}
}

func TestCatalogReadsLatest(t *testing.T) {
	registry := NewRegistry()
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	schema := BuildSchema("events", []string{"x"}, nil)
	if _, err := registry.Update(schema, now); err != nil {
		t.Fatal(err)
	}
	catalog := NewCatalog(registry)
	got, err := catalog.SchemaFor("events")
	if err != nil {
		t.Fatal(err)
	}
	if got.Table != "events" {
		t.Fatal("catalog returned the wrong table")
	}
}
