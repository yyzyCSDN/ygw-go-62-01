package query

import (
	"columnstore/internal/col"
	"columnstore/internal/scan"
)

// Aggregate carries rolling totals while a query executes.
type Aggregate struct {
	Column string
	Count  int64
	Sum    float64
}

// Add folds one scanned value into the aggregate.
func (a *Aggregate) Add(v col.Value) {
	a.Count++
	a.Sum += v.Num
}

// FoldAccumulator folds a scan accumulator into the query aggregate.
func (a *Aggregate) FoldAccumulator(acc *scan.Accumulator) {
	a.Count += acc.Count
	a.Sum += acc.Sum
}
