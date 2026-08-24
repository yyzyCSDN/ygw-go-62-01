package load

// Source streams batches into a loader.
type Source interface {
	Batches() []Batch
}

// SliceSource is an in-memory source over a static batch list.
type SliceSource struct {
	batch []Batch
}

// NewSliceSource wraps a static batch list as a source.
func NewSliceSource(batches []Batch) *SliceSource {
	return &SliceSource{batch: batches}
}

// Batches returns the wrapped batch list.
func (s *SliceSource) Batches() []Batch {
	return s.batch
}
