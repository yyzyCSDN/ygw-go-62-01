package load

import (
	"fmt"
	"time"

	"columnstore/internal/col"
)

// Batch is one row group destined for a single partition.
type Batch struct {
	Table   string
	Column  string
	Rows    []col.Value
	Start   time.Time
	End     time.Time
	Seal    bool
	loadSeq int
}

// Describe renders the batch header for logs.
func (b Batch) Describe() string {
	return fmt.Sprintf("%s.%s rows=%d [%s,%s)", b.Table, b.Column, len(b.Rows),
		b.Start.Format("15:04:05"), b.End.Format("15:04:05"))
}
