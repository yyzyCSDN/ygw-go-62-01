package scan

import "columnstore/internal/col"

// PushdownScan visits only the rows of a block that match a predicate after
// the block-level stats say the block cannot be skipped entirely.
func PushdownScan(block *col.ColumnBlock, chunk int, pred Predicate, visit func([]col.Value)) int {
	stats := Stats(block)
	if PushdownEligible(stats, pred) {
		return 0
	}
	rows := col.ReadColumn(block)
	return ScanFiltered(rows, chunk, pred, visit)
}
