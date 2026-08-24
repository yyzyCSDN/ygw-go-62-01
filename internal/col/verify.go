package col

import (
	"encoding/json"
	"fmt"
)

// VerifySegment reads a block's segment file and confirms every expected row
// survived the round trip. It returns an error when the payload is truncated
// or the row count disagrees with the in-memory block.
func (s *Store) VerifySegment(block *ColumnBlock) error {
	if block == nil {
		return ErrBlockNotFound
	}
	payload, err := s.ReadBlockBytes(block)
	if err != nil {
		return err
	}
	if len(payload) == 0 {
		return ErrSegmentMissing
	}
	var rows []Value
	if err := json.Unmarshal(payload, &rows); err != nil {
		return err
	}
	if len(rows) != len(block.Rows) {
		return fmt.Errorf("segment row count %d does not match block rows %d", len(rows), len(block.Rows))
	}
	return nil
}

// VerifyAllSegments checks every block that has a backing segment.
func (s *Store) VerifyAllSegments() (int, error) {
	verified := 0
	for _, block := range s.ListBlocks() {
		if block.Segment == "" {
			continue
		}
		if err := s.VerifySegment(block); err != nil {
			return verified, err
		}
		verified++
	}
	return verified, nil
}
