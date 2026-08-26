package col

import (
	"testing"
	"time"
)

func TestBlockLifecycleSealAndRead(t *testing.T) {
	store := NewStore(t.TempDir())
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	block := NewBlock("t", "c", []Value{Int(1), Int(2)}, now)
	id, err := store.AddBlock(block)
	if err != nil {
		t.Fatal(err)
	}
	if err := block.Seal(block.Segment); err != nil {
		t.Fatal(err)
	}
	if err := store.WriteBlock(block); err != nil {
		t.Fatal(err)
	}
	got, err := store.ReadSegment(block)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[1].Num != 2 {
		t.Fatalf("unexpected rows after segment round trip: %+v", got)
	}
	if _, err := store.GetBlock(id); err != nil {
		t.Fatal(err)
	}
}

func TestBlockReaderIteratesValues(t *testing.T) {
	block := NewBlock("t", "c", []Value{Int(7), Int(8), Int(9)}, time.Now())
	reader := NewReader(block)
	sum := 0
	for reader.Next() {
		sum += int(reader.Value().Num)
	}
	if sum != 24 {
		t.Fatalf("expected 24, got %d", sum)
	}
}

func TestReclaimRemovesOldBlock(t *testing.T) {
	store := NewStore(t.TempDir())
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	block := NewBlock("t", "c", []Value{Int(1)}, now.Add(-48*time.Hour))
	_, err := store.AddBlock(block)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.ReclaimBlock(block); err != nil {
		t.Fatal(err)
	}
	if store.BlockCount() != 0 {
		t.Fatalf("reclaimed block still visible")
	}
}
