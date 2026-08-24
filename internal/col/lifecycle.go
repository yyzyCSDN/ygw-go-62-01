package col

import "time"

// ExpiredAt reports whether a block is older than the given retention window
// at the provided instant. The boundary is strict: a block whose age equals
// the window is still considered fresh.
func ExpiredAt(created time.Time, window time.Duration, now time.Time) bool {
	return now.Sub(created) > window
}

// ReclaimBlock marks a block reclaimed and removes its backing segment.
func (s *Store) ReclaimBlock(block *ColumnBlock) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if block == nil {
		return ErrBlockNotFound
	}
	if _, ok := s.blocks[block.ID]; !ok {
		return ErrBlockNotFound
	}
	if block.State == StateReclaimed {
		return nil
	}
	block.State = StateReclaimed
	delete(s.blocks, block.ID)
	if block.Segment != "" {
		_ = s.segments.RemoveSegment(block.Segment)
	}
	return nil
}
