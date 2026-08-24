package gc

import (
	"time"

	"columnstore/internal/col"
)

// ExpiredAt decides whether a block has outlived the retention window at the
// given instant. The boundary is strict: a block whose age equals the window
// is still fresh and must never be reclaimed.
func ExpiredAt(block *col.ColumnBlock, ttl time.Duration, now time.Time) bool {
	return col.ExpiredAt(block.CreatedAt, ttl, now)
}
