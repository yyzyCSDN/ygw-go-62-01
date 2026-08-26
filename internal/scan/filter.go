package scan

import "columnstore/internal/col"

// Filter applies a predicate to a full row slice without pushdown.
func Filter(rows []col.Value, pred Predicate) []col.Value {
	out := make([]col.Value, 0, len(rows))
	for _, row := range rows {
		if pred.Match(row) {
			out = append(out, row)
		}
	}
	return out
}
