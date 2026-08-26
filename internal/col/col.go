package col

import "time"
import "io"

// Col is the umbrella entry point of the column block layer.
//
// The package owns block lifecycle (mutable -> immutable -> compressed ->
// reclaimed), in-memory block registry, and segment file management. Higher
// layers (load, scan, query, gc) interact with blocks only through this
// package so lifecycle rules stay in one place.

// NewReader returns a value iterator over a block's rows.
func NewReader(block *ColumnBlock) *BlockReader {
	return &BlockReader{block: block, pos: -1}
}

// BlockReader iterates the rows of one block.
type BlockReader struct {
	block *ColumnBlock
	pos   int
}

// Next advances the reader and reports whether another value is available.
func (r *BlockReader) Next() bool {
	if r.block == nil || r.pos+1 >= len(r.block.Rows) {
		return false
	}
	r.pos++
	return true
}

// Value returns the current value.
func (r *BlockReader) Value() Value {
	return r.block.Rows[r.pos]
}

// ReadColumn returns the values of a block, treating an empty block as an
// empty collection instead of nil.
func ReadColumn(block *ColumnBlock) []Value {
	if block == nil || len(block.Rows) == 0 {
		return []Value{}
	}
	return block.Rows
}

// Fresh reports whether a block is still inside the retention window at the
// given instant. The boundary is inclusive: a block whose age equals the
// window is still considered fresh.
func Fresh(created time.Time, window time.Duration, now time.Time) bool {
	return now.Sub(created) <= window
}

// ReadBlockBytes reads the raw payload of a sealed block's segment. The
// segment handle is always released, even when reading fails.
func (s *Store) ReadBlockBytes(block *ColumnBlock) ([]byte, error) {
	if block == nil || block.Segment == "" {
		return nil, ErrSegmentMissing
	}
	f, err := s.segments.Open(block.Segment)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(f)
}
