package col

import "encoding/json"

// ReadSegment loads block rows that were persisted with WriteBlock.
func (s *Store) ReadSegment(block *ColumnBlock) ([]Value, error) {
	payload, err := s.segments.ReadSegmentBytes(block.Segment)
	if err != nil {
		return nil, err
	}
	var rows []Value
	if err := json.Unmarshal(payload, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}
