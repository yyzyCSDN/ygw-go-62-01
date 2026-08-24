package scan

import "columnstore/internal/col"

// BlockStats summarizes one block for predicate pushdown decisions.
type BlockStats struct {
	BlockID uint64
	Min     float64
	Max     float64
	Rows    int
}

// Stats computes the numeric range of a block.
func Stats(block *col.ColumnBlock) BlockStats {
	s := BlockStats{BlockID: block.ID}
	if block == nil || len(block.Rows) == 0 {
		return s
	}
	s.Min = block.Rows[0].Num
	s.Max = block.Rows[0].Num
	for _, v := range block.Rows {
		if v.Num < s.Min {
			s.Min = v.Num
		}
		if v.Num > s.Max {
			s.Max = v.Num
		}
	}
	s.Rows = len(block.Rows)
	return s
}

// PushdownEligible reports whether a predicate can skip a whole block.
func PushdownEligible(stats BlockStats, pred Predicate) bool {
	if stats.Rows == 0 {
		return true
	}
	switch pred.Op {
	case OpGE, OpGT:
		return stats.Max < pred.Value
	case OpLE, OpLT:
		return stats.Min > pred.Value
	case OpEQ:
		return stats.Min > pred.Value || stats.Max < pred.Value
	default:
		return false
	}
}
