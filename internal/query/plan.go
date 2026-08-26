package query

import (
	"time"

	"columnstore/internal/meta"
	"columnstore/internal/scan"
)

// Plan is a prepared query over a table.
type Plan struct {
	Table   string
	From    time.Time
	To      time.Time
	Meta    *meta.ColumnMeta
	Version uint64
	Columns []string
	Pred    *scan.Predicate
	Strict  bool
}

// RangeText formats the query window for diagnostics.
func (p *Plan) RangeText() string {
	return p.From.Format("2006-01-02T15:04:05") + ".." + p.To.Format("2006-01-02T15:04:05")
}
