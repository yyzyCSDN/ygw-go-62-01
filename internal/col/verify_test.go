package col_test

import (
	"testing"
	"time"

	"columnstore/internal/col"
)

func TestColumnSegmentHandleClosed(t *testing.T) {
	store := col.NewStore(t.TempDir())
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 20; i++ {
		block := col.NewBlock("t", "c", []col.Value{col.Int(int64(i))}, now)
		if _, err := store.AddBlock(block); err != nil {
			t.Fatal(err)
		}
		if err := block.Seal(block.Segment); err != nil {
			t.Fatal(err)
		}
		if err := store.WriteBlock(block); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := store.VerifyAllSegments(); err != nil {
		t.Fatal(err)
	}
	if store.OpenHandles() != 0 {
		t.Fatalf("segment handles leaked after verify: %d", store.OpenHandles())
	}
}
