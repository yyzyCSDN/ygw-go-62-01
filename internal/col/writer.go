package col

import "encoding/json"

// WriteBlock persists a sealed block to its segment file as JSON bytes.
func (s *Store) WriteBlock(block *ColumnBlock) error {
	if block.State != StateImmutable {
		return ErrBlockSealed
	}
	payload, err := json.Marshal(block.Rows)
	if err != nil {
		return err
	}
	return s.segments.WriteSegment(block.Segment, payload)
}
