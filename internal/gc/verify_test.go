package gc_test

import (
	"testing"
	"time"

	"columnstore/internal/col"
	"columnstore/internal/gc"
)

func TestReclaimKeepsFreshBlocks(t *testing.T) {
	store := col.NewStore(t.TempDir())
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	block := col.NewBlock("t", "c", []col.Value{col.Int(1)}, now)
	if _, err := store.AddBlock(block); err != nil {
		t.Fatal(err)
	}
	ttl := time.Hour
	reclaimer := gc.New(store, gc.DefaultPolicy().WithTTL(ttl))
	if _, err := reclaimer.Run(now.Add(ttl)); err != nil {
		t.Fatal(err)
	}
	if store.BlockCount() != 1 {
		t.Fatal("a block whose age equals the retention window must stay fresh")
	}
}
