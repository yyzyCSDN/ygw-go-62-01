package load

import (
	"testing"
	"time"

	"columnstore/internal/col"
	"columnstore/internal/meta"
	"columnstore/internal/part"
)

// Throwaway verification (NOT one of the existing tests):
// batch load fails mid-way. Prior batches fully sealed (block + segment +
// partition) and an in-flight batch (block + segment written, Seal=false so
// no partition) must all leave ZERO residue after rollback.
func TestVerifyZZ_RollbackLeavesNoResidue(t *testing.T) {
	dir := t.TempDir()
	store := col.NewStore(dir)
	parts := part.NewRegistry()
	metaReg := meta.NewRegistry()
	loader := NewLoader(store, parts, metaReg)
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)

	// Batch 1: sealed  -> block + partition + segment file.
	// Batch 2: in-flight (Seal=false) -> block + segment file, NO partition.
	// Batch 3: invalid (empty rows) -> triggers rollback.
	src := NewSliceSource([]Batch{
		{Table: "events", Column: "latency_ms", Rows: []col.Value{col.Int(10)}, Start: now, End: now.Add(time.Hour), Seal: true},
		{Table: "events", Column: "latency_ms", Rows: []col.Value{col.Int(20)}, Start: now.Add(time.Hour), End: now.Add(2 * time.Hour), Seal: false},
		{Table: "events", Column: "latency_ms", Rows: []col.Value{}, Start: now.Add(2 * time.Hour), End: now.Add(3 * time.Hour), Seal: true},
	})

	if _, err := loader.Load("events", src); err == nil {
		t.Fatal("expected load to fail on the invalid batch")
	}

	if got := store.BlockCount(); got != 0 {
		t.Errorf("residue: store.BlockCount()=%d, want 0 (in-memory blocks leaked)", got)
	}
	if got := parts.Count(); got != 0 {
		t.Errorf("residue: parts.Count()=%d, want 0 (partitions leaked)", got)
	}
	if n, err := store.VerifyAllSegments(); err == nil && n != 0 {
		t.Errorf("residue: %d segment files still verifiable, want 0", n)
	}
	if got := loader.Manifest().TableBlocks("events"); got != 0 {
		t.Errorf("residue: manifest TableBlocks=%d, want 0", got)
	}
}
