package query

import (
	"columnstore/internal/col"
	"columnstore/internal/scan"
)

// executor owns the scan configuration shared by every block in a query.
type executor struct {
	chunk int
}

// newExecutor prepares a scan with the given chunk size.
func newExecutor(chunk int) *executor {
	return &executor{chunk: chunk}
}

// scanRows runs a chunked scan over one block's rows.
func (e *executor) scanRows(rows []col.Value, visit func([]col.Value)) int {
	return scan.ScanChunked(rows, e.chunk, visit)
}

// scanBlockFiltered applies pushdown stats before scanning a whole block.
func (e *executor) scanBlockFiltered(block *col.ColumnBlock, pred scan.Predicate, strict bool, visit func([]col.Value)) int {
	if strict {
		return scan.ScanWithStrict(col.ReadColumn(block), e.chunk, pred, true, visit)
	}
	return scan.PushdownScan(block, e.chunk, pred, visit)
}
