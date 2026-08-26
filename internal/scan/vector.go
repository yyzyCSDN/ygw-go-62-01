package scan

import "columnstore/internal/col"

// Accumulator aggregates numeric values while scanning.
type Accumulator struct {
	Count int64
	Sum   float64
	Min   float64
	Max   float64
	seen  bool
}

// Add folds one value into the accumulator.
func (a *Accumulator) Add(v col.Value) {
	a.Count++
	if !a.seen {
		a.Min = v.Num
		a.Max = v.Num
		a.seen = true
	}
	if v.Num < a.Min {
		a.Min = v.Num
	}
	if v.Num > a.Max {
		a.Max = v.Num
	}
	a.Sum += v.Num
}

// AddRows folds a row slice into the accumulator.
func (a *Accumulator) AddRows(rows []col.Value) {
	for _, v := range rows {
		a.Add(v)
	}
}

// ChunkEnd computes the exclusive end of a chunk that never runs past the
// row slice, so the final partial chunk is always processed.
func ChunkEnd(start, chunk, total int) int {
	end := start + chunk
	if end > total {
		return total
	}
	return end
}
