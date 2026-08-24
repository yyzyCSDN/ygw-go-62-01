package scan_test

import (
	"testing"

	"columnstore/internal/col"
	"columnstore/internal/scan"
)

func TestVectorizedScanKeepsTrailingRows(t *testing.T) {
	rows := make([]col.Value, 10)
	for i := range rows {
		rows[i] = col.Int(int64(i + 1))
	}
	visited := 0
	scan.ScanChunked(rows, 4, func(batch []col.Value) {
		visited += len(batch)
	})
	if visited != 10 {
		t.Fatalf("chunked scan must cover all rows, visited %d of 10", visited)
	}
}
