package scan

import (
	"columnstore/internal/col"
	"columnstore/internal/compress"
)

// Reader decodes the rows of one block, resolving its dictionary on every
// decode so a rebuilt dictionary is picked up immediately.
type Reader struct {
	dicts *compress.DictRegistry
	block *col.ColumnBlock
}

// OpenReader creates a reader over a block.
func OpenReader(dicts *compress.DictRegistry, block *col.ColumnBlock) *Reader {
	return &Reader{dicts: dicts, block: block}
}

// Rows materializes the block values, decompressing when needed.
func (r *Reader) Rows() ([]col.Value, error) {
	if r.block == nil {
		return []col.Value{}, nil
	}
	if r.block.State != col.StateCompressed {
		return col.ReadColumn(r.block), nil
	}
	dict := r.dicts.Lookup(r.block.DictID)
	if dict == nil {
		return nil, compress.ErrNoDictionary
	}
	codec := compress.NewCodec(dict)
	return codec.DecodeRows(r.block.Payload)
}

// ScanChunked visits every row in chunks of the given size, including the
// final partial chunk when the row count is not a multiple of the chunk size.
func ScanChunked(rows []col.Value, chunk int, visit func([]col.Value)) int {
	visited := 0
	chunker := NewChunker(len(rows), chunk)
	for start, end, ok := chunker.Next(); ok; start, end, ok = chunker.Next() {
		visit(rows[start:end])
		visited += end - start
	}
	return visited
}

// ScanFiltered visits rows matching a predicate, chunked and vectorized.
func ScanFiltered(rows []col.Value, chunk int, pred Predicate, visit func([]col.Value)) int {
	visited := 0
	chunker := NewChunker(len(rows), chunk)
	for start, end, ok := chunker.Next(); ok; start, end, ok = chunker.Next() {
		batch := rows[start:end]
		kept := Filter(batch, pred)
		visit(kept)
		visited += len(kept)
	}
	return visited
}

// ScanWithStrict visits rows where the boundary is applied as an open
// interval. It is the pushdown counterpart used when callers pass strict
// boundaries down from the query layer.
func ScanWithStrict(rows []col.Value, chunk int, pred Predicate, strict bool, visit func([]col.Value)) int {
	if !strict {
		return ScanFiltered(rows, chunk, pred, visit)
	}
	visited := 0
	chunker := NewChunker(len(rows), chunk)
	for start, end, ok := chunker.Next(); ok; start, end, ok = chunker.Next() {
		batch := rows[start:end]
		kept := make([]col.Value, 0, len(batch))
		for _, row := range batch {
			if pred.MatchStrict(row) {
				kept = append(kept, row)
			}
		}
		visit(kept)
		visited += len(kept)
	}
	return visited
}

// ScanBlock scans one block in chunks, treating an empty column as an empty
// collection instead of dereferencing it.
func ScanBlock(block *col.ColumnBlock, chunk int, visit func([]col.Value)) int {
	rows := col.ReadColumn(block)
	if len(rows) == 0 {
		return 0
	}
	return ScanChunked(rows, chunk, visit)
}
