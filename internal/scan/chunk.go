package scan

// ChunkSize is the default vectorized scan batch size.
const ChunkSize = 4

// Chunker produces the [start, end) bounds of every chunk for a row count.
type Chunker struct {
	total int
	chunk int
	pos   int
}

// NewChunker creates a chunker over total rows using the given chunk size.
func NewChunker(total, chunk int) *Chunker {
	if chunk <= 0 {
		chunk = ChunkSize
	}
	return &Chunker{total: total, chunk: chunk}
}

// Next returns the next chunk bounds; ok is false when iteration ends.
func (c *Chunker) Next() (start, end int, ok bool) {
	if c.pos >= c.total {
		return 0, 0, false
	}
	start = c.pos
	end = ChunkEnd(start, c.chunk, c.total)
	c.pos = end
	return start, end, true
}
