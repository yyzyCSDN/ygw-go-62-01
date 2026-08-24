package load

import (
	"errors"
	"fmt"

	"columnstore/internal/col"
)

var (
	ErrEmptyBatch   = errors.New("batch carries no rows")
	ErrKindMismatch = errors.New("batch mixes value kinds")
)

// Validate checks a batch before it is persisted.
func Validate(b Batch) error {
	if len(b.Rows) == 0 {
		return ErrEmptyBatch
	}
	if b.End.Before(b.Start) {
		return fmt.Errorf("batch interval reversed")
	}
	kind := b.Rows[0].Kind
	for _, row := range b.Rows {
		if row.Kind != kind {
			return fmt.Errorf("%w at row %d", ErrKindMismatch, 0)
		}
	}
	return nil
}

// ColumnKind returns the uniform kind of a validated batch.
func ColumnKind(b Batch) col.Kind {
	if len(b.Rows) == 0 {
		return col.KindInt
	}
	return b.Rows[0].Kind
}
