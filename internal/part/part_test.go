package part

import (
	"testing"
	"time"
)

func TestRegistryRegisterRemove(t *testing.T) {
	registry := NewRegistry()
	start := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	p := NewPartition("events", start, start.Add(time.Hour), 0)
	id, err := registry.Register(p)
	if err != nil {
		t.Fatal(err)
	}
	if registry.Count() != 1 {
		t.Fatalf("expected one partition")
	}
	if err := registry.Remove(id); err != nil {
		t.Fatal(err)
	}
	if registry.Count() != 0 {
		t.Fatalf("expected zero partitions after remove")
	}
}

func TestPruneKeepsOnlyOverlapping(t *testing.T) {
	base := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	early := NewPartition("events", base.Add(-2*time.Hour), base.Add(-time.Hour), 1)
	mid := NewPartition("events", base, base.Add(time.Hour), 2)
	late := NewPartition("events", base.Add(3*time.Hour), base.Add(4*time.Hour), 3)
	kept := Prune([]*Partition{early, mid, late}, base, base.Add(2*time.Hour))
	if len(kept) != 1 || kept[0].ID != 2 {
		t.Fatalf("expected only the middle partition, got %+v", kept)
	}
}

func TestClampRangeNarrowsWindow(t *testing.T) {
	base := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	p := NewPartition("events", base.Add(30*time.Minute), base.Add(time.Hour), 1)
	lo, hi := ClampRange(p, base, base.Add(2*time.Hour))
	if !lo.Equal(base.Add(30 * time.Minute)) {
		t.Fatalf("expected lower bound clamped to partition start, got %v", lo)
	}
	if !hi.Equal(base.Add(time.Hour)) {
		t.Fatalf("expected upper bound clamped to partition end, got %v", hi)
	}
}

func TestPartitionStates(t *testing.T) {
	base := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	p := NewPartition("events", base, base.Add(time.Hour), 1)
	if !p.AddBlock(10) {
		t.Fatal("active partition must accept blocks")
	}
	p.Seal()
	if p.AddBlock(11) {
		t.Fatal("sealed partition must reject blocks")
	}
	if !p.Contains(base.Add(30 * time.Minute)) {
		t.Fatal("partition must contain interior instants")
	}
}
