package scan

import (
	"testing"
	"time"

	"columnstore/internal/col"
)

func TestChunkerCoversDivisibleRows(t *testing.T) {
	rows := make([]col.Value, 8)
	chunker := NewChunker(len(rows), 4)
	chunks := 0
	for _, _, ok := chunker.Next(); ok; _, _, ok = chunker.Next() {
		chunks++
	}
	if chunks != 2 {
		t.Fatalf("expected two chunks, got %d", chunks)
	}
}

func TestAccumulatorAggregates(t *testing.T) {
	acc := &Accumulator{}
	acc.AddRows([]col.Value{col.Int(3), col.Int(5), col.Int(2)})
	if acc.Count != 3 || acc.Sum != 10 || acc.Min != 2 || acc.Max != 5 {
		t.Fatalf("unexpected aggregate: %+v", acc)
	}
}

func TestFilterStrictBounds(t *testing.T) {
	rows := []col.Value{col.Int(1), col.Int(2), col.Int(3), col.Int(4)}
	pred := Predicate{Column: "x", Op: OpGT, Value: 2}
	kept := Filter(rows, pred)
	if len(kept) != 2 || kept[0].Num != 3 {
		t.Fatalf("unexpected filter result: %+v", kept)
	}
}

func TestStatsAndPushdown(t *testing.T) {
	block := col.NewBlock("t", "x", []col.Value{col.Int(4), col.Int(9), col.Int(6)}, time.Now())
	stats := Stats(block)
	if stats.Min != 4 || stats.Max != 9 || stats.Rows != 3 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
	if !PushdownEligible(stats, Predicate{Op: OpGT, Value: 100}) {
		t.Fatal("block with max below threshold must be pushdown-eligible")
	}
}

func TestScanChunkedDivisible(t *testing.T) {
	rows := []col.Value{col.Int(1), col.Int(2), col.Int(3), col.Int(4), col.Int(5), col.Int(6), col.Int(7), col.Int(8)}
	visited := ScanChunked(rows, 4, func([]col.Value) {})
	if visited != 8 {
		t.Fatalf("expected 8 visited rows, got %d", visited)
	}
}

func TestScanChunkedRemainder(t *testing.T) {
	// 10 rows, chunk size 4: the final partial chunk (rows[8:10]) must be
	// visited, not dropped.
	rows := []col.Value{col.Int(1), col.Int(2), col.Int(3), col.Int(4), col.Int(5), col.Int(6), col.Int(7), col.Int(8), col.Int(9), col.Int(10)}
	seen := make(map[int]bool, len(rows))
	visited := ScanChunked(rows, 4, func(batch []col.Value) {
		for _, v := range batch {
			seen[int(v.Num)] = true
		}
	})
	if visited != 10 {
		t.Fatalf("expected 10 visited rows, got %d", visited)
	}
	for i := 1; i <= 10; i++ {
		if !seen[i] {
			t.Fatalf("row %d was not visited by the chunked scan", i)
		}
	}
}

func TestScanChunkedAggregatesRemainder(t *testing.T) {
	// The whole-block aggregate must match a full scan even when the row
	// count is not a multiple of the chunk size; the tail rows must fold in.
	rows := []col.Value{col.Int(1), col.Int(2), col.Int(3), col.Int(4), col.Int(5), col.Int(6), col.Int(7), col.Int(8), col.Int(9), col.Int(10)}
	acc := &Accumulator{}
	ScanChunked(rows, 4, func(batch []col.Value) { acc.AddRows(batch) })
	if acc.Count != 10 || acc.Sum != 55 || acc.Min != 1 || acc.Max != 10 {
		t.Fatalf("aggregate lost tail rows: %+v", acc)
	}
}
